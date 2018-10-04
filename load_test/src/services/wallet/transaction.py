"""
Provide wallet transaction functionality.
"""
from ethereum import transactions
import codecs

import rlp
from request.wrapper import RequestWrapper
from request.templates import get_request_json

from settings.transaction import MAX_SEND_RAW_RETRIES_COUNT
from utils.values import (
    add_hex_prefix,
    hex_to_int,
    strip_hex_prefix,
)

TRANSACTION_ERROR_CODE = -32000


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


class FailedToGetTransactionInfo(Exception):
    """
    Failed to get transaction info.
    """

    pass


class WalletTransaction(RequestWrapper):
    """
    Wallet transaction implementation.
    """

    def __init__(self, *args, **kwargs):
        super(WalletTransaction, self).__init__(*args, **kwargs)

        self.internal_counter = 0

    def get(self, hash_):
        """
        Get transaction info by hash.
        """
        transaction_info, error = self.send(json=get_request_json('eth_getTransactionByHash', hash_))

        if error:
            raise FailedToGetTransactionInfo(error)

        if transaction_info is None:
            raise FailedToGetTransactionInfo

        return transaction_info

    def is_mined(self, node_number, tx_number, hash_):
        """
        Check if transaction block is mined.
        """
        transaction_block_number = self.get(hash_=hash_).get('blockNumber')

        print(f'\n\n>>>>>>>> N{node_number}#{tx_number} TX HASH: {hash_}; TX BLOCKN: {transaction_block_number}\n\n')

        if transaction_block_number is None:
            return False

        transaction_block_number = hex_to_int(transaction_block_number)

        return transaction_block_number > 0

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
        transaction_hash = None

        if nonce is None:
            nonce = self.get_count(address=from_)

        print(f'<<<<<<<<<<<<<<<< NONCE IS: {nonce}')

        for _ in range(MAX_SEND_RAW_RETRIES_COUNT + 1):
            transaction_data = self._get_raw_data(
                to=to, gas=gas, gas_price=gas_price, value=value, private_key=private_key, nonce=nonce, data=data,
            )

            transaction_hash, error = self.send(json=get_request_json('eth_sendRawTransaction', transaction_data))

            if transaction_hash:
                return transaction_hash

            if error:
                # if error.get('code') == TRANSACTION_ERROR_CODE:
                #     nonce += self.get_count(address=from_)
                #
                #     continue

                raise FailedToCreateTransaction(error)

        raise FailedToCreateTransaction(f'Failed to create transaction after {MAX_SEND_RAW_RETRIES_COUNT} tries!')

    @staticmethod
    def _get_raw_data(to, gas, gas_price, value, private_key, nonce=0, data=None):
        """
        Get a RAW Transaction data.
        """

        # print(f'<<<<<<<<<<<<<<<< NONCE IS: {nonce}')

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

        return transactions_count

        # print(f'Tx COUNT | RCVD: {transactions_count}; INTR: {self.internal_counter}')
        #
        # if self.internal_counter < transactions_count:
        #     self.internal_counter = transactions_count
        #
        # self.internal_counter += 1
        #
        # return self.internal_counter
