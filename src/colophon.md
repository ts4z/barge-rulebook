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

mdbook has print support in conjunction with a web browswer, but it isn't as
usable at the poker table as the old LaTeX version, which produces a really
nice book with a proper Table of Contents and page numbers.

So we also have the "dead tree PDF edition" of the rulebook, intended for paper
and to be used at a poker table. This is produced via a perl script that
extracts the file list from the mdbook configuration, uses CommonMark (a
Markdown library) to produce [LaTeX](https://www.latex-project.org/), then
doctors that up with some minor hacks.

All of the tools are either common free software except the perl script and
build glue, which is included in the `git` repository.
