package mdbook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

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

// ParseSummary reads and parses SUMMARY.md into a structured representation
func ParseSummary(summaryPath string) (*SummaryStructure, error) {
	content, err := os.ReadFile(summaryPath)
	if err != nil {
		return nil, fmt.Errorf("reading summary file: %w", err)
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(content))

	structure := &SummaryStructure{
		Sections: make([][]*SummaryNode, 0),
	}

	currentSection := make([]*SummaryNode, 0)
	var listStack []*SummaryNode // stack to track nesting

	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.ThematicBreak:
			// HR separates sections
			if len(currentSection) > 0 || len(structure.Sections) == 0 {
				structure.Sections = append(structure.Sections, currentSection)
				currentSection = make([]*SummaryNode, 0)
				listStack = nil
			}

		case *ast.List:
			// Track list depth for nesting level
			// The list itself doesn't create a node, but its items do

		case *ast.ListItem:
			// Extract link from list item
			summaryNode := extractLinkFromNode(node, content)
			if summaryNode != nil {
				// Determine nesting level based on parent lists
				level := 0
				parent := node.Parent()
				for parent != nil {
					if _, ok := parent.(*ast.List); ok {
						level++
					}
					parent = parent.Parent()
				}
				summaryNode.Level = level - 1 // 0-indexed

				// Attach to parent or section
				if level == 1 {
					// Top-level item
					currentSection = append(currentSection, summaryNode)
					if len(listStack) > 0 {
						listStack = listStack[:0]
					}
					listStack = append(listStack, summaryNode)
				} else if level > 1 && len(listStack) >= level-1 {
					// Nested item
					parent := listStack[level-2]
					parent.Children = append(parent.Children, summaryNode)
					// Update stack: trim to keep only parents up to level-1
					if len(listStack) >= level-1 {
						listStack = listStack[:level-1]
					}
					listStack = append(listStack, summaryNode)
				}
			}

		case *ast.Paragraph:
			// Check if this paragraph contains a link (for non-list items)
			if _, ok := node.Parent().(*ast.ListItem); !ok {
				// Extract ALL links from this paragraph (not just the first one)
				links := extractAllLinksFromNode(node, content)
				for _, summaryNode := range links {
					summaryNode.Level = 0
					currentSection = append(currentSection, summaryNode)
				}
			}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	// Add final section if exists
	if len(currentSection) > 0 {
		structure.Sections = append(structure.Sections, currentSection)
	}

	return structure, nil
}

// extractAllLinksFromNode extracts all link information from a node
func extractAllLinksFromNode(node ast.Node, source []byte) []*SummaryNode {
	var result []*SummaryNode

	// Walk children to find all links (may be nested in paragraphs or text blocks)
	var findLinks func(ast.Node)
	findLinks = func(n ast.Node) {
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			if link, ok := child.(*ast.Link); ok {
				// Extract title from link text
				title := ""
				for linkChild := link.FirstChild(); linkChild != nil; linkChild = linkChild.NextSibling() {
					if text, ok := linkChild.(*ast.Text); ok {
						title += string(text.Segment.Value(source))
					}
				}

				// Extract and normalize file path
				dest := string(link.Destination)
				dest = strings.TrimPrefix(dest, "./")
				dest = strings.TrimSuffix(dest, ".md")

				result = append(result, &SummaryNode{
					Title:    title,
					FilePath: dest,
					Children: make([]*SummaryNode, 0),
				})
			}
			// Recursively search in paragraphs and text blocks
			if _, ok := child.(*ast.Paragraph); ok {
				findLinks(child)
			}
			if _, ok := child.(*ast.TextBlock); ok {
				findLinks(child)
			}
		}
	}

	findLinks(node)
	return result
}

// extractLinkFromNode extracts link information from a node (first link only)
func extractLinkFromNode(node ast.Node, source []byte) *SummaryNode {
	var link *ast.Link

	// Walk children to find link (may be nested in paragraph or text block)
	var findLink func(ast.Node) *ast.Link
	findLink = func(n ast.Node) *ast.Link {
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			if l, ok := child.(*ast.Link); ok {
				return l
			}
			// Recursively search in paragraphs and text blocks
			if _, ok := child.(*ast.Paragraph); ok {
				if found := findLink(child); found != nil {
					return found
				}
			}
			if _, ok := child.(*ast.TextBlock); ok {
				if found := findLink(child); found != nil {
					return found
				}
			}
		}
		return nil
	}

	link = findLink(node)

	if link == nil {
		return nil
	}

	// Extract title from link text
	title := ""
	for child := link.FirstChild(); child != nil; child = child.NextSibling() {
		if text, ok := child.(*ast.Text); ok {
			title += string(text.Segment.Value(source))
		}
	}

	// Extract and normalize file path
	dest := string(link.Destination)
	dest = strings.TrimPrefix(dest, "./")
	dest = strings.TrimSuffix(dest, ".md")

	return &SummaryNode{
		Title:    title,
		FilePath: dest,
		Children: make([]*SummaryNode, 0),
	}
}

// GetFileBasename extracts the basename without extension from a file path
func GetFileBasename(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
