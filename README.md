BARGE Rulebook Redux
====================

This is my revision of the BARGE rulebook intended for a multi-author audience.

This is an attempt to produce readable input text, HTML, and PDF output from
one set of documents.

Markdown is an easy-to-edit but subtly ambiguous format that has become
popular.  This format has a few things going for it: it reads as text, so the
input format is readable and looks great in email; it renders to HTML well and
supports links; and it can be turned into LaTeX, which renders well on a
printed page.

In order to make this acceptable in terms of text representation, I've used
Github-Flavored Markdown, which allows for tables, footnotes, and other
niceities.

For the HTML version, a program needs to exisst that provides a usable table of
contents.  I wrote one that does a marginally acceptable job.

TBD: The LaTeX conversion is not easy, for reasons I don't understand.

Building the Rulebook
---------------------

HTML: `make rulebook.html` should do it.

LaTeX: Unsolved.  Markdown can easily be translated to LaTeX, but I haven't
found a tool that does the right work to make it into a "book" with a table of
contents.  This is not a huge amount of work, but it hasn't been required yet
(and I'm hoping a hyperlinked document is more useful).

Markdown notes
--------------

Markdown is deeply ambiguous.  This tries to stick to GitHub-flavored Markdown
which allows for footnotes, tables, and triple-quote "code" blocks that are 
occasionally useful for rudimentary diagrams.

### markdown-toc

`markdown-toc` is required.  This is available on GitHub as of this writing
[here](http://github.com/ts4z/markdown-toc).  This is a thin wrapper around
Goldmark, a Go-based Markdown parser that supports a table-of-contents
extension.  `markdown-toc` allows for "front matter" to be prepended and a
linked TOC to be built out of the Markdown.  This rulebook is the reason the
tool exists, so it's possible we could add other features as needed.

`markdown-toc` is set to generate internal anchors (`id` attributes) on the
HTML nodes that can be linked to with a hash anchor in the URL.  These are
based on the section headings, and can be obtained from the table of contents
section.
