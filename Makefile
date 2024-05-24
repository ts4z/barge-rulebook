#
# This is not a good Makefile, but we have a simple project.
#
# Note that the LaTeX rules don't really work at all.  I haven't figured out
# how the Markdown->LaTeX tools are intended to work and haven't been satisied
# with the results.  This means printed copies won't be very good (no index).
#

ALL=rulebook.html
PROBLEMATIC=rulebook.pdf

SUBSTANCE=frontmatter.md main.md
BOILERPLATE=front-boilerplate.latex end-boilerplate.latex

TITLE=Barge Rulebook 2024

all: $(ALL)

# this is, unfortunately, reliant on a utility of my own design.
# see https://github.com/ts4z/markdown-toc
rulebook.html: frontmatter.md main.md
	markdown-toc -o $@ -f frontmatter.md -i main.md -t "$(TITLE)"

# utter crap
rulebook.latex: $(BOILERPLATE) $(SUBSTANCE)
	-rm $@
	cat front-boilerplate.latex >> $@ || (rm $@ && false)
	(cat frontmatter.md main.md) | cmark -t latex >> $@ || (rm $@ && false)
	cat end-boilerplate.latex >> $@ || (rm $@ && false)

%.dvi : %.latex
	latex $^

%.pdf : %.dvi
	dvipdf $^ $@

clean:
	-rm $(ALL) $(PROBLEMATIC) *.dvi rulebook.latex *.log *~
