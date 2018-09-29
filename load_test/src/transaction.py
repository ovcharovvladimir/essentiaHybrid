"""
Provide transaction functionality.
"""
import json
from multiprocessing import Pool
from time import time
import requests

import rpyc

DEFAULT_TRANSACTIONS_COUNT = 1000

HOST = '18.224.0.169'
PORT = 8545
FROM_ADDRESS = '0x99e8b93282e722070b3f1865207adc6aff497f9c'
TO_ADDRESS = ''


class TransactionPool:
    """
    Transaction pool implementation.
    """

    def __init__(self, transactions_count=1000):
        self._pool = Pool(transactions_count)
        self._start_time = 0
        self._end_time = 0
        self.is_running = False
        self.transactions_count = transactions_count
        self.url = f'http://{HOST}:{PORT}'

    def start(self):
        """
        Start pool of transactions.
        """
        self._start_time = time()
        self._end_time = 0
        self.is_running = True

        self._pool.map(self._send_transaction, [i for i in range(self.transactions_count)])

        self.is_running = False

    def _send_transaction(self):
        """
        Perform single transaction send.
        """
        pass

    def get_elapsed_time(self):
        """
        Return time it took to run transactions.
        """
        if self.is_running:
            return -1

        return self._start_time - self._end_time

    def test_transaction(self):
        """
        Test transaction run via RPC.
        """
        connection = rpyc.classic.connect(host='localhost', port=51903)

        connection.eth_sendTransacstion(json.dumps({
            'from': FROM_ADDRESS,
            'to': TO_ADDRESS,
            'value': 1,
            'gasLimit': 1,
            'gasPrice': 1,
        }))

        # connection.eth_getBalance(json.dumps({
        #     'from': FROM_ADDRESS,
        #     'to': TO_ADDRESS,
        #     'value': 1,
        #     'gasLimit': 1,
        #     'gasPrice': 1,
        # }))
