"""Process all of the solutions to a puzzle."""
__author__ = "Pierce Edmiston <pierce.edmiston@gmail.com>"

import unipath

def find_puzzles_with_available_data():
    playbook_results = 'playbooks/puzzles_with_all_solutions.txt'
    assert unipath.Path(playbook_results).exists(), \
        "File '{}' not found. Run ansible playbook 'find_puzzles_with_available_data'"

    puzzle_ids = []
    for puzzle_dir_with_all in open(playbook_results):
        solution_dir = unipath.Path(puzzle_dir_with_all).parent.stem
        puzzle_id = int(solution_dir.split('_')[1])
        puzzle_ids.append(puzzle_id)
    return puzzle_ids

if __name__ == '__main__':
    all_puzzles = find_puzzles_with_available_data()
