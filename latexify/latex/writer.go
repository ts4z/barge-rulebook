package latex

import (
	"fmt"
	"io"
)

// SectionAwareWriter implements LatexWriter for writing to a file
type SectionAwareWriter struct {
	writer io.Writer
	level  int // current section nesting level
	parent *SectionAwareWriter
}

// NewSectionAwareWriter creates a new LatexWriter that writes to the given writer
func NewSectionAwareWriter(w io.Writer) *SectionAwareWriter {
	return &SectionAwareWriter{
		writer: w,
		level:  0,
		parent: nil,
	}
}

// WriteText writes raw text/LaTeX to the output
func (w *SectionAwareWriter) WriteText(s string) {
	fmt.Fprint(w.writer, s)
}

// WriteComment writes a LaTeX comment
func (w *SectionAwareWriter) WriteComment(comment string) {
	fmt.Fprintf(w.writer, "%%%% %s\n", comment)
}

// BeginSubsection creates a nested writer for a subsection
func (w *SectionAwareWriter) BeginSubsection() *SectionAwareWriter {
	return &SectionAwareWriter{
		writer: w.writer,
		level:  w.level + 1,
		parent: w,
	}
}

// Close ends the current subsection (no-op for this implementation)
func (w *SectionAwareWriter) Close() {
	// In this implementation, we don't need to do anything special
	// when closing a subsection, but we keep this for interface compatibility
}

// WriteChapter writes a chapter heading
func (w *SectionAwareWriter) WriteChapter(title string) {
	fmt.Fprintf(w.writer, "\\chapter{%s}\n", EscapeSpecials(title))
}

// WriteSection writes a section heading
func (w *SectionAwareWriter) WriteSection(title string) {
	fmt.Fprintf(w.writer, "\\section{%s}\n", EscapeSpecials(title))
}

// WriteSubsection writes a subsection heading
func (w *SectionAwareWriter) WriteSubsection(title string) {
	fmt.Fprintf(w.writer, "\\subsection{%s}\n", EscapeSpecials(title))
}

// WriteSectionByLevel writes a section at the appropriate level
func (w *SectionAwareWriter) WriteSectionByLevel(level int, title string) {
	switch level {
	case 0:
		w.WriteChapter(title)
	case 1:
		w.WriteSection(title)
	case 2:
		w.WriteSubsection(title)
	default:
		// For deeper nesting, use subsubsection
		fmt.Fprintf(w.writer, "\\subsubsection{%s}\n", EscapeSpecials(title))
	}
}
