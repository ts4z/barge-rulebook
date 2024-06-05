#
# barge-rulebook-md Makefile
#

all: $(ALL)

book:
	mdbook build

clean:
	-rm $(ALL) $(PROBLEMATIC) *.dvi rulebook.latex *.log *~ src/*~
	-rm -r book
