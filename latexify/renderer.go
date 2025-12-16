package main

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// LatexRenderer renders Goldmark AST to LaTeX
type LatexRenderer struct {
	sectionOffset       int            // offset for section levels (0=chapter, 1=section, etc)
	seenFirstH1         bool           // track if we've seen the first h1 to skip it
	inSkippedH1         bool           // track if we're currently in the skipped h1
	footnotes           map[int]string // map of footnote index to content
	currentFootnote     *bytes.Buffer  // buffer for current footnote being collected
	currentDefinitionDD int            // count of DefinitionDescription nodes for current term
}

// NewLatexRenderer creates a new LaTeX renderer
func NewLatexRenderer(sectionOffset int) *LatexRenderer {
	return &LatexRenderer{
		sectionOffset:       sectionOffset,
		seenFirstH1:         false,
		inSkippedH1:         false,
		footnotes:           make(map[int]string),
		currentFootnote:     nil,
		currentDefinitionDD: 0,
	}
}

// RegisterFuncs implements renderer.NodeRenderer
func (r *LatexRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Block elements
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// Table elements
	reg.Register(extast.KindTable, r.renderTable)
	reg.Register(extast.KindTableHeader, r.renderTableHeader)
	reg.Register(extast.KindTableRow, r.renderTableRow)
	reg.Register(extast.KindTableCell, r.renderTableCell)

	// Definition list elements
	reg.Register(extast.KindDefinitionList, r.renderDefinitionList)
	reg.Register(extast.KindDefinitionTerm, r.renderDefinitionTerm)
	reg.Register(extast.KindDefinitionDescription, r.renderDefinitionDescription)

	// Footnote elements
	reg.Register(extast.KindFootnoteLink, r.renderFootnoteLink)
	reg.Register(extast.KindFootnote, r.renderFootnote)
	reg.Register(extast.KindFootnoteList, r.renderFootnoteList)
	reg.Register(extast.KindFootnoteBacklink, r.renderFootnoteBacklink)

	// Inline elements
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
}

