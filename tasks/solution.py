import pandas

class TopSolution:
    TOP_SOLUTION_FIELDS = 'PuzzleID UserID GroupID Score RankType Rank LenHistory'.split()
    def __init__(self, data):
        self._data = data
        try:
            self._data['LenHistory'] = len(self._data['History'])
        except TypeError:
            self._data['LenHistory'] = 0

    def to_record(self):
        return pandas.Series(
            [self._data[field] for field in self.TOP_SOLUTION_FIELDS],
            index=self.TOP_SOLUTION_FIELDS
        )

    def get_actions(self):
        actions = pandas.Series(self._data['Actions'])
        actions.name = 'Count'
        actions.index.name = 'Action'
        actions = actions.reset_index()

        for i, field in enumerate(self.TOP_SOLUTION_FIELDS):
            actions.insert(i, field, self._data[field])

        return actions
