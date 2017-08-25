#!/usr/bin/env python
import sys
import glob
import json
import pandas
import unipath

PLAYBOOKS = unipath.Path('~/foldit/playbooks')
TOP_DATA = PLAYBOOKS + 'data/top'
SOLUTION_DATA = TOP_DATA + '/run-*/solution_data.json'


def extract_solution_filenames(solution_data):
    solution_filenames = []
    return solution_filenames


if __name__ == '__main__':
    assert len(sys.argv) == 2 and unipath.Path(sys.argv[1]).exists()
    available_solutions = pandas.read_table(sys.argv[1], names=['path']).path

    already_downloaded = []
    for solution_data in glob.glob(SOLUTION_DATA):
        solution_filenames = extract_solution_filenames(solution_data)
        already_downloaded.extend(solution_filename)

    solutions_not_downloaded = available_solutions.ix[
        ~available_solutions.isin(already_downloaded)
    ]
    solutions_not_downloaded.to_csv(sys.stdout, index=False)
