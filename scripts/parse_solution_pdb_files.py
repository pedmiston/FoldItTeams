from parse_solution_paths import read_solution_paths


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('puzzle_id')
    args = parser.parse_args()


    solutions = pandas.read_csv('data/top_bid_solutions.csv')
