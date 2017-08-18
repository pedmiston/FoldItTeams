#!/usr/bin/env python
import pandas
import unipath

PLAYBOOKS = unipath.Path(__file__).absolute().ancestor(2)
TOP_DATA = unipath.Path(PLAYBOOKS, 'data/top')

available_solutions = pandas.read_table(
    unipath.Path(TOP_DATA, 'available_solutions.txt'),
    names=['path']
)

downloaded_solutions_csv = unipath.Path(TOP_DATA, 'solution_data.csv')

if downloaded_solutions_csv.exists():
    downloaded_solutions = pandas.read_csv(downloaded_solutions_csv)
    downloaded_solutions = downloaded_solutions[['path']]
    solutions_not_downloaded = available_solutions.path[
        available_solutions.path.notin(downloaded_solutions.path)
    ]
else:
    solutions_not_downloaded = available_solutions

solutions_not_downloaded.path.to_csv(
    unipath.Path(TOP_DATA, 'solutions_not_downloaded.txt'),
    index=False
)
