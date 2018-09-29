"""
Provide node interface.
"""
from services.accounts import Account
from services.wallet.balance import WalletBalance
from services.wallet.transaction import WalletTransaction


class Node:
    """
    Node interface implementation.
    """

    def __init__(self, host):
        self.host = host
        self.account = Account(url=self.host)
        self.wallet_balance = WalletBalance(url=self.host)
        self.wallet_transaction = WalletTransaction(url=self.host)
