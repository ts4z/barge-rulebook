#!/bin/sh

# Add a bunch of meta http-equiv refresh redirects for book pages
# that either used to exist, or look like they should exist, but
# don't.
#
# mdbook doesn't support this natively, but we can run this script
# after we have built the book.

# This is not an ideal solution but works fine in the CI environment
# and in a wide variety of deployment environments.

target_dir=book

log() {
    echo 2>&1 "$@"
}

set -eu
# set -o pipefail || log "can't set -o pipefail"

list_aliases() {
    cat <<_EOF
common-rules-and-variations:common-rules
drawmadugi:dramadugi
drawmaha-21:dramaha-21
drawmaha-49:dramaha-49
drawmaha-zero:dramaha-zero
drawmaha:dramaha
holdem:texas-holdem
lowball:california-lowball
mississippi-stud-and-variants:mississippi-stud
nit-cards:nit-levels
omaha-8:omaha-high-low-eight-or-better
omaha-high:omaha-high-only
omaha:omaha-high-only
stud:seven-card-stud
variants:poker-games
_EOF
}

alias() {
    cat <<_EOF
<!DOCTYPE html>
<html>
  <head>
    <title>document moved</title>
    <link rel="canonical" href="${1}">
    <meta name="robots" content="noindex">
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <meta http-equiv="refresh" content="0; url=${1}">
  </head>
</html>
_EOF
}

list_aliases | (cd "$target_dir" && while IFS= read -r line ; do
    src=$(echo "$line" | cut -d: -f1).html
    dst=$(echo "$line" | cut -d: -f2).html

    log "alias: file $src redirects to $dst"
    alias "$dst" > "$src"
done
)
