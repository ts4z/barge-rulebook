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

### Doing the above automatically

This book is automatically built by a GitHub action and published to
[ts4z.github.io/barge-rulebook/](https://ts4z.github.io/barge-rulebook/).

If this stops automatically working, it is likely that a GitHub personal access
token has expired and must be updated.  This token belongs to the GitHub user
`ts4z` and is stored as a secret in the barge-rulebook repository.

Building LaTeX Rulebook
-----------------------

Math majors everywhere rejoice: we have a way of producing a LaTeX rulebook,
with nice page numbers.  There are a few problems with this:

- There are some bugs with the translation.
- XeTeX is required, since the input documents use Unicode.
- Perl is required, and CommonMark must be installed (under Debian, install
  `libcommonmark-perl`).
- The perl script that produces the LaTeX input is a total hack.  (I could say
  the same about CommonMark.)
- Tables in CommonMark don't get turned into tables in LaTeX.
  These get completely mangled.


Markdown notes
--------------

Markdown is deeply ambiguous.  This tries to stick to GitHub-flavored Markdown
which allows for footnotes, tables, and triple-quote "code" blocks that are 
occasionally useful for rudimentary diagrams.

### mdbook

`mdbook` is required.  This is available if Rust is installed.  If you are
running on something like Debian Linux, you'll need a recent Rust to build
`mdbook`.
