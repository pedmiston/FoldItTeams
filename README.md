# The impact of inheritance on open-ended problem solving

## Contributing

1. Install the python requirements in a virtualenv.
2. Configure `osf-cli`.
3. Export `OSF_PASSWORD` and `ANSIBLE_VAULT_PASSWORD_FILE` environment variables.
4. Create `.vault_pass.txt`.

```bash
python3 -m venv ~/.venvs/foldit
source ~/.venvs/foldit/bin/activate
pip install -r requirements.txt
osf init
echo "export OSF_PASSWORD=my-osf-password" > .env
echo "export ANSIBLE_VAULT_PASSWORD_FILE=~/path/to/project/.vault_pass.txt" > .env
vim .vault_pass.txt
source .env
```
