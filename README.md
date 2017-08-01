# The impact of inheritance on open-ended problem solving

This project investigates the impact of inheritance on problem solving ability in the scientific discovery game **FoldIt**. Problem solving in **FoldIt** involves folding amino acid structures into three-dimensional protein structures. Individuals may work alone or together by sharing partial solutions for others to inherit. This investigation measures the impact of inheritance on problem solving ability by comparing the effectiveness of teams of problem solvers to individuals working alone. The goal of this investigation is to uncover the tradeoffs involved in sharing partial solutions to problems among groups of problem solvers in open-ended problem solving contexts.

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
