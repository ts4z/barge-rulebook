#
# barge-rulebook-md Makefile
#

ALL=rulebook.pdf
all: $(ALL)

book:
	mdbook build

%.pdf: %.latex
	xelatex $<

rulebook.latex: latexify.pl src/*.md
	./latexify.pl

clean:
	-rm $(ALL) $(PROBLEMATIC) *.dvi rulebook.latex *.log *~ src/*~
	-rm -r book
