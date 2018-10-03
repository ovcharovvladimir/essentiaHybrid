"""
Provide wallet transaction functionality.
"""
from ethereum import transactions
import codecs

import rlp
from request.wrapper import RequestWrapper
from request.templates import get_request_json

from utils.values import (
    add_hex_prefix,
    hex_to_int,
    strip_hex_prefix,
)


class FailedToCreateTransaction(Exception):
    """
    Failed to create transaction.
    """

    pass


class FailedToGetTransactionCount(Exception):
    """
    Failed to get transaction count.
    """

    pass


class WalletTransaction(RequestWrapper):
    """
    Wallet transaction implementation.
    """

    def __init__(self, *args, **kwargs):
        super(WalletTransaction, self).__init__(*args, **kwargs)

        self.internal_counter = 0

    def create(self, from_, to, gas, gas_price, value):
        """
        Create transaction
        """
        transaction_hash, error = self.send(
            json=get_request_json(
                'eth_sendTransaction',
                from_=from_,
                to=to,
                gas=hex(int(gas)),
                gasPrice=hex(int(gas_price)),
                value=hex(int(value)),
            )
        )

        if error:
            raise FailedToCreateTransaction(error)

        if transaction_hash is None:
            raise FailedToCreateTransaction

        return transaction_hash

    def create_raw(self, from_, to, gas, gas_price, value, private_key, nonce=None, data=None):
        """
        Create a RAW Transaction.
        """
        if nonce is None:
            nonce = self.get_count(address=from_)

        transaction_data = self._get_raw_data(
            to=to, gas=gas, gas_price=gas_price, value=value, private_key=private_key, nonce=nonce, data=data,
        )

        transaction_hash, error = self.send(json=get_request_json('eth_sendRawTransaction', transaction_data))

        if error:
            raise FailedToCreateTransaction(error)

        if transaction_hash is None:
            raise FailedToCreateTransaction

        return transaction_hash

    @staticmethod
    def _get_raw_data(to, gas, gas_price, value, private_key, nonce=0, data=None):
        """
        Get a RAW Transaction data.
        """
        if data is None:
            data = '0x'

        binary_data = codecs.decode(strip_hex_prefix(data), 'hex')

        transaction = transactions.Transaction(nonce, gas_price, gas, to, value, binary_data)

        signed_transaction = transaction.sign(private_key)
        raw_transaction = add_hex_prefix(rlp.encode(signed_transaction).hex())

        return raw_transaction

    def get_count(self, address, quantity='pending'):
        """
        Get transactions count for a given address.
        """
        transactions_count, error = self.send(json=get_request_json('eth_getTransactionCount', address, quantity))

        if error:
            raise FailedToGetTransactionCount(error)

        if transactions_count is None:
            raise FailedToGetTransactionCount

        transactions_count = hex_to_int(transactions_count)

        print(f'Tx COUNT | RCVD: {transactions_count}; INTR: {self.internal_counter}')

        if self.internal_counter < transactions_count:
            self.internal_counter = transactions_count

        elif self.internal_counter == transactions_count:
            self.internal_counter += 1

        return self.internal_counter
