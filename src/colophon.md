Colophon
========

This book is available in one of three ways:

* Web: [https://ts4z.github.io/barge-rulebook](https://ts4z.github.io/barge-rulebook)
* PDF: (location TBD)
* Source: [https://github.com/ts4z/barge-rulebook](https://github.com/ts4z/barge-rulebook)

This document is written in Markdown, with a tool called
[mdbook](https://rust-lang.github.io/mdBook/) that makes it pretty easy to
navigate.  mdbook produces a pretty good printed copy.

A GitHub workflow is used to build the rulebook into this version hosted in
GitHub pages (at the web link above).

The "dead tree PDF edition" of the rulebook, intended for paper and to be used
at a poker table, is produced via a perl script that extracts the file list
from the mdbook configuration, uses CommonMark (a Markdown library) to produce
LaTeX, then doctors that up with a perl script.  All of the pieces are either
free software or included at with the rest of the source code above.
