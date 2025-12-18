# Notes on Why this is the Way it Is

Tim Showalter / 17 Dec 2025

# History

BARGE had a website, and it had a list of poker variants on it.  That was
enough of a rulebook.

Around 2015, BARGE had a 25-game mixed tournament, and in 2016, we had a
26-game tournament.  The poker room needed actual written down rules for these
games.  This was, I think, the motivation for Chris Mecklin to make the first
rulebook.  Chris has a math PhD, so naturally his authoring tool of choice is
LaTeX.  After that BARGE, he continued to maintain the rulebook.

The 2015/2016 rulebook was a reference manual, but it supported hyperlinks
within the PDF, and can be indexed by search engines that understand PDF, but
you can't link directly to a game within the PDF.  It also was unwieldy to use
on a phone.

Chris stopped updating the rulebook after BARGE 2021.  
The rulebook wasn't updated to sort out date references or verb tenses.  New
games, like [Omaha
X-or-better](https://www.barge.org/rulebook/omaha-x-or-better.html), couldn't
be added to the rulebook without re-making the source for the rulebook.

In 2024, I was playing on [https://craftpoker.com/](craftpoker.com) somewhat
regularly, and the game Archie got picked as part of a lot of mixed game lists.
I like Archie, but I keep forgetting how to play it.  I found the PDF rulebook
hard to use.

With some encouragement from Rich, I took a pass at recreating the rulebook.  I
had used Markdown very successfully for long-lived high-content tech
documentation, and that's what I wanted to use.

My goals were these:

1. Rulebook on barge.org.  Good HTML version.  A community resource is
   convenient and authoritative.  Providing a mixed game reference is good
   advertising.
2. Allow for multiple authors.  Enable other people to make corrections and
   additions.
3. Good paper version.  Paper is awesome, but it doesn't work well if we can't
   enforce organization.
4. Avoid future ground-up rewrites.  BARGE has been around a long time.  The
   rulebook should be durable.  We should stick to technologies that we are
   sure will still work in 10 years.
5. Don't allow anyone to be the "only" maintainer.  Plan to flake out and quit
   without warning, even if it is not my intent.

Markdown has a few things going for it.  It's easy to edit, since it mostly
reads like plain text.  I hoped less technical people would contribute.
It is easy to turn into very readable HTML.  It's easy to search, and if you
check it in to version control, you can easily see a complete edit history.

Markdown also has its limitations.  It is poorly versioned and routinely (and
silently) extended by most of its major users.  It has no built-in support for
footnotes or tables.

I considered LaTeX and raw HTML.

I didn't use HTML because it's something of an attractive nuisance for people
to edit it with tools.  I wanted to restrict people (including myself) to
mostly plain text.

Similarly, I didn't want to use LaTeX because I find it to be something of a
major impediment for non-technical type folks.  I don't want to be only one who
edits the rulebook, and I *really* don't want to hand off a LaTeX environment
off should I get hit by a bus.

My hope was to be able to use a Markdown parser that could also output LaTeX.
This would allow us to recreate something like the 2021 rulebook.

I cut-and-pasted the 2021 rulebook into an Emacs buffer and turned it into
Markdown.  (This is easy, but tedious.)  The initial version rendered as one
big web page, and started to write a tool to generate the table of contents.
Cliff Matthews recommended I use `mdbook` instead, which provides a better web
viewer.  I was a little hesitant, mostly because it meant a major refactor, but
the quality was so much better I got on with it.

This solved a lot of my problems.  The book was web-navigable.  It looks good
on a phone.  I had a Table of Contents.  I could link directly to the Archie
rules.

This was the initial new 2024 rulebook.  It was hosted on ts4z.github.io,
because hosting dozens of files wasn't something our previous web host
environment could reasonably do.  But, for BARGE 2024, at least the rulebook
was somewhat up to date.

The only negative feedback was that the printed version was not good.  In the
fall, I decided to do this.  I wrote some Perl to output LaTeX, using the
CommonMark library.  CommonMark has LaTeX output, so I could use its
translation.  But CommonMark doesn't understand tables, or footnotes, both of
which I had used in the text.  I made my own workarounds: I wrote LaTeX
versions of the big tables and subbed them in, and I removed the footnotes
(there weren't many).  I worked around a couple other bugs.

Then I figured out how to add colored suit symbols to the LaTeX version, a
popular feature of the 2021 rulebook.  (I didn't colorize the rank, which I
prefer, but it is also a little harder to do.)

The BARGE website was also in need of a refresh, so that happened in parallel.
Ultimately, we cut over to the "new" website when we saw a hiccup in the
hosting provider.  The new website can be edited by more than just one person,
has version control, the next events' dates on nearly every page, and now, a
new home for the rulebook.

In December of 2025, I figured out how to get colored suit symbols in the
`mdbook` version.  This is something of a hack, but it is working well enough for
now.

Also in December of 2025, I replaced my LaTeX-generating Perl script with
something written in Go, using GoldMark, that understands tables and footnotes.
(LLMs are neat.)  This allowed fixing some of the other minor glitches in the
LaTeX generation at the expense of having to maintain a somewhat larger
utility.

## Hacks

### `latexify`

Our utility for turning `mdbook` into LaTeX is called `latexify` and it lives
in the `latexify` directory.  It demands a recent version of Go but it will
probably work with an older one.

The only user of `latexify` is this book, so while the approach is general,
it is not well-tested or particularly universal.

But it works.

`latexify` ended up learning how to fix section numbering.  `mdbook` works best
when each page starts with an `h1`, but LaTeX doesn't want each subsection to
start with a `title` equivalent.  Having a custom translator allowed me to fix
that up.  (There are probably bugs here.)  This actually improved the `mdbook`
version because the left-hand section index no longer has a duplicated title
for each game.

`latexify` knows about some special files that are used at various points in
the output.  One file sets up the LaTeX environment (setup.pl), then there are
files for declaring the frontmatter, mainmatter, and backmatter, as well as
being able to entirely override a Markdown file if a LaTeX version is
available.

#### Capitalizing LaTeX correctly

One of the quirks of TeX documents is that they inevitably spell TeX (and
LaTeX) correctly because there is a command for that facility.  Well, we do it
by simple string substitution, replacing "LaTeX" with "\\LaTeX{}".  Until the
word "latex" appears in a poker book for any other reason, this hack will
probably continue to exist.

#### Maintenance burden

`latexify` represents a moderate amount of code that can't just be downloaded.
This isn't desirable, and if a better path to a book becomes available, we
should switch.  Because this code was originally LLM-authored, it should not be
that hard to recreate, if we have to.

### ♠♥♣♦

Colored suit symbols exist in both the LaTeX version and the `mdbook` version.

In the LaTeX version, there is a small amount of TeX code that replaces all of
these characters with colorized versions, as well as replacing ✕ with a
mathematical version that TeX doesn't mind.

In the `mdbook` version, we use a string substitution to inject a `span` around
suit symbols which declares a CSS class, and we attach a color to the CSS
class.

### HTML entities and left/right quotes

The LaTeX translator is a little loosey-goosey and not carefully thought out.
Mostly, this is LLM-authored code and it appers to be correct for all of the
cases that we care about.  This may require more care.

This was a particular issue for "smart" quotes, but they seem to be working
now.

### Tables

The 2021 rulebook had a large table for lowball hand numbers.  Initially, we
built the LaTeX version by providing an "alternate version" of this file that
could be inlined into the LaTeX, and tables could not appear in most Markdown
at all.  Now, tables are supported in Markdown, but we have preserved the LaTeX
tables because they look better, and I added a couple more that look *much*
better rotated.

### Footnotes

I like footnotes in text.  GoldMark supports them, as does `latexify`.  These
are used sparingly but are very useful for references to source material.

## CI/CD

My GitHub account hosts the repository.  A pipeline does ``mdbook` build` and
pushes all changes to the barge.org repository, also hosted in my GitHub
account.

A Makefile exists in the repo to facilitate this, but it is more used for
building and deploying the PDF version.  The Makefile can build the PDF,
check out barge.org, and update rulebook.pdf there.

Unfortunately, I don't think this can be easily run on GitHub, at least without
paying GitHub for the privilege.  We're well above the free tier for the
full-fat TeX installation that's required to generate our book for our naive
install.  A real TeXpert could boil this down to something that could run on
GitHub, but geting XeTeX to run in that environment wasn't easy, or fun, and I
don't want to leave a maintenance burden for later.

Additionally, to run LaTeX, one must first produce the input file, which requires
running a Go program; so we'd have to *also* install Go.

For now, I'll do this part outside CI.
