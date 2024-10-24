#
# barge-rulebook-md Makefile
#

ALL=rulebook.pdf
all: $(ALL)

book: src/*.md
	mdbook build

%.pdf: %.latex
	xelatex $< || rm $@

rulebook.latex: latexify.pl src/*.md src/*.latex
	./latexify.pl

clean:
	-rm $(ALL) *.dvi rulebook.aux rulebook.latex *.log *~ src/*~
	-rm -r book
