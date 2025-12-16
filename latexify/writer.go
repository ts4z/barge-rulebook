package main

import (
	"fmt"
	"io"
)

// SectionWriter implements LatexWriter for writing to a file
type SectionWriter struct {
	writer io.Writer
	level  int // current section nesting level
	parent *SectionWriter
}

// NewLatexWriter creates a new LatexWriter that writes to the given writer
func NewLatexWriter(w io.Writer) *SectionWriter {
	return &SectionWriter{
		writer: w,
		level:  0,
		parent: nil,
	}
}

// WriteText writes raw text/LaTeX to the output
func (w *SectionWriter) WriteText(s string) {
	fmt.Fprint(w.writer, s)
}

// WriteComment writes a LaTeX comment
func (w *SectionWriter) WriteComment(comment string) {
	fmt.Fprintf(w.writer, "%%%% %s\n", comment)
}

// BeginSubsection creates a nested writer for a subsection
func (w *SectionWriter) BeginSubsection() *SectionWriter {
	return &SectionWriter{
		writer: w.writer,
		level:  w.level + 1,
		parent: w,
	}
}

// Close ends the current subsection (no-op for this implementation)
func (w *SectionWriter) Close() {
	// In this implementation, we don't need to do anything special
	// when closing a subsection, but we keep this for interface compatibility
}

// WriteChapter writes a chapter heading
func (w *SectionWriter) WriteChapter(title string) {
	fmt.Fprintf(w.writer, "\\chapter{%s}\n", escapeLatex(title))
}

// WriteSection writes a section heading
func (w *SectionWriter) WriteSection(title string) {
	fmt.Fprintf(w.writer, "\\section{%s}\n", escapeLatex(title))
}

// WriteSubsection writes a subsection heading
func (w *SectionWriter) WriteSubsection(title string) {
	fmt.Fprintf(w.writer, "\\subsection{%s}\n", escapeLatex(title))
}

// WriteSectionByLevel writes a section at the appropriate level
func (w *SectionWriter) WriteSectionByLevel(level int, title string) {
	switch level {
	case 0:
		w.WriteChapter(title)
	case 1:
		w.WriteSection(title)
	case 2:
		w.WriteSubsection(title)
	default:
		// For deeper nesting, use subsubsection
		fmt.Fprintf(w.writer, "\\subsubsection{%s}\n", escapeLatex(title))
	}
}

// escapeLatex escapes special LaTeX characters
func escapeLatex(s string) string {
	// Basic escaping - can be expanded as needed
	replacements := map[rune]string{
		'\\': `\textbackslash{}`,
		'#':  `\#`,
		'$':  `\$`,
		'%':  `\%`,
		'&':  `\&`,
		'_':  `\_`,
		'{':  `\{`,
		'}':  `\}`,
		'~':  `\textasciitilde{}`,
		'^':  `\textasciicircum{}`,
	}

	result := ""
	for _, r := range s {
		if repl, ok := replacements[r]; ok {
			result += repl
		} else {
			result += string(r)
		}
	}
	return result
}
