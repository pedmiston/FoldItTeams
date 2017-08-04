"""Extract relevant data from paths to solutions."""
__author__ = "Pierce Edmiston <pierce.edmiston@gmail.com>"

import re
import pandas

re_solution_path = re.compile(r"""
    /home/pierce/fetch/solution_(?P<puzzle_id>\d+)/top/
    solution_(?P<prefix>[a-z]+)_(?P<rank>\d+)_\d+_\d+_\d+.ir_solution.pdb
    """, re.X)

if __name__ == '__main__':
    solutions = pandas.read_csv('playbooks/data/puzzles_with_top_bid_solutions.txt',
                                names=['path'])

    # Extract data from path to solution data
    solution_path_data = solutions.path.str.extract(re_solution_path, expand=True)
    solution_path_data['puzzle_id'] = pandas.to_numeric(solution_path_data.puzzle_id)
    solution_path_data['rank'] = pandas.to_numeric(solution_path_data['rank'])

    solutions = (solutions.join(solution_path_data)
                          .sort_values(['puzzle_id', 'rank'])
                          .reset_index(drop=True))

    output_cols = 'puzzle_id prefix rank path'.split()
    assert set(output_cols) == set(solutions.columns)
    solutions[output_cols].to_csv('data/top_bid_solutions.csv', index=False)
