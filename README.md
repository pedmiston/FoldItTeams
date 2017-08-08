# The impact of inheritance on open-ended problem solving

This project investigates the impact of inheritance on problem solving ability in the scientific discovery game **FoldIt**. Problem solving in **FoldIt** involves folding amino acid structures into three-dimensional protein structures. Individuals may work alone or together by sharing partial solutions for others to inherit. This investigation measures the impact of inheritance on problem solving ability by comparing the effectiveness of teams of problem solvers to individuals working alone. The goal of this investigation is to uncover the tradeoffs involved in sharing partial solutions to problems among groups of problem solvers in open-ended problem solving contexts.

## Finding puzzles

Solution data is stored on the server at <analytics.fold.it>. Solution data
exists in two forms: the raw `ir_solution` data (binary), and the converted
`pdb` data (plaintext). The downloading of the `ir_solution` data as well as
it's conversion to `pdb` is handled automatically.

The "find_puzzles.yml" playbook is designed to help locate the puzzles
with available data so the solutions can be downloaded and processed
before they disappear from the server.

```bash
# Find all puzzles with "all/" directory containing all solutions.
# Returns a list of directories.
# Results are saved at "data/puzzles_with_all_data.txt"
ansible-playbook playbooks/find_puzzles.yml -e type=all

# Find all puzzles with "top/solution_bid*.pdb" files.
# Returns a list of files.
# Results are saved at "data/puzzles_with_top_bid_data.txt"
ansible-playbook playbooks/find_puzzles.yml -e "{type: 'top', rank: 'bid'}"
```

## Parsing data from file paths

There is a lot of data contained in the paths to `pdb` files. The script
"parse_solution_paths.py" extracts the data in a txt file of paths
and outputs the results to a csv.

```bash
python scripts/parse_solution_paths.py  # outputs data/top_bid_solutions.csv
```

# Data

## SketchbookPuzzles

```bash
# Create a txt file of all top sketchbook puzzles on the analytics server
ansible-playbook playbooks/find_puzzles.yml -e @playbooks/vars/sketchbook.yml"

# Download the pdb files from the analytics server
ansible-playbook playbooks/download_puzzles.yml -e @playbooks/vars/sketchbook.yml

# Extract the histories from the top sketchbook puzzles
python scripts/extract_histories.py playbooks/data/sketchbook/top_pdb_files.txt playbooks/data/top_solutions/puzzle_2003996
```

## BestScores



## Contributing

1. Install the python requirements in a virtualenv.

```bash
python3 -m venv ~/.venvs/foldit
source ~/.venvs/foldit/bin/activate
pip install -r requirements.txt
```

2. Configure `osf-cli`.

```bash
osf init  # project id is "72txj"
echo "export OSF_PASSWORD=my-osf-password" >> .env
```

3. Set up Ansible Vault for securing secrets.

```bash
echo "my-vault-pass" > .vault_pass.txt
echo "export ANSIBLE_VAULT_PASSWORD_FILE=~/path/to/project/.vault_pass.txt" >> .env
```
