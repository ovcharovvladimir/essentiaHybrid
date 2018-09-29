"""
Provide wallet transaction functionality.
"""
from request.wrapper import RequestWrapper
from request.templates import get_request_json


class FailedToCreateTransaction(Exception):
    """
    Failed to create transaction.
    """

    pass


class WalletTransaction(RequestWrapper):
    """
    Wallet transaction implementation.
    """

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
