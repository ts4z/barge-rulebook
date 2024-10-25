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

Building the Rulebook (mdbook)
------------------------------

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
- Perl is required, and cpan-type CommonMark must be available.
- The perl script that produces the LaTeX input is a total hack.  (I could say
  the same about CommonMark.)
- Tables in CommonMark don't get turned into tables in LaTeX.
  These get completely mangled.  So we have alternate versions of those pages.

All of that said, if you're set up with the right packages, `make rulebook.pdf`
will do it.

If you are on a Debian (or Ubuntu) Linux system, try installing these packages:

- libcommonmark-perl
- texlive-xetex
- texlive-latex-base
- make
- perl

If that works, it's news to me.  I haven't verified that this list is complete.
Please let me know.

Markdown notes
--------------

Markdown is deeply ambiguous.  This tries to stick to GitHub-flavored Markdown
which allows for footnotes, tables, and triple-quote "code" blocks that are 
occasionally useful for rudimentary diagrams.

### mdbook

`mdbook` is required.  This is available if Rust is installed.  If you are
running on something like Debian Linux, you'll need a recent Rust to build
`mdbook`.

Style
-----

Please match the style of the document as a whole.  Generally, games have the
same format from game to game, but for some games this is quite redundant,
so they just refer to the original game.

Footnotes and tables are not permitted in Markdown source as they can't be
reliably rendered by CommonMark in the LaTeX version.

Some of the strucutre of the content can be improved.  Don't be afraid to make
it better.

All titles in SUMMARY.md should match the title at the top of the page.
(Pages that exist to hold alternate versions of tables are excepted.)

Technical Aspects
-----------------

The rulebook is in the same order in both the LaTeX and mdbook versions.  This
order is declared in src/SUMMARY.md which is parsed directly by mdbook.
Logically, if you catenate the files in SUMMARY.md together, you get a nice
Markdown document.

- All differences betwen src/X.md and src/X.latex represent only practical
  formatting differences intended to produce similar (not identical) products in both contexts.
- If src/X.latex exists and src/X.md does not, the file is an interstitial, and
  its position is controlled by latexify.pl.
  - Interstitials are used to inject the Table of Contents, control chapter
    numbering, and inject LaTeX boilerplate.
- If src/X.md exists and src/X.latex does not exist, X.md will be rendered as
  LaTeX when building the LaTeX book.
- If src/X.md and src/X.latex both exist, X.md is ignored when building the
  LaTeX version.  This is intended to facilitate tables present in the
  appendix.  (If we come up with a better way of rendering tables directly from
  Markdown, we could stop this.)
- Unicode characters are permitted, but note LaTeX is limited in which ones it knows about.
  Some of the common ones require work in `setup.latex`.

Future Work
-----------

Make all of the people happy all of the time.

I'd love to find a way to get rid of the duplicated tables.
