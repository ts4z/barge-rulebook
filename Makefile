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
	-rm $(ALL) *.dvi rulebook.aux rulebook.latex rulebook.toc \
		rulebook.out *.log *~ src/*~
	-rm -rf book tmp

REV=`git rev-parse --short=7 HEAD`

checkout-gh-pages:
	if ! [ -d tmp/gh-pages ]; then \
		rm -rf tmp; \
		mkdir tmp; \
		(cd tmp; \
		 git clone git@github.com:ts4z/ts4z.github.io gh-pages);\
	fi

sync-gh-pages: checkout-gh-pages
	cd tmp/gh-pages && git pull --rebase

tmp/gh-pages/barge-rulebook/rulebook.pdf: sync-gh-pages rulebook.pdf
	cp rulebook.pdf tmp/gh-pages/barge-rulebook/rulebook.pdf

publish-pdf: tmp/gh-pages/barge-rulebook/rulebook.pdf
	cd tmp/gh-pages && git add barge-rulebook/rulebook.pdf
	cd tmp/gh-pages && git commit -m \
		"Update rulebook.pdf as of source repo rev $(REV)."
	cd tmp/gh-pages && git push

.PHONY: tmp/barge.org

tmp/barge.org:
	if ! [ -d tmp/barge.org ]; then \
		rm -rf tmp; \
		mkdir tmp; \
		cd tmp && git clone git@github.com:ts4z/barge.org barge.org; \
	fi

sync-barge-site: tmp/barge.org
	cd tmp/barge.org && git pull --rebase

tmp/barge.org/static/rulebook.pdf: rulebook.pdf
	cp rulebook.pdf tmp/barge.org/static/rulebook.pdf
	( cd tmp/barge.org && git add static/rulebook.pdf && git commit -m \
		"Update rulebook.pdf as of source repo rev $(REV)." ) || \
	( cd tmp/barge.org && git checkout -- static/rulebook.pdf && false )

push-barge-site: 
	cd tmp/barge.org && git push

update-rulebook-on-barge-site: sync-barge-site tmp/barge.org/static/rulebook.pdf push-barge-site
