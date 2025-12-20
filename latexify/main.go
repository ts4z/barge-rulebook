package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/pflag"

	"github.com/ts4z/barge-rulebook/latexify/latex"
	"github.com/ts4z/barge-rulebook/latexify/mdbook"
)

// Config holds command-line options
type Config struct {
	SourceDir   string
	OutputFile  string
	SummaryFile string
}

// Regular expression to match sub-appendix titles like "Appendix A.1:", "Appendix D.2:"
var subAppendixPattern = regexp.MustCompile(`^Appendix [A-Z]+\.[0-9]+:`)

// isSubAppendix checks if a title represents a sub-appendix (e.g., "Appendix A.1:")
func isSubAppendix(title string) bool {
	return subAppendixPattern.MatchString(title)
}

func main() {
	config := parseFlags()

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	pflag.StringVar(&config.SourceDir, "source-dir", "src", "Source directory containing markdown files")
	pflag.StringVar(&config.OutputFile, "output", "rulebook.latex", "Output LaTeX file")
	pflag.StringVar(&config.SummaryFile, "summary", "src/SUMMARY.md", "Path to SUMMARY.md file")

	pflag.Parse()

	return config
}

func run(config *Config) error {
	fmt.Printf("Parsing %s...\n", config.SummaryFile)

	// Parse SUMMARY.md
	summary, err := mdbook.ParseSummary(config.SummaryFile)
	if err != nil {
		return fmt.Errorf("parsing summary: %w", err)
	}

	// Open output file
	out, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	writer := latex.NewSectionAwareWriter(out)

	// Write header comment
	writer.WriteComment("")
	writer.WriteComment(fmt.Sprintf(" %s (produced by latexify)", config.OutputFile))
	writer.WriteComment("")

	// Special files to insert at HR boundaries
	specialFiles := []string{"setup", "frontmatter", "mainmatter", "backmatter", "endmatter"}
	specialIndex := 0

	// Render special file at the beginning
	if err := renderSpecialFile(writer, config.SourceDir, specialFiles[specialIndex]); err != nil {
		return err
	}
	specialIndex++

	// Process each section
	for sectionIdx, section := range summary.Sections {
		fmt.Printf("Processing section %d with %d entries...\n", sectionIdx, len(section))

		// Insert special file at section boundaries (after first section)
		if sectionIdx > 0 && specialIndex < len(specialFiles) {
			if err := renderSpecialFile(writer, config.SourceDir, specialFiles[specialIndex]); err != nil {
				return err
			}
			specialIndex++
		}

		// Section 0 is the title page - render file content directly without chapter
		if sectionIdx == 0 {
			for _, node := range section {
				if node.FilePath != "" {
					if err := renderFile(writer, config.SourceDir, node.FilePath, 0); err != nil {
						return err
					}
				}
			}
			continue
		}

		// Determine if this is the appendix section (last section with content)
		isAppendix := sectionIdx == len(summary.Sections)-1
		appendixEmitted := false

		// Render each node in the section
		for _, node := range section {
			if err := renderNode(writer, config.SourceDir, node, isAppendix, &appendixEmitted); err != nil {
				return err
			}
		}
	}

	// Render final special file
	if specialIndex < len(specialFiles) {
		if err := renderSpecialFile(writer, config.SourceDir, specialFiles[specialIndex]); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully generated %s\n", config.OutputFile)
	return nil
}

// renderNode renders a summary node and its children
func renderNode(writer *latex.SectionAwareWriter, sourceDir string, node *mdbook.SummaryNode, isAppendix bool, appendixEmitted *bool) error {
	// Render section heading based on level
	if node.Level == 0 {
		if isAppendix && !*appendixEmitted {
			writer.WriteText("\\appendix\n")
			*appendixEmitted = true
		}
		// Check if this is a sub-appendix (e.g., "Appendix A.1:")
		if isAppendix && isSubAppendix(node.Title) {
			// Skip section heading for sub-appendices - don't emit anything
		} else {
			writer.WriteChapter(node.Title)
		}
	} else {
		writer.WriteSectionByLevel(node.Level, node.Title)
	}

	// Render the file content
	if node.FilePath != "" {
		// For sub-appendices, use sectionOffset 2 so h1 becomes subsection
		// For regular level 0 nodes, use level+1 (which is 1) so h1 becomes section
		sectionOffset := node.Level + 1
		if isAppendix && node.Level == 0 && isSubAppendix(node.Title) {
			sectionOffset = 2 // h1 in sub-appendix files becomes subsection
		}
		if err := renderFile(writer, sourceDir, node.FilePath, sectionOffset); err != nil {
			return fmt.Errorf("rendering %s: %w", node.FilePath, err)
		}
	}

	// Render children
	for _, child := range node.Children {
		if err := renderNode(writer, sourceDir, child, isAppendix, appendixEmitted); err != nil {
			return err
		}
	}

	return nil
}

// renderFile renders either a .latex file or a .md file converted to LaTeX
func renderFile(writer *latex.SectionAwareWriter, sourceDir string, basePath string, sectionOffset int) error {
	// Try .latex file first
	latexPath := filepath.Join(sourceDir, basePath+".latex")
	if content, err := os.ReadFile(latexPath); err == nil {
		writer.WriteComment(fmt.Sprintf(" %s (raw LaTeX)", basePath))
		writer.WriteText(string(content))
		fmt.Printf("  rendered %s as LaTeX\n", basePath)
		return nil
	}

	// Fall back to .md file
	mdPath := filepath.Join(sourceDir, basePath+".md")
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("reading markdown file: %w", err)
	}

	// Convert markdown to LaTeX
	// Pass the filename (basePath + ".md") for label generation
	filename := filepath.Base(basePath) + ".md"
	latex, err := latex.RenderMarkdownToLatex(content, sectionOffset, filename)
	if err != nil {
		return fmt.Errorf("converting markdown to latex: %w", err)
	}

	writer.WriteComment(fmt.Sprintf(" %s (rendered Markdown)", basePath+".md"))
	writer.WriteText(latex)
	fmt.Printf("  rendered %s as Markdown\n", basePath)
	return nil
}

// renderSpecialFile renders a special LaTeX file (setup, frontmatter, etc.)
func renderSpecialFile(writer *latex.SectionAwareWriter, sourceDir string, name string) error {
	path := filepath.Join(sourceDir, name+".latex")
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading special file %s: %w", name, err)
	}

	writer.WriteComment(fmt.Sprintf(" %s (special file)", name))
	writer.WriteText(string(content))
	fmt.Printf("  rendered special: %s\n", name)
	return nil
}
