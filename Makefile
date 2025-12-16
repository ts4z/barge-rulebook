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
	xelatex -interaction=batchmode $<
	-rm $@			# output isn't right, aux files are good
	xelatex -interaction=batchmode $<

latexify/latexify: latexify/*.go latexify/go.mod latexify/go.sum
	(cd latexify && go build .)

rulebook.latex: latexify/latexify src/*.md src/*.latex
	latexify/latexify

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

tmp/barge.org/static/rulebook.pdf: rulebook.pdf sync-barge-site
	[ -z "`git status --porcelain`" ]
	cp rulebook.pdf tmp/barge.org/static/rulebook.pdf
	cd tmp/barge.org && (\
		git commit -m \
		"Update rulebook.pdf as of source repo rev `cd ../.. && git rev-parse --short=7 HEAD`." \
		static/rulebook.pdf ) || \
		(echo "commit failed, clean up:"; git reset --hard origin/main && false)

push-barge-site: 
	cd tmp/barge.org && git push

update-rulebook-on-barge-site: sync-barge-site tmp/barge.org/static/rulebook.pdf push-barge-site
