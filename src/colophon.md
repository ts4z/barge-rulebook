Colophon
========

The "source code" for this document is here:

[https://github.com/ts4z/barge-rulebook](https://github.com/ts4z/barge-rulebook)

This document is written in Markdown with a tool called
[mdbook](https://rust-lang.github.io/mdBook/) that makes it pretty easy to
navigate.

A GitHub workflow is used to build the rulebook into this version hosted in
GitHub pages here:

[https://ts4z.github.io/barge-rulebook](https://ts4z.github.io/barge-rulebook)

The "dead tree edition" of the rulebook, intended for paper and to be used at a
poker table, is produced via a perl script that wraps CommonMark and LaTeX and
produces PDF. This is available with the rest of the source code above.  It is
a very limited translator but it works for most of the use cases.

The PDF version omits two appendicies that are tables.  Unfortunately, Markdown does
not have a specific way of representing tables nicely, and the CommonMark and mdbook
renderers do not get along.
