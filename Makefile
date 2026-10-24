#
# barge-rulebook Makefile
#

.PHONY: all clean spotless-repo

ALL=rulebook.pdf
all: $(ALL)

#.......890123456789012345678901234567890123456789012345678901234567890123456789
help:
	-@test -x /usr/bin/banner && banner rulebook
	@echo "barge-rulebook targets:"
	@echo " rulebook.pdf ................................build PDF rulebook"
	@echo " update-rulebook-on-barge-site ... publish PDF to barge.org repo"
	@echo
	@echo "That is, you probably want:"
	@echo "  make update-rulebook-on-barge-site"
	@echo
	@echo "The primary purpose of this Makefile is to capture the dance of"
	@echo "running LaTeX thrice (!) in order to get the pagination right,"
	@echo "checking out the barge.org repo, updating it, and pushing our"
	@echo "outputs into it."
	@echo
	@echo "update-rulebook-on-barge-site is the key rule here.  HTML"
	@echo "builds are fully automated (see the GitHub workflows).  PDF"
	@echo "builds are *not* automatically triggered, and need to use this."
	@echo
	@echo "The PDF machinery requires a Unicode-capable LaTeX install,"
	@echo "a working recent Go compiler, and a left-handed flathead"
	@echo "screwdriver in order to operate.  This makes for an absurdly"
	@echo "large (and so unavailable) Docker image, at least on our free"
	@echo "GitHub hosting plan.  Sorry not sorry."
	@echo

book: src/*.md src/*.latex
	mdbook build

# Note: Three passes are needed:
# 1. longtable needs two passes to compute column widths
# 2. marginnote needs position data from the aux file to resolve
#    x-positions on even/odd pages
# 3. a third pass to with the data marginnote needs
%.pdf: %.latex
	-rm rulebook.toc rulebook.aux # throw out aux files
	echo First pagesetting collects TOC to aux file
	xelatex -interaction=batchmode $< || (rm $@; false)
	rm $@
	echo Second pagesetting gets marginnotes fixed
	xelatex -interaction=batchmode $< || (rm $@; false)
	rm $@
	echo Generate final PDF with recomputed things
	xelatex -interaction=batchmode $< || (rm $@; false)

latexify/latexify: latexify/*.go latexify/go.mod latexify/go.sum
	(cd latexify && go build .)

rulebook.latex: latexify/latexify src/*.md src/*.latex
	latexify/latexify --need-lines 8

clean:
	-rm $(ALL) *.dvi rulebook.aux rulebook.latex rulebook.toc \
		rulebook.out *.log *~ src/*~
	-rm -rf book tmp

checkout-gh-pages:
	if ! [ -d tmp/gh-pages ]; then \
		rm -rf tmp; \
		mkdir tmp; \
		(cd tmp; \
		 git clone --depth 1 git@github.com:ts4z/ts4z.github.io gh-pages);\
	fi

sync-gh-pages: checkout-gh-pages
	cd tmp/gh-pages && git pull --rebase

tmp/gh-pages/barge-rulebook/rulebook.pdf: sync-gh-pages rulebook.pdf
	cp rulebook.pdf tmp/gh-pages/barge-rulebook/rulebook.pdf

publish-pdf: tmp/gh-pages/barge-rulebook/rulebook.pdf
	cd tmp/gh-pages && git add barge-rulebook/rulebook.pdf
	cd tmp/gh-pages && git commit -m \
		"Update rulebook.pdf from barge-rulebook repo."
	cd tmp/gh-pages && git push

tmp/barge.org:
	if ! [ -d tmp/barge.org ]; then \
		rm -rf tmp; \
		mkdir tmp; \
		cd tmp && git clone --depth 1 git@github.com:ts4z/barge.org barge.org; \
	fi
	cd tmp/barge.org && git reset --hard origin/main
	touch tmp/barge.org

sync-barge-site: tmp/barge.org
	cd tmp/barge.org && git fetch && git reset --hard origin/main

spotless-repo:
	@[ -z "`git status --porcelain`" ] || ( \
		echo 1>&2 "ERROR: Dirty git repo detected." ; \
		echo 1>&2 "Commit your changes before publishing." ; \
	false)

tmp/barge.org/static/rulebook.pdf: rulebook.pdf sync-barge-site spotless-repo
	cp rulebook.pdf tmp/barge.org/static/rulebook.pdf
	cd tmp/barge.org && (\
		git commit -m \
		"Update rulebook.pdf as of source repo rev `cd ../.. && git rev-parse --short=7 HEAD`." \
		static/rulebook.pdf ) || \
		(echo "commit failed, clean up:"; git reset --hard origin/main && false)

push-barge-site: 
	cd tmp/barge.org && git push

update-rulebook-on-barge-site: sync-barge-site tmp/barge.org/static/rulebook.pdf push-barge-site

love:
	@echo "not war?"
