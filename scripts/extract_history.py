import collections
import re

import unipath


fields = collections.OrderedDict([
('pdl', re.compile('^IRDATA PDL')),
('timestamp', re.compile('^IRDATA TIMESTAMP')),
('energy', re.compile('^IRDATA ENERGY')),
('history', re.compile('^IRDATA HISTORY'))
])

pdl_fields = ['username', 'groupname', 'uid', 'gid', 'buildid', 'current_score', 'score_valid', 'best_score', 'action_log']

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

                    if field in pdl.keys():
                        data[field] += [pdl]
                    else:
                        data[field] = [pdl]

                # if field == 'history': ...

                else:
                    data[field] = [line.split()[2:]]

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
    args = parser.parse_args()
    assert path.exists(args.solution_pdb_paths), "solution paths not found"

    local_data_dir = unipath.Path(local_data_dir)
    if not local_data_dir.isdir():
        local_data_dir.mkdir(True)

    with open(args.solution_pdb_paths, 'r') as pdb_paths:
        for pdb_path in pdb_paths.read():
            pdb_path = pdb_path.strip()
            assert path.exists(pdb_path), "solution pdb file not found"
            data = extract_data(pdb_path, local_data_dir)
