from os import environ

import sqlalchemy
import ansible_vault
from unipath import Path

import foldit


def connect_to_db():
    url = "mysql+pymysql://{user}:{password}@{host}:{port}/{dbname}".format(
        user='foldit',
        password=get_from_vault('foldit_password'),
        host='localhost',
        port='3306',
        dbname='Totems',
    )
    con = sqlalchemy.create_engine(url)
    return con


def get_from_vault(key=None, vault_file='playbooks/vars/secrets.yml'):
    try:
        ansible_vault_password_file = environ['ANSIBLE_VAULT_PASSWORD_FILE']
    except KeyError:
        raise AssertionError('Set the ANSIBLE_VAULT_PASSWORD_FILE environment variable')
    ansible_vault_password = open(ansible_vault_password_file).read().strip()
    vault = ansible_vault.Vault(ansible_vault_password)
    secrets_yaml = Path(foldit.PROJ_ROOT, vault_file)
    data = vault.load(open(secrets_yaml).read())
    if key is None:
        return data
    else:
        return data.get(key)