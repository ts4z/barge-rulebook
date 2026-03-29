# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build the PDF (runs latexify then xelatex twice)
make

# Generate LaTeX source only (without compiling to PDF)
make rulebook.latex

# Build the HTML version
mdbook build

# Rebuild the latexify Go tool after editing latexify/*.go
(cd latexify && go build .)

# Clean all build artifacts
make clean

# Publish PDF to barge.org (requires git access)
make update-rulebook-on-barge-site
```

The PDF build runs xelatex twice intentionally — `longtable` requires two passes.

Use `make` to access build tools and to document the build process.

## Architecture

This is a poker rulebook with two output formats: HTML (via mdbook) and PDF (via a custom LaTeX pipeline).

### Source Files

- `src/SUMMARY.md` — defines the book structure and file ordering; drives both the HTML and PDF builds
- `src/*.md` — rulebook content in Markdown
- `src/*.latex` — raw LaTeX fragments that override the markdown-to-LaTeX conversion for specific pages, or are inserted at structural boundaries

### PDF Pipeline

```
src/*.md + src/*.latex
        ↓
latexify/latexify  (Go program)
        ↓
rulebook.latex
        ↓
xelatex (twice)
        ↓
rulebook.pdf
```

The `latexify` program (in `latexify/`) parses `SUMMARY.md` to understand document structure, then converts each Markdown file to LaTeX using the goldmark library. If a `.latex` file exists alongside a `.md` file (e.g., `src/foo.latex`), the raw LaTeX is used instead of converting the Markdown.

Special LaTeX files are injected at section boundaries in this order: `setup`, `frontmatter`, `mainmatter`, `backmatter`, `endmatter` (all in `src/`). `setup.latex` contains the `\documentclass` declaration and preamble. The suit symbols (♠♥♦♣) are handled differently in each output: HTML uses a `mdbook-replace` preprocessor with CSS spans; LaTeX uses `\newunicodechar` with color macros.

### HTML Pipeline

`mdbook build` uses the `mdbook-replace` preprocessor (configured in `book.toml`) to replace suit symbols with colored HTML `<span>` elements, then renders to `book/`.

### latexify internals

- `latexify/main.go` — entry point; orchestrates file rendering and special file injection
- `latexify/mdbook/summary.go` — parses `SUMMARY.md` into a tree of `SummaryNode`s
- `latexify/latex/renderer.go` — goldmark AST → LaTeX renderer
- `latexify/latex/writer.go` — wraps output with section-aware helpers
