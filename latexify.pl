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
# However, you can override any md file by writing a file ending in latex, so
# just do the tables in latex.  We do something similar for the title page,
# since we customize that for both versions.
#
# This seems to work but you are allowed to clean it up.

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
            &$collector(undef);
        } elsif ($type == CommonMark::NODE_LINK) {
            print "${indent}LINK: ", $node->get_url(), "\n";
            &$collector($node->get_url());
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
        if ($file) {
            $file =~ m:([^/]+)\.md$: or croak "can't parse filename $file";
            push @nodes, $1;
        } else {
            push @nodes, undef;
        }
    };
    walk_tree_collecting($root, 0, $collector);
    return \@nodes;
}

my $doc = CommonMark->parse(file => $summary)
     or die "can't parse summary";
my $files = walk_tree($doc);

sub render_cm_file_as_latex {
    my $OUT = shift;
    my $fn = shift;

    my $real_fn = "$rel/$fn.md";
    open(my $fh, '<:utf8', "$real_fn")
         or die "can't open md fragment file $real_fn";
    my $doc = CommonMark->parse(file => $fh);
    close ($fh);
    my $latex = $doc->render_latex;

    # CommonMark outputs things starting at the section label, but our document
    # class wants chapters.  We want section numbering to be correct, so here
    # we are.
    $latex =~ s/\\section/\\chapter/g;
    $latex =~ s/\\subsection/\\section/g;

    print $OUT "%%%% $real_fn (rendered Markdown)\n";
    print $OUT $latex;
}

sub render_latex_file {
    my $OUT = shift;
    my $fn = shift;

    # Latex override exists.
    my $latex_fn = "$rel/$fn.latex";
    open(my $IN, '<:utf8', "$latex_fn") or return undef;
    print $OUT "%%%% $fn (raw LaTeX)\n";
    while(<$IN>) { print $OUT $_ }
    close($IN);
    return 1;
}

sub render_file {
    my $OUT = shift;
    my $fn = shift;

    if (!$fn) {
        render_special($OUT);
        return;
    }

    if (render_latex_file($OUT, $fn)) {
        print "rendered $fn as latex\n";
        return;
    }
    
    print "rendered $fn as md\n";
    render_cm_file_as_latex($OUT, $fn);
}

my @special_files = (
                     "setup",
                     "frontmatter",
                     "mainmatter",
                     "backmatter",
                     "endmatter",
                    );

my $next_special = 0;

sub render_special {
    my $OUT = shift;
    my $fn = $special_files[$next_special];
    print "render special: $fn\n";
    ++$next_special;
    render_latex_file($OUT, $fn)
}

my $out_filename = "rulebook.latex";
open(my $OUT, '>:utf8', "$out_filename")
     or croak "can't open $out_filename for write";

print $OUT <<EOF
%%
%% $out_filename (produced by $0)
%%
EOF
     ;
render_special($OUT);
for my $file (@$files) { render_file($OUT, $file); }
render_special($OUT);

if ($next_special != scalar(@special_files)) {
    croak "$next_special special files printed, expected "
         . scalar(@special_files);
}

close $OUT;
