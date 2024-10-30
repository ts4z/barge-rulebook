#
# barge-rulebook-md Makefile
#

ALL=rulebook.pdf
all: $(ALL)

book: src/*.md
	mdbook build

# Note: This generates the book twice because longtable seems to indicate
# that it can't do its job the first time.
%.pdf: %.latex
	-rm rulebook.toc rulebook.aux # throw out aux files
	xelatex $< || rm $@
	rm $@			# output isn't right, aux files are good
	xelatex $< || rm $@

rulebook.latex: latexify.pl src/*.md src/*.latex
	./latexify.pl

clean:
	-rm $(ALL) *.dvi rulebook.aux rulebook.latex rulebook.toc rulebook.out *.log *~ src/*~
	-rm -r book

build-docker-image:
	docker build -t debian-xetex-commonmark:`git rev-parse --short=7 HEAD` .

push-docker-image: build-docker-image
	docker push debian-xetex-commonmark
