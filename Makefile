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

DHUSER=ts4z
REV=`git rev-parse --short=7 HEAD`
REPO=$(DHUSER)/debian-xetex-commonmark

build-docker-image:
	docker build -t "$(REPO):$(REV)" .
	docker tag "$(REPO):$(REV)" "$(REPO):latest"
	@echo "HEAD 7: $(REV)"

push-docker-image: build-docker-image
	docker push -a "$(REPO)"
	@echo "pushed: $(REPO):$(REV) and all tags"
