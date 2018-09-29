"""
Provide balance functionality.
"""
import json
import requests

from request.templates import get_request_json
from request.wrapper import RequestWrapper


class FailedToGetBalance(Exception):
    """
    Failed to get balance error.
    """


class WalletBalance(RequestWrapper):
    """
    Balance implementation.
    """

    def get(self, address, period='latest'):
        """
        Get balance.
        """
        balance, error = self.send(json=get_request_json('eth_getBalance', address, period))

        if error:
            raise FailedToGetBalance(error)

        if balance is None:
            raise FailedToGetBalance

        return balance
