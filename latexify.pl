#!/usr/bin/perl -w

# This script idiosyncratically produces a merged LaTeX document from the
# mdbook SUMMARY.md input by reaching out to all of the links.
#
# This is not a quality product, and it has many, many bugs.
#
# You may think to yourself, "Why not use one of the fine, high-quality
# mdbook-to-latex backends?"  Well, I looked for them, and I tried a couple.
# And I'm pretty disappointed with what you see here.  And yet, here it is.
#
# One problem with this is that it doesn't know how to render tables
# (CommonMark doens't understand tables, and we use it for rendering LaTeX).
#
# There are two scripts in here, struggling to get out.  One is a parser for
# SUMMARY.md that extracts just the files in order that have to be catted
# together to produce the "big" Markdown file.  The other is the md->latex
# converter.  Both use CommonMark.  The first one works well enough; the second
# one does not, particularly for tables.

use strict;
use Carp;
use CommonMark;
use Data::Dumper;
use File::Copy;
use utf8;

binmode(STDOUT, ":utf8");

my $rel = "src";

open(my $summary, '<:utf8', "$rel/SUMMARY.md")
     or die "can't open summary file";

sub walk_tree_collecting {
    my $node = shift;
    my $depth = shift;
    my $collector = shift;
    my $indent = "." x $depth;

    do {
        my $type = $node->get_type;
        if ($type == CommonMark::NODE_THEMATIC_BREAK) {
            &$collector({});
        } elsif ($type == CommonMark::NODE_LINK) {
            print "${indent}LINK: ", $node->get_url(), "\n";
            &$collector({filename=>$node->get_url()})
        } else {
            # print "${indent}looking at type $type\n";
            my $subtree = $node->first_child;
            if ($subtree) {
                # print "${indent}recurse into: ", $subtree->render_commonmark, "\n";
                walk_tree_collecting($subtree, 1+$depth, $collector)
            }
        }
    } while ($node = $node->next);
}

# Walk the tree of CM nodes and collect all of the links.  (We ignore the
# structure of SUMMARY.md.)
sub walk_tree {
    my $root = shift;
    my @nodes = ();
    my $collector = sub {
        my $file = shift;
        push @nodes, $file;
    };
    walk_tree_collecting($root, 0, $collector);
    return \@nodes;
}

my $doc = CommonMark->parse(file => $summary)
     or die "can't parse summary";
my $files = walk_tree($doc);

sub render_cm_file_as_latex {
    my $OUT = shift;
    my $info = shift;
    my $fn = $info->{filename};

    if (!defined($fn)) {
        return;
    }

    my $real_fn = "$rel/$fn";
    open(my $fh, '<:utf8', "$real_fn") or die "can't open summary file";
    my $doc = CommonMark->parse(file => $fh);
    close ($fh);
    my $latex = $doc->render_latex;

    # CommonMark outputs things starting at the section label, but our document
    # class wants chapters.  We want section numbering to be correct, so here
    # we are.
    $latex =~ s/\\section/\\chapter/;
    $latex =~ s/\\subsection/\\section/;

    print $OUT "%%%% $real_fn\n";
    print $OUT $latex;
}

my $out_filename = "rulebook.latex";
open(my $out, '>:utf8', "$out_filename")
     or croak "can't open $out_filename for write";

# Write LaTeX boilerplate.  This is idiosyncratic.
#
# We fix up a few symbols that XeLaTeX seems to have trouble with.
#
# Note this is filled with double escapes due to Perl "helping".
print $out <<EOF
\\documentclass[letterpaper,12pt,hidelinks]{book}
\\usepackage{fontspec}
\\usepackage{newcomputermodern}
\\usepackage{graphicx}
\\usepackage{hyperref}
\\usepackage{newunicodechar}
\\usepackage{bookmark}
\\usepackage[dvipsnames]{xcolor}

% Substitute Unicode characters to add color.  I've messed with this a few
% times, they seem to be better behaved in math mode.

\\newunicodechar{♥}{{\\color{red}\\ensuremath{♥}}}
\\newunicodechar{♠}{\\ensuremath{♠}}
\\newunicodechar{♣}{{\\color{ForestGreen}\\ensuremath{♣}}}
\\newunicodechar{♦}{{\\color{RoyalBlue}\\ensuremath{♦}}}

% Times character gets complainy if not in math mode.

\\newunicodechar{✕}{\\ensuremath{\\times}}

\\author{Christopher J. Mecklin\\\\
---
\\and
  Tim Showalter\\\\
\\texttt{tjs\@psaux.com}}
\\begin{document}
\\frontmatter
\\begin{titlepage}
\\centering
\\includegraphics[width=8cm]{src/barge-logo.png}
\\vfill
{\\LARGE BARGE Rulebook}
\\vfill
Christopher J. Mecklin\\\\
Tim Showalter
\\vfill
\\today
\\vfill
\\end{titlepage}
\\tableofcontents
EOF
     ;

render_cm_file_as_latex($out, {filename=>'./preface.md'});

     print $out "\n\n\\mainmatter\n\n";

my %skip = (
            'colophon' => 1,
            'killer-cards' => 1,
            'lowball-scales' => 1,
            'preface' => 1,
            'sevens-rule' => 1,
            'title' => 1,
            'what-beats-what' => 1,
           );
sub skippable {
    my $fn = shift;
    if (!defined($fn->{filename})) {
        return undef;
    }
    $fn->{filename} =~ m:([^/]+)\.md$: or return undef;
    my $bare = $1;
    return defined($skip{$bare});
}
my @without_frontmatter = grep { !skippable($_) } (@$files);

map { render_cm_file_as_latex($out, $_) } (@without_frontmatter);

print $out "\n\n\\backmatter\n\n";

# Bug -- we should deduce this from the splits in SUMMARY.md.
for my $md (# 'lowball-scales.md',
            'sevens-rule.md',
            'what-beats-what.md',
            # 'killer-cards.md',
            'colophon.md') {
    render_cm_file_as_latex($out, {filename=>$md});
}

print $out <<EOF
\\end{document}
EOF
     ;

close $out;
