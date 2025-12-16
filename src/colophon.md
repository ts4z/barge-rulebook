# Colophon

Both web and print versions of this book are both produced from the same files[^1].

[^1]: [https://github.com/ts4z/barge-rulebook/](https://github.com/ts4z/barge-rulebook/)

Our web version of this book is produced with
[mdbook](https://rust-lang.github.io/mdBook/), which produces a group of nice,
easy-to-navigate web pages.  The text is written in an extended flavor of the
[Markdown](https://daringfireball.net/projects/markdown/) format, which
is easy to write and easy to read even if you don't understand the technical
bits.

A [GitHub](https://github.com/) workflow is used to build the rulebook into the
version hosted at [barge.org](https://www.barge.org/).  GitHub (and the
underlying `git` utility) provide our version control.

Our PDF version, intended for printing and reading at the poker table,
is produced by converting those files into
[LaTeX](https://www.latex-project.org/).  When we did the 2024 rulebook,
we discovered mdbook's PDF support just wasn't as good as our
previous [LaTeX](https://www.latex-project.org/) version.  (LaTeX presents
a high bar and is renowned for its output quality.)
So we converted the mdbook version to LaTeX, and now that is available
for those wanting a hard copy.

The external tools we use are open source software and freely available
All of the other parts are available in our GitHub repository with the
rest of the book.