func (r *LatexRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Nothing to do for document node
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)

	// Skip the first h1 heading (it's redundant with SUMMARY.md title)
	if n.Level == 1 && !r.seenFirstH1 {
		if entering {
			r.seenFirstH1 = true
			r.inSkippedH1 = true
			return ast.WalkSkipChildren, nil
		}
	}

	// If we're exiting the skipped h1, don't write closing brace
	if r.inSkippedH1 && !entering {
		r.inSkippedH1 = false
		return ast.WalkContinue, nil
	}

	if entering {
		// Note: Headings in the individual MD files will be rendered according to their level
		// The SUMMARY.md structure determines the actual chapter/section hierarchy
		level := n.Level - 1 + r.sectionOffset
		switch level {
		case 0:
			w.WriteString("\\chapter{")
		case 1:
			w.WriteString("\\section{")
		case 2:
			w.WriteString("\\subsection{")
		case 3:
			w.WriteString("\\subsubsection{")
		default:
			w.WriteString("\\paragraph{")
		}
	} else {
		w.WriteString("}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Start paragraph
	} else {
		// End paragraph with double newline
		w.WriteString("\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\\begin{quotation}\n")
	} else {
		w.WriteString("\\end{quotation}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.CodeBlock)
		w.WriteString("\\begin{verbatim}\n")
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}
		w.WriteString("\\end{verbatim}\n\n")
	}
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.FencedCodeBlock)
		lang := string(n.Language(source))

		if lang != "" {
			// Use listings package for syntax highlighting
			fmt.Fprintf(w, "\\begin{lstlisting}[language=%s]\n", lang)
		} else {
			w.WriteString("\\begin{verbatim}\n")
		}

		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}

		if lang != "" {
			w.WriteString("\\end{lstlisting}\n\n")
		} else {
			w.WriteString("\\end{verbatim}\n\n")
		}
	}
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Skip HTML blocks
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if entering {
		if n.IsOrdered() {
			w.WriteString("\\begin{enumerate}\n")
		} else {
			w.WriteString("\\begin{itemize}\n")
		}
	} else {
		if n.IsOrdered() {
			w.WriteString("\\end{enumerate}\n\n")
		} else {
			w.WriteString("\\end{itemize}\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\\item ")
	} else {
		w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\\vspace{1em}\\hrule\\vspace{1em}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderDefinitionList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\\begin{description}\n")
	} else {
		w.WriteString("\\end{description}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderDefinitionTerm(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Count how many DefinitionDescription siblings follow this term
		r.currentDefinitionDD = 0
		for sibling := node.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
			if sibling.Kind() == extast.KindDefinitionDescription {
				r.currentDefinitionDD++
			} else if sibling.Kind() == extast.KindDefinitionTerm {
				// Stop at the next term
				break
			}
		}

		w.WriteString("\\item[")
	} else {
		w.WriteString("] ")
		// If there are multiple definitions, start an enumerate
		if r.currentDefinitionDD > 1 {
			w.WriteString("\\begin{enumerate}\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderDefinitionDescription(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// If we have multiple descriptions, wrap each in \item
		if r.currentDefinitionDD > 1 {
			w.WriteString("\\item ")
		}
	} else {
		// Check if this is the last DefinitionDescription for this term
		isLast := true
		for sibling := node.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
			if sibling.Kind() == extast.KindDefinitionDescription {
				isLast = false
				break
			} else if sibling.Kind() == extast.KindDefinitionTerm {
				break
			}
		}

		if isLast && r.currentDefinitionDD > 1 {
			w.WriteString("\\end{enumerate}\n")
		}
		w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.Table)
	if entering {
		// Build column specification from alignments
		colSpec := ""
		for _, align := range n.Alignments {
			switch align {
			case extast.AlignLeft:
				colSpec += "l"
			case extast.AlignRight:
				colSpec += "r"
			case extast.AlignCenter:
				colSpec += "c"
			default:
				colSpec += "l" // default to left
			}
		}
		// Add spacing before table and center it
		w.WriteString("\n\\vspace{\\baselineskip}\n")
		w.WriteString("\\begin{center}\n")
		w.WriteString("\\begin{tabular}{" + colSpec + "}\n")
		w.WriteString("\\hline\n")
	} else {
		w.WriteString("\\hline\n")
		w.WriteString("\\end{tabular}\n")
		w.WriteString("\\end{center}\n")
		w.WriteString("\\vspace{\\baselineskip}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// TableHeader is just a container, don't output anything special
	// But add row ending and \hline after the header
	if !entering {
		w.WriteString(" \\\\\n")
		w.WriteString("\\hline\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.WriteString(" \\\\\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Add & separator between cells (except for first cell)
		if node.PreviousSibling() != nil {
			w.WriteString(" & ")
		}
		// Make header cells bold
		if node.Parent() != nil && node.Parent().Kind() == extast.KindTableHeader {
			w.WriteString("\\textbf{")
		}
	} else {
		// Close bold for header cells
		if node.Parent() != nil && node.Parent().Kind() == extast.KindTableHeader {
			w.WriteString("}")
		}
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderFootnoteLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*extast.FootnoteLink)
		// Output inline footnote with collected content
		if content, ok := r.footnotes[n.Index]; ok {
			w.WriteString("\\footnote{" + content + "}")
		} else {
			// Fallback if content not found
			w.WriteString(fmt.Sprintf("\\footnote{[footnote %d]}", n.Index))
		}
	}
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderFootnote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Footnotes are collected in the first pass and rendered inline
	// Skip them in the main rendering pass
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderFootnoteList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// FootnoteList container is not rendered - footnotes are inline
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderFootnoteBacklink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Backlinks are for HTML - skip them
	return ast.WalkSkipChildren, nil
}

func (r *LatexRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		segment := n.Segment
		value := segment.Value(source)

		// Escape special LaTeX characters
		escaped := escapeLatexText(string(value))
		w.WriteString(escaped)

		if n.HardLineBreak() {
			w.WriteString("\\\\")
		} else if n.SoftLineBreak() {
			w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.String)
		escaped := escapeLatexText(string(n.Value))
		w.WriteString(escaped)
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if entering {
		if n.Level == 1 {
			w.WriteString("\\textit{")
		} else {
			w.WriteString("\\textbf{")
		}
	} else {
		w.WriteString("}")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		// For now, just render the link text
		// TODO: proper hyperref support
		w.WriteString("\\href{")
		w.Write(n.Destination)
		w.WriteString("}{")
	} else {
		w.WriteString("}")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Image)
		w.WriteString("\\begin{figure}[h]\n\\centering\n\\includegraphics{")
		w.Write(n.Destination)
		w.WriteString("}\n")

		// Render caption if there's alt text
		if n.FirstChild() != nil {
			w.WriteString("\\caption{")
		}
	} else {
		n := node.(*ast.Image)
		if n.FirstChild() != nil {
			w.WriteString("}")
		}
		w.WriteString("\n\\end{figure}\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\\texttt{")
		// CodeSpan contains text nodes as children, walk them
	} else {
		w.WriteString("}")
	}
	return ast.WalkContinue, nil
}

func (r *LatexRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Skip raw HTML
	return ast.WalkSkipChildren, nil
}

// escapeLatexText escapes special LaTeX characters in text
func escapeLatexText(s string) string {
	// First, decode HTML entities
	s = html.UnescapeString(s)

	// Convert smart quotes and common typographic characters to LaTeX equivalents
	// Do this BEFORE escaping special LaTeX characters
	s = strings.ReplaceAll(s, "\u201C", "``")  // left double quotation mark
	s = strings.ReplaceAll(s, "\u201D", "''")  // right double quotation mark
	s = strings.ReplaceAll(s, "\u2018", "`")   // left single quotation mark
	s = strings.ReplaceAll(s, "\u2019", "'")   // right single quotation mark
	s = strings.ReplaceAll(s, "\u2013", "--")  // en dash
	s = strings.ReplaceAll(s, "\u2014", "---") // em dash

	// Convert ASCII straight quotes to LaTeX quotes
	// This is a simple state-based approach: alternate between opening and closing quotes
	s = convertASCIIQuotes(s)

	// Now escape special LaTeX characters (but preserve our LaTeX quote markers)
	replacer := strings.NewReplacer(
		"\\", "\\textbackslash{}",
		"#", "\\#",
		"$", "\\$",
		"%", "\\%",
		"&", "\\&",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
	)
	return replacer.Replace(s)
}

// convertASCIIQuotes converts ASCII straight quotes (") to proper LaTeX quotes
// Opening quotes become “ and closing quotes become ”
func convertASCIIQuotes(s string) string {
	if !strings.Contains(s, "\"") {
		return s
	}

	var result strings.Builder
	result.Grow(len(s))

	openQuote := true

	for i, r := range s {
		if r == '"' {
			// Determine if this should be an opening or closing quote
			// based on context
			prevChar := rune(0)
			nextChar := rune(0)

			if i > 0 {
				prevChar = rune(s[i-1])
			}
			if i < len(s)-1 {
				nextChar = rune(s[i+1])
			}

			// Heuristic: opening quote typically follows whitespace, punctuation, or start of string
			// Closing quote typically precedes whitespace, punctuation, or end of string
			isOpening := false

			if prevChar == 0 || prevChar == ' ' || prevChar == '\n' || prevChar == '\t' ||
				prevChar == '(' || prevChar == '[' || prevChar == '{' {
				isOpening = true
			} else if nextChar == 0 || nextChar == ' ' || nextChar == '\n' || nextChar == '\t' ||
				nextChar == '.' || nextChar == ',' || nextChar == ';' || nextChar == ':' ||
				nextChar == '!' || nextChar == '?' || nextChar == ')' || nextChar == ']' || nextChar == '}' {
				isOpening = false
			} else {
				// Use the toggle state as fallback
				isOpening = openQuote
			}

			if isOpening {
				result.WriteString("``")
				openQuote = false
			} else {
				result.WriteString("''")
				openQuote = true
			}
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// RenderMarkdownToLatex renders markdown content to LaTeX
func RenderMarkdownToLatex(source []byte, sectionOffset int) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.NewTable(),
			extension.NewFootnote(),
			extension.DefinitionList,
		),
	)
	doc := md.Parser().Parse(text.NewReader(source))

	latexRenderer := NewLatexRenderer(sectionOffset)

	// First pass: collect footnote content
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if fn, ok := n.(*extast.Footnote); ok && entering {
			// Collect footnote content
			var contentBuf bytes.Buffer
			contentWriter := &bufWriter{buf: &contentBuf}

			// Walk children to get content
			for child := fn.FirstChild(); child != nil; child = child.NextSibling() {
				ast.Walk(child, func(cn ast.Node, cEntering bool) (ast.WalkStatus, error) {
					// Render child nodes but skip footnote backlink
					if _, isBacklink := cn.(*extast.FootnoteBacklink); isBacklink {
						return ast.WalkSkipChildren, nil
					}
					return latexRenderer.renderNode(contentWriter, source, cn, cEntering)
				})
			}

			latexRenderer.footnotes[fn.Index] = strings.TrimSpace(contentBuf.String())
		}
		return ast.WalkContinue, nil
	})

	// Second pass: render document
	var buf bytes.Buffer
	writer := &bufWriter{buf: &buf}

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		return latexRenderer.renderNode(writer, source, n, entering)
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// bufWriter implements util.BufWriter using bytes.Buffer
type bufWriter struct {
	buf *bytes.Buffer
}

func (w *bufWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *bufWriter) WriteByte(c byte) error {
	return w.buf.WriteByte(c)
}

func (w *bufWriter) WriteRune(r rune) (int, error) {
	return w.buf.WriteRune(r)
}

func (w *bufWriter) WriteString(s string) (int, error) {
	return w.buf.WriteString(s)
}

func (w *bufWriter) Flush() error {
	return nil
}

func (w *bufWriter) Available() int {
	return 1024 // arbitrary
}

func (w *bufWriter) Buffered() int {
	return w.buf.Len()
}

func (r *LatexRenderer) renderNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	switch n := node.(type) {
	case *ast.Document:
		return r.renderDocument(w, source, n, entering)
	case *ast.Heading:
		return r.renderHeading(w, source, n, entering)
	case *ast.Paragraph:
		return r.renderParagraph(w, source, n, entering)
	case *ast.Blockquote:
		return r.renderBlockquote(w, source, n, entering)
	case *ast.CodeBlock:
		return r.renderCodeBlock(w, source, n, entering)
	case *ast.FencedCodeBlock:
		return r.renderFencedCodeBlock(w, source, n, entering)
	case *ast.HTMLBlock:
		return r.renderHTMLBlock(w, source, n, entering)
	case *ast.List:
		return r.renderList(w, source, n, entering)
	case *ast.ListItem:
		return r.renderListItem(w, source, n, entering)
	case *ast.ThematicBreak:
		return r.renderThematicBreak(w, source, n, entering)
	case *extast.Table:
		return r.renderTable(w, source, n, entering)
	case *extast.TableHeader:
		return r.renderTableHeader(w, source, n, entering)
	case *extast.TableRow:
		return r.renderTableRow(w, source, n, entering)
	case *extast.TableCell:
		return r.renderTableCell(w, source, n, entering)
	case *extast.DefinitionList:
		return r.renderDefinitionList(w, source, n, entering)
	case *extast.DefinitionTerm:
		return r.renderDefinitionTerm(w, source, n, entering)
	case *extast.DefinitionDescription:
		return r.renderDefinitionDescription(w, source, n, entering)
	case *extast.FootnoteLink:
		return r.renderFootnoteLink(w, source, n, entering)
	case *extast.Footnote:
		return r.renderFootnote(w, source, n, entering)
	case *extast.FootnoteList:
		return r.renderFootnoteList(w, source, n, entering)
	case *extast.FootnoteBacklink:
		return r.renderFootnoteBacklink(w, source, n, entering)
	case *ast.Text:
		return r.renderText(w, source, n, entering)
	case *ast.String:
		return r.renderString(w, source, n, entering)
	case *ast.Emphasis:
		return r.renderEmphasis(w, source, n, entering)
	case *ast.Link:
		return r.renderLink(w, source, n, entering)
	case *ast.Image:
		return r.renderImage(w, source, n, entering)
	case *ast.CodeSpan:
		return r.renderCodeSpan(w, source, n, entering)
	case *ast.RawHTML:
		return r.renderRawHTML(w, source, n, entering)
	}
	return ast.WalkContinue, nil
}
