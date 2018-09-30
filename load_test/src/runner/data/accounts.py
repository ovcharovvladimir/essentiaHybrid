"""
Provide accounts storage.
"""
import json

from settings.accounts import DEFAULT_ACCOUNT_PASSWORD

ACCOUNTS_FILE_NAME = 'accounts.json'


class AccountsData:
    """
    Accounts storage implementation.
    """

    __instance = None

    def __new__(cls):
        if not AccountsData.__instance:
            AccountsData.__instance = object.__new__(cls)

        return AccountsData.__instance

    def __init__(self):
        self.accounts = {}
        self.load()

    def add_account(self, node_host, address, password=DEFAULT_ACCOUNT_PASSWORD):
        self.accounts.setdefault(node_host, []).append({'address': address, 'password': password})

    def save(self):
        """
        Save current accounts to file.
        """
        accounts_json = json.dumps(self.accounts)

        with open(ACCOUNTS_FILE_NAME, 'w') as accounts_file:
            accounts_file.write(accounts_json)

    def load(self):
        """
        Load account from file.
        """
        accounts_json = '{}'

        try:
            with open(ACCOUNTS_FILE_NAME, 'r') as accounts_file:
                accounts_json = accounts_file.read()

        except FileNotFoundError:
            pass

        self.accounts = json.loads(accounts_json)
