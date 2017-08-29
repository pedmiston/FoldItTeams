import glob
import json
import sqlite3
import invoke
import unipath
import pandas
import boto3

from tasks import s3
from tasks import r
from tasks import paths


@invoke.task
def csv(ctx):
    """Create structured csvs from unstructured json."""
    filename = paths.R_PKG_DATA_RAW + '/top.json'
    assert filename.exists(), 'run "inv s3.pull" to download the data'

    top_scores_dst = paths.R_PKG_DATA_RAW + '/top-scores.csv'
    top_actions_dst = paths.R_PKG_DATA_RAW + '/top-actions.csv'

    first_write = True
    for solution_json in open(filename, 'r'):
        solution_data = json.loads(solution_json)
        solution = TopSolution(solution_data)
        solution.get_score().to_csv(top_scores_dst, header=first_write,
                                    mode='w' if first_write else 'a')
        solution.get_actions().to_csv(top_actions_dst, header=first_write,
                                      mode='w' if first_write else 'a')



class TopSolution:
    TOP_SOLUTION_FIELDS = 'PuzzleID UserID GroupID Score RankType Rank LenHistory'.split()
    def __init__(self, data):
        self._data = data
        try:
            self._data['LenHistory'] = len(self._data['History'])
        except TypeError:
            self._data['LenHistory'] = 0

    def get_score(self):
        return pandas.DataFrame(
            {field: self._data[field] for field in self.TOP_SOLUTION_FIELDS},
            index=[0],
        )[self.TOP_SOLUTION_FIELDS]

    def get_actions(self):
        actions = pandas.Series(self._data['Actions'])
        actions.name = 'Count'
        actions.index.name = 'Action'
        actions = actions.reset_index()

        for i, field in enumerate(self.TOP_SOLUTION_FIELDS):
            actions.insert(i, field, self._data[field])

        return actions


namespace = invoke.Collection(s3, r, csv)
