Colophon
========

This book is available in one of two ways:

* Web: [https://www.barge.org/rulebook/](https://www.barge.org/rulebook/)
* PDF: [https://www.barge.org/rulebook.pdf](https://www.barge.org/rulebook.pdf)

These are both produced from the same source code available at
[https://github.com/ts4z/barge-rulebook/](https://github.com/ts4z/barge-rulebook/).

This document is primarily written with
[mdbook](https://rust-lang.github.io/mdBook/), which produces a group of nice,
easy-to-navigate web pages.  The text is written in the Markdown format, which
is easy to write and easy to read even if you don't understand the technical
bits.

A [GitHub](https://github.com/) workflow is used to build the rulebook into the
version hosted at [barge.org](https://www.barge.org/).  GitHub (and the
underlying `git` utility) provide our version control.

We also have a PDF version produced by converting the mdbook/Markdown files
into LaTeX.  When we did the 2024 rulebook, we discovered mdbook's PDF support
just wasn't as good as our previous [LaTeX](https://www.latex-project.org/)
version.  So now, we have a small perl script that extracts the file list from
the mdbook configuration, uses CommonMark (a Markdown library) to produce
LaTeX, then doctors that up with some minor hacks and renders it to PDF.

All of the tools are common free software except the perl script and
build glue, which is included in the `git` repository.
