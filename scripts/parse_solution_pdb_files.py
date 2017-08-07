from parse_solution_paths import read_solution_paths


def parse_puzzle_solutions(puzzle_id):
    pdb_files = filter_solutions_by_puzzle_id(puzzle_id)
   


def filter_solutions_by_puzzle_id(puzzle_id):
    solutions = pandas.read_csv('data/top_bid_solutions.csv')
    return solutions.ix[solutions.puzzle_id == puzzle_id, ]


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('puzzle_id')
    args = parser.parse_args()
    parse_puzzle_solutions(args.puzzle_id)
