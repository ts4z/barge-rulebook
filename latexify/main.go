package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

// Config holds command-line options
type Config struct {
	SourceDir   string
	OutputFile  string
	SummaryFile string
}

// DocumentSection represents a major section of the document
type DocumentSection int

const (
	SectionTitle DocumentSection = iota
	SectionFrontmatter
	SectionMainmatter
	SectionBackmatter
)

// SummaryNode represents an entry in SUMMARY.md
type SummaryNode struct {
	Title    string
	FilePath string // relative path without .md extension
	Level    int    // nesting level (0 = top level)
	Children []*SummaryNode
}

// SummaryStructure represents the parsed SUMMARY.md
type SummaryStructure struct {
	Sections [][]*SummaryNode // sections separated by HRs
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
	summary, err := ParseSummary(config.SummaryFile)
	if err != nil {
		return fmt.Errorf("parsing summary: %w", err)
	}

	// Open output file
	out, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	writer := NewLatexWriter(out)

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
func renderNode(writer *SectionWriter, sourceDir string, node *SummaryNode, isAppendix bool, appendixEmitted *bool) error {
	// Render section heading based on level
	if node.Level == 0 {
		if isAppendix && !*appendixEmitted {
			writer.WriteText("\\appendix\n")
			*appendixEmitted = true
		}
		writer.WriteChapter(node.Title)
	} else {
		writer.WriteSectionByLevel(node.Level, node.Title)
	}

	// Render the file content
	if node.FilePath != "" {
		if err := renderFile(writer, sourceDir, node.FilePath, node.Level); err != nil {
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
func renderFile(writer *SectionWriter, sourceDir string, basePath string, level int) error {
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
	// Use level+1 as section offset since the chapter heading is already written
	latex, err := RenderMarkdownToLatex(content, level+1)
	if err != nil {
		return fmt.Errorf("converting markdown to latex: %w", err)
	}

	writer.WriteComment(fmt.Sprintf(" %s (rendered Markdown)", basePath+".md"))
	writer.WriteText(latex)
	fmt.Printf("  rendered %s as Markdown\n", basePath)
	return nil
}

// renderSpecialFile renders a special LaTeX file (setup, frontmatter, etc.)
func renderSpecialFile(writer *SectionWriter, sourceDir string, name string) error {
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
