"""
Provide accounts storage.
"""
from utils.cycle_list import CycleList


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

    def add_account(self, node_host, address):
        self.accounts.setdefault(node_host, []).append(address)
