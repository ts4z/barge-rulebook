#
# This is not a good Makefile, but we have a simple project.
#
# Note that the LaTeX rules don't really work at all.  I haven't figured out
# how the Markdown->LaTeX tools are intended to work and haven't been satisied
# with the results.  This means printed copies won't be very good (no index).
#

ALL=rulebook.html

# Order here is important.
SECTIONS=\
	src/preface.md \
	src/common.md \
	src/variants.md \
	src/game-ace-to-5-triple-draw.md \
	src/game-action-razz.md \
	src/game-action-razzdugi.md \
	src/games.md \
	src/game-two-or-five-omaha-8-or-better.md \
	src/appendix-a.md \
	src/appendix-b.md \
	src/appendix-c.md \
	src/colophon.md
SUBSTANCE=src/frontmatter.md $(SECTIONS)
BOILERPLATE=front-boilerplate.latex end-boilerplate.latex

TITLE=Barge Rulebook 2024

all: $(ALL)

# this is, unfortunately, reliant on a utility of my own design.
# see https://github.com/ts4z/markdown-toc
rulebook.html: $(SUBSTANCE)
	markdown-toc -o $@ -t "$(TITLE)" -f src/frontmatter.md $(SECTIONS)

clean:
	-rm $(ALL) $(PROBLEMATIC) *.dvi rulebook.latex *.log *~
