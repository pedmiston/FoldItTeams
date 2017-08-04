from parse_solution_paths import read_solution_paths

if __name__ == '__main__':
    solutions = read_solution_paths('data/puzzles_with_top_bid_solutions.txt')
    print('\n'.join(map(str, solutions.puzzle_id.unique())))
