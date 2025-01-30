BARGE Rulebook Redux
====================

**VIEW THE RULEBOOK** [here](https://ts4z.github.io/barge-rulebook/).

**GET NICE PDF** [here](https://ts4z.github.io/barge-rulebook.pdf).

This is a revision of the BARGE rulebook intended for a multi-author audience.

This produces readable input text, and HTML, and a proper print-to-book PDF
from one set of input files that (hopefully) anyone can edit.

Almost all of the book text is in Markdown, a popular easy-to-edit format that
is easily confused for "old fashioned plain text".  It reads as text, so the
input format is readable and looks great in email; it renders to HTML well and
supports links; and it can be turned into LaTeX, which renders well on a
printed page.

In order to produce a very web-friendly version, we use md
[`mdbook`](https://rust-lang.github.io/mdBook/).  In order to produce a
print-friendly version, we use a pipeline that converts Markdown to LaTeX and
applies some small fixes.

Building the mdbook Rulebook
----------------------------

Install `mdbook`, then run `mdbook build`.

mdbook is easy to install with the rust `cargo` utility.  If you are running on
something like Debian Linux, you'll need a recent Rust to build `mdbook`.

### Doing the above automatically

This book is automatically built by a GitHub action and published to
[ts4z.github.io/barge-rulebook/](https://ts4z.github.io/barge-rulebook/).

If this stops automatically working, it is likely that a GitHub personal access
token has expired and must be updated.  This token belongs to the GitHub user
`ts4z` and is stored as a secret in the barge-rulebook repository.

This is sufficiently robust (and Markdown is sufficiently forgiving as a
format) that I haven't had problems just updating and assuming the right thing
will happen.

Building LaTeX Rulebook
-----------------------

Math majors everywhere rejoice: we have a way of producing a LaTeX rulebook,
with nice page numbers, that should look great printed.  There are a few
problems with this:

- There are some bugs with the translation.
- XeTeX is required, since the input documents use Unicode.
- Perl is required, and cpan-type CommonMark must be available.
- The perl script that produces the LaTeX input is a total hack.  (I could say
  the same about CommonMark.)
- Tables in CommonMark don't get turned into tables in LaTeX.
  These get completely mangled.  So we have alternate versions of those pages.

All of that said, if you're set up with the right packages, `make rulebook.pdf`
will do it.

If you are on a Debian Linux system, you will need these packages:

- libcommonmark-perl
- make
- perl
- texlive-fonts-extra
- texlive-xetex

I believe this list is complete.  It is possible to wrap these up into a Docker
container (see the Dockerfile in this directory); however the container is
prohibitively large for the free GitHub runners, so I build the LaTeX rulebook
manually.

Style
-----

Please match the style of the document as a whole.  Generally, games have the
same format from game to game, but for some games this is quite redundant,
so they just refer to the original game.

Footnotes and tables are not permitted in Markdown source as they can't be
reliably rendered by CommonMark in the LaTeX version.  Tables are permissible
only if two versions of the source file is provided, one in GFM and one in
LaTeX; the file should contain nothing other than the table.

Some of the strucutre of the content can be improved.  Don't be afraid to make
it better.

All titles in SUMMARY.md should match the title at the top of the page.  (Pages
that exist to hold alternate versions of tables are excepted.)  Note that
titles are used only by mdbook, but both mdbook and pdf versions use SUMMARY.md
to determine the chapter order.

Do not have two versions of a file (X.md and X.latex) unless X.md has data that
can't be translated to LaTeX.  This probably means a table.  Keeping the two in
sync is unlikely.

Do not depend on implementation details.  For instance, someday, there will
only be one version of all input files.

### Markdown notes

Markdown comes in several flavors.  Basic Markdown lacks support for tables,
which particularly hurts.  Many technical users support a variant called 
"GitHub Flavored Markdown", which we can partially use.

Avoid footnotes (unless latexify.pl has been modified to be able to process
them).

Technical Aspects
-----------------

The rulebook is in the same order in both the LaTeX and mdbook versions.  This
order is declared in src/SUMMARY.md, whose format is required by `mdbook` and
repurposed as an input for `latexify.pl` for the PDF version.  Logically, if
you catenate the files in SUMMARY.md together, you get a nice Markdown
document, albeit without a Table of Contents.

- All differences betwen src/X.md and src/X.latex represent only practical
  formatting differences intended to produce similar (not identical) products
  in both contexts.
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
- Unicode characters are permitted, but note LaTeX is limited in which ones it
  knows about.  Some of the common ones require work in `setup.latex`.

### Limitations

A multi-output format tends to work to the lowest common denominator, and
that's the case here.  The LaTeX is limited by the fact that it's translated
from Markdown (aside from some boilerplate), and we can't use advanced (GFM)
Markdown features because CommonMark, our md->LaTeX translator, doesn't
understand them.

This mostly affects the PDF version, but the mdbook version is missing some
nice features, too (like color suits).

Future Work
-----------

Make all of the people happy all of the time.

I'd love to find a way to get rid of the duplicated tables.

For the PDF version, `latexify.pl` works, but it's clearly a hack.  It would be
nice to have more control over the structure of the output and remove all of
the simple text substitutions.

Some things, like the Chowaha board, could change from text to images.
(Both mdbook and LaTeX can render images.)

Replace CommonMark with GoldMark, which supports tables and footnotes (I
think).  Probably means rewriting latexify.pl in Go.

Improve a few images, like the logo on the title page, to increase size.

Consider changing margins as to save paper (as per a discussion with Patrick
Milligan).
