import collections
import re

import unipath
import pandas

from parse_solution_paths import read_solution_paths


fields = collections.OrderedDict([
('pdl', re.compile('^IRDATA PDL')),
('timestamp', re.compile('^IRDATA TIMESTAMP')),
('energy', re.compile('^IRDATA ENERGY')),
('history', re.compile('^IRDATA HISTORY'))
])

pdl_fields = ['username', 'groupname', 'uid', 'gid', 'buildid', 'current_score', 'score_valid', 'best_score', 'action_log']

PROJ = unipath.Path(__file__).absolute().ancestor(2)
SCRIPTS = unipath.Path(PROJ, 'scripts')
LOCAL_PDB_DIR = unipath.Path(PROJ, 'playbooks/data/top_solutions/puzzle_2003996')


class Solution:
    local_data_dir = None

    def __init__(self, solution_pdb):
        solution_pdb = unipath.Path(solution_pdb)

        if solution_pdb.exists():
            local_solution_pdb = solution_pdb
        else:
            local_solution_pdb = unipath.Path(self.local_data_dir, solution_pdb.name)
            if not local_solution_pdb.exists():
                download_solution_pdb(solution_pdb, local_solution_pdb)

        self.data = extract_data(local_solution_pdb)
        self.data['path'] = solution_pdb

    def get_best_scores(self):
        return self.get_row('uid', 'gid', 'timestamp', 'energy', 'path')

    def get_total_actions(self):
        action_log = pandas.Series(self.data['pdl']['action_log'])
        actions = (action_log.rename({'|': '|UnknownAction'})
                             .rename(lambda x: x.strip('|')))
        actions.index.name = 'name'
        actions.name = 'count'
        actions = actions.reset_index()

        for id_var in ['uid', 'gid', 'path']:
            actions[id_var] = self.get(id_var)

        col_order = ['uid', 'gid', 'name', 'count', 'path']
        return actions[col_order]

    def get_solution_history(self):
        history = pandas.Series(self.get('history').split(','), name='solution_id')
        history.index.name = 'solution_ix'
        history = history.reset_index()

        for id_var in ['uid', 'gid', 'path']:
            history[id_var] = self.get(id_var)

        col_order = ['uid', 'gid', 'solution_ix', 'solution_id', 'path']
        return history[col_order]

    def get_row(self, *data_args):
        row_data = [self.get(arg) for arg in data_args]
        return pandas.Series(row_data, index=data_args)

    def get(self, arg):
        data = self.data.get(arg)
        if data is None:
            try:
                data = getattr(self, arg)
            except AttributeError:
                raise NotImplementedError("don't know how to extract arg '%s'" % arg)
        return data


    @property
    def uid(self):
        return self.data['pdl']['uid']

    @property
    def gid(self):
        return self.data['pdl']['gid']


def extract_data(solution_pdb):
    solution_pdb_handle = open(solution_pdb, 'r')

    data = {}
    last_uid = 0
    for line in solution_pdb_handle.readlines():
        for field in fields:
            if fields[field].match(line):
                if field == 'pdl':
                    pdl = {}
                    splt = line.split(',')
                    v1 = splt[0].split()
                    if v1[2] != '.': # this is merged data, ignore this for now
                        continue
                    pdl['username'] = " ".join(v1[3:])
                    pdl['groupname'] = splt[1]
                    pdl['uid'] = splt[2]
                    if last_uid != splt[2]:
                        last_uid = splt[2]
                    else:
                        continue #skip this if we are just seeing multiple entries for the same player in a row
                    pdl['gid'] = splt[3]
                    pdl['buildid'] = splt[4]
                    pdl['current_score'] = splt[5]
                    pdl['score_valid'] = splt[6]

                    v1 = splt[7].split()


                    pdl['best_score'] = v1[0]

                    try:
                        if len(v1) > 2:
                            pdl['action_log'] = {key: value for (key, value) in [(t.split('=')[0],t.split('=')[1]) for t in v1[3:]]}
                    except Exception as e:
                        continue

                    data[field] = pdl

                # if field == 'history': ...

                else:
                    data[field] = line.split()[2:][0]

    solution_pdb_handle.close()
    return data


def download_solution_pdb(src, dst):
    raise NotImplementedError('looking for: %s\nat: %s' % (src, dst))


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('solution_pdb_paths')
    parser.add_argument('dest')
    args = parser.parse_args()

    solution_pdb_paths = unipath.Path(args.solution_pdb_paths)
    pdb_paths = [unipath.Path(path.strip()) for path in
                 open(args.solution_pdb_paths).readlines()]

    if not LOCAL_PDB_DIR.isdir():
        LOCAL_PDB_DIR.mkdir(True)

    Solution.local_data_dir = LOCAL_PDB_DIR

    # Concurrency here??

    best_scores = []
    total_actions = []
    solution_history = []

    for pdb_path in pdb_paths:
        pdb_path = pdb_path.strip()
        solution = Solution(pdb_path)
        best_scores.append(solution.get_best_scores())
        total_actions.append(solution.get_total_actions())
        solution_history.append(solution.get_solution_history())

    best_scores_data = pandas.DataFrame.from_records(best_scores)
    total_actions_data = pandas.concat(total_actions)
    solution_histories = pandas.concat(solution_history)

    path_data = read_solution_paths(args.solution_pdb_paths)

    def merge_path_data(frame):
        len_before_merge = len(frame)
        frame = frame.merge(path_data)
        assert len(frame) == len_before_merge
        del frame['path']
        return frame

    best_scores_data = merge_path_data(best_scores_data)
    total_actions_data = merge_path_data(total_actions_data)
    solution_histories = merge_path_data(total_actions_data)

    dest = unipath.Path(args.dest)
    if not dest.isdir():
        dest.mkdir()

    best_scores_data.to_csv(unipath.Path(dest, 'best_scores.csv'), index=False)
    total_actions_data.to_csv(unipath.Path(dest, 'total_actions.csv'), index=False)
    solution_histories.to_csv(unipath.Path(dest, 'solution_histories.csv'), index=False)
