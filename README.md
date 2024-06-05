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

`make` should do it for now.

Nice printed version: Unsolved.  Markdown can easily be translated to LaTeX.
There are several tools that do this well.  Unfortunately I have chosen to use
Unicode instead of things that predate Unicode.  I'm hoping a hyperlinked
document is more useful than page numbers.

Future Work
-----------

Right now the rulebook is two files, one that goes before the Table of Contents
and one with all of the other text.  This isn't ideal.  mdbook's approach is
better, but I haven't done all the work to switch to it yet.

Markdown notes
--------------

Markdown is deeply ambiguous.  This tries to stick to GitHub-flavored Markdown
which allows for footnotes, tables, and triple-quote "code" blocks that are 
occasionally useful for rudimentary diagrams.

### mdbook

`mdbook` is required.  This is available if Rust is installed.  If you are
running on something like Debian Linux, you'll need a recent Rust to build
`mdbook`.
