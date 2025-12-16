//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			if err := processFile(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func processFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse the markdown
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(content))

	// Find the first heading
	var firstHeading *ast.Heading
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if h, ok := n.(*ast.Heading); ok && firstHeading == nil {
			firstHeading = h
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	if firstHeading == nil {
		fmt.Printf("No heading found in %s\n", path)
		return nil
	}

	// Check if it's a level 1 heading
	if firstHeading.Level != 1 {
		fmt.Printf("First heading in %s is not level 1 (it's level %d), skipping\n", path, firstHeading.Level)
		return nil
	}

	// Find the position of the heading in the source
	lines := firstHeading.Lines()
	if lines.Len() == 0 {
		fmt.Printf("No lines for heading in %s\n", path)
		return nil
	}

	// Get the start and end of the heading
	firstLine := lines.At(0)

	// Convert byte positions to line numbers
	startPos := firstLine.Start

	// Handle setext-style headings (underlined with ===)
	// Check if the next line after the heading might be the underline
	sourceLines := strings.Split(string(content), "\n")

	// Find which line number contains startPos
	currentPos := 0
	startLineNum := 0
	for i, line := range sourceLines {
		if currentPos >= startPos {
			startLineNum = i
			break
		}
		currentPos += len(line) + 1 // +1 for newline
	}

	// Check if this is a setext heading by looking at the next line
	endLineNum := startLineNum
	if startLineNum+1 < len(sourceLines) {
		nextLine := strings.TrimSpace(sourceLines[startLineNum+1])
		if strings.HasPrefix(nextLine, "===") || strings.HasPrefix(nextLine, "---") {
			endLineNum = startLineNum + 1
		}
	}

	// For ATX-style headings (starting with #), just remove that line
	// Build new content without the heading
	var newLines []string
	for i, line := range sourceLines {
		if i < startLineNum || i > endLineNum {
			newLines = append(newLines, line)
		}
	}

	// Remove leading blank lines
	for len(newLines) > 0 && strings.TrimSpace(newLines[0]) == "" {
		newLines = newLines[1:]
	}

	newContent := strings.Join(newLines, "\n")

	// Write back
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Removed first heading from %s\n", path)
	return nil
}
