Colophon
========

This book is available in one of three ways:

* Web: [https://www.barge.org/rulebook/](https://www.barge.org/rulebook/)
* PDF: [https://www.barge.org/rulebook.pdf](https://www.barge.org/rulebook.pdf)
* Source: [https://github.com/ts4z/barge-rulebook](https://github.com/ts4z/barge-rulebook)

This document is written in Markdown, with a tool called
[mdbook](https://rust-lang.github.io/mdBook/) that makes it pretty easy to
navigate.

A GitHub workflow is used to build the rulebook into this version hosted in
GitHub pages (at the web link above), and another one publishes it to
[barge.org](https://www.barge.org/).

mdbook has print support in conjunction with a web browswer, but it isn't as
usable at the poker table as the old LaTeX version, which produces a really
nice book with a proper Table of Contents and page numbers.

So we also have the "dead tree PDF edition" of the rulebook, intended for paper
and to be used at a poker table. This is produced via a perl script that
extracts the file list from the mdbook configuration, uses CommonMark (a
Markdown library) to produce LaTeX, then doctors that up with a perl script.

All of the tools are either common free software or included with the other
pieces above.
