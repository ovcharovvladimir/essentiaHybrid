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
        self._all_accounts = {}
        self.accounts = {}
        self.load()

    def _update_all_accounts(self):
        """
        Actualize all_accounts list with current accounts list.
        """
        for host in self.accounts:
            accounts_list = self.accounts[host]
            for account in accounts_list:
                if self._all_accounts.get(host):
                    if account in self._all_accounts.get(host):
                        continue

                self._all_accounts.setdefault(host, []).append(account)

    def add_account(self, node_host, address, password=DEFAULT_ACCOUNT_PASSWORD):
        """
        Add account data under specified node host.
        """
        self.accounts.setdefault(node_host, []).append({'address': address, 'password': password})

    def set_actual_accounts_for_node(self, node_host, accounts_list):
        """
        Replace account list for a given node host with a new one.
        """
        self.accounts[node_host] = accounts_list

    def save(self):
        """
        Save current accounts to file.
        """
        self._update_all_accounts()
        accounts_json = json.dumps(self._all_accounts, indent=4)

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
        self._all_accounts = json.loads(accounts_json)
