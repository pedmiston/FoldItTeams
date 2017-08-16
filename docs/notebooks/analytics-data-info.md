From an email from Jeff Flatten:

Once logged in, you will find a directory 'fetch' in your home directory. Inside this dir is the solution data from Foldit puzzles. Each recorded puzzle has a directory in the form of solution_XXXXXX/ where the XX is the puzzle ID in the URL of that puzzle on the website. 

For example, http://fold.it/portal/node/2003903 is the puzzle "1397: Revisiting Puzzle 124: PDZ Domain", and has an ID of 2003903. So you would expect to find data for this puzzle in the directory solution_2003903. 

Inside these directories, data is classified into several different types.

The all/ directory contains all solutions recorded for this puzzle. It is organized into smaller directories just to keep a single directory from having tens of thousands of files.

The shared/ directory contains solutions that were specifically shared by players when they don't think the score function is capturing something interesting.

The top/ directory contains solutions ranked by score in several different categories. These solutions take the form of solution_PREFIX_RANK_XXXX_XXXX_XXXX.pdb

The XXXX numbers can be ignored.

The PREFIX is one of the following, and describes the group in which this solution is of rank RANK.

bid = best solution per group (where non-grouped players count as their own group)
gid = best per group
uid = best per player (evolver or soloist)
sid = best soloist per player
eid = best evolver per player
