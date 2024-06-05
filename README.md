BARGE Rulebook Redux
====================

This is my revision of the BARGE rulebook intended for a multi-author audience.

This is an attempt to produce readable input text and HTML (and someday, PDF)
from one set of documents that is easy to maintain.

Markdown is an easy-to-edit but subtly ambiguous format that has become
popular.  This format has a few things going for it: it reads as text, so the
input format is readable and looks great in email; it renders to HTML well and
supports links; and it can be turned into LaTeX, which renders well on a
printed page.

In order to make this acceptable in terms of text representation, I've used
Github-Flavored Markdown, which allows for tables, footnotes, and other
niceities.  Fortunately, these seem to work very well with
[`mdbook`](https://rust-lang.github.io/mdBook/).

TBD: A nice PDF version with functional page numbers should be doable, but I
haven't figured out the easiest path for it yet.  (A minor complication is that
suit characters, and some other Unicode things, appear in the text.  Some TeX
derivatives aren't believers in Unicode yet, so some of the tools don't "just
work").

Building the Rulebook
---------------------

Install `mdbook`.  `mdbook build` will do it.

Or, `make` should do it, too.

Future Work
-----------

Nice printed version: Unsolved.  Markdown can easily be translated to LaTeX.
There are several tools that do this well.  Unfortunately I have chosen to use
Unicode instead of things that predate Unicode.  I'm hoping a hyperlinked
document is more useful than page numbers.

Markdown notes
--------------

Markdown is deeply ambiguous.  This tries to stick to GitHub-flavored Markdown
which allows for footnotes, tables, and triple-quote "code" blocks that are 
occasionally useful for rudimentary diagrams.

### mdbook

`mdbook` is required.  This is available if Rust is installed.  If you are
running on something like Debian Linux, you'll need a recent Rust to build
`mdbook`.
