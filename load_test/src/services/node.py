"""
Provide node interface.
"""
from services.accounts import Account
from services.tx_pool import TransactionPool
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
        self.tx_pool = TransactionPool(url=self.host)

    def get_next_nonce_for(self, address):
        """
        Get next valid nonce for address.
        """
        transactions_count = self.wallet_transaction.get_count(address=address)
        queued_transactions = self.tx_pool.get_content().get('queued').get(address)

        try:
            ql = len(queued_transactions)
        except TypeError:
            ql = 0

        print(f'TC: {transactions_count}; QT: {ql}')

        if queued_transactions is None:
            return transactions_count

        return transactions_count + len(queued_transactions)
