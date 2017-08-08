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

def process_solution_pdb(solution_pdb, solution_dir):
    solution = Solution(solution_pdb, solution_dir)
    return extract_best_scores(solution),

class Solution:
    local_data_dir = None

    def __init__(self, solution_pdb):
        self.data = extract_data(solution_pdb, self.local_data_dir)
        self.data['path'] = solution_pdb

    def get_best_scores(self):
        return self.get_row('uid', 'gid', 'timestamp', 'energy', 'path')

    def get_row(self, *data_args):
        row_data = []
        for arg in data_args:
            data = self.data.get(arg)
            if data is None:
                try:
                    data = getattr(self, arg)
                except AttributeError:
                    raise NotImplementedError("don't know how to extract arg '%s'" % arg)
            row_data.append(data)
        return pandas.Series(row_data, index=data_args)

    @property
    def uid(self):
        return self.data['pdl']['uid']

    @property
    def gid(self):
        return self.data['pdl']['gid']


def extract_data(solution_pdb, solution_dir):
    solution_pdb_handle = get_or_download(solution_pdb, solution_dir)
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


def get_or_download(solution_pdb, local_data_dir):
    remote_path = unipath.Path(solution_pdb)
    expected_file = remote_path.name
    local_path = unipath.Path(local_data_dir, remote_path.name)
    if not local_path.exists():
        raise NotImplementedError()

    return open(local_path, 'r')


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('solution_pdb_paths')
    parser.add_argument('local_data_dir')
    parser.add_argument('dest')
    args = parser.parse_args()
    assert unipath.Path(args.solution_pdb_paths).exists(), "solution paths not found"

    local_data_dir = unipath.Path(args.local_data_dir)
    if not local_data_dir.isdir():
        local_data_dir.mkdir(True)

    pdb_paths = [unipath.Path(path.strip()) for path in
                 open(args.solution_pdb_paths, 'r').readlines()]

    Solution.local_data_dir = args.local_data_dir

    # Concurrency here??

    best_scores = []

    for pdb_path in pdb_paths:
        pdb_path = pdb_path.strip()
        solution = Solution(pdb_path)
        best_scores.append(solution.get_best_scores())

    best_scores_data = pandas.DataFrame.from_records(best_scores)

    path_data = read_solution_paths(args.solution_pdb_paths)

    len_before_merge = len(best_scores_data)
    best_scores_data = best_scores_data.merge(path_data)
    assert len(best_scores_data) == len_before_merge
    del best_scores_data['path']

    dest = unipath.Path(args.dest)
    if not dest.parent.isdir():
        dest.parent.mkdir()

    best_scores_data.to_csv(dest, index=False)
