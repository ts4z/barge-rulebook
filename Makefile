#
# barge-rulebook-md Makefile
#
# This is not a good Makefile, but we have a simple project.
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
	src/game-archie.md \
	src/game-badacey.md \
	src/game-badeucy.md \
	src/game-badugi.md \
	src/game-bidirectional-chowaha-high-low.md \
	src/game-big-o-five-card-omaha-high-low-eight-or-better.md \
	src/game-binglaha.md \
	src/game-california-lowball.md \
	src/game-chowaha.md \
	src/game-courchevel-high-only.md \
	src/game-crazy-pineapple-high-low-eight-or-better.md \
	src/game-deuce-to-seven-lowball.md \
	src/game-deuce-to-seven-razz.md \
	src/game-deuce-to-seven-triple-draw.md \
	src/game-dramadugi.md \
	src/game-dramaha.md \
	src/game-dramaha-49.md \
	src/game-duck-flush.md \
	src/game-five-card-draw.md \
	src/game-five-card-stud.md \
	src/game-four-card-chowaha-high-low-eight-or-better.md \
	src/game-holdem.md \
	src/game-irish.md \
	src/game-korean.md \
	src/game-lazy-pineapple-high-only.md \
	src/game-lazy-pineapple-high-low-eight-or-better.md \
	src/game-london-lowball.md \
	src/game-mexican-poker.md \
	src/game-mississippi-stud-stud-high-low-eight-or-better-razz.md \
	src/game-murder.md \
	src/game-oklahoma.md \
	src/game-omaha-high-only.md \
	src/game-omaha-high-low-eight-or-better.md \
	src/game-omaha-x-or-better.md \
	src/game-paradise-road-pickem.md \
	src/game-razz.md \
	src/game-razzdugi.md \
	src/game-rio-bravo.md \
	src/game-scrotum.md \
	src/game-short-deck-omaha.md \
	src/game-short-deck-texas-holdem.md \
	src/game-sohe-simultaneous-omaha-holdem.md \
	src/game-stud.md \
	src/game-stud-high-low-eight-or-better.md \
	src/game-stud-high-low-no-qualifier.md \
	src/game-super-stud-stud-high-low-eight-or-better-razz.md \
	src/game-texas-holdem.md \
	src/game-texas-holdem-high-low-eight-or-better.md \
	src/game-triple-draw-dramaha.md \
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
