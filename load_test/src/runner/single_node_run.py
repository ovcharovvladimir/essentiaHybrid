"""
Provide single test run functionality.
"""
from time import sleep
import random
from multiprocessing import Pool, Process

from utils.cycle_list import CycleList
from runner.data.accounts import AccountsData
from runner.nodes import GessNodes
from services.wallet.transaction import FailedToCreateTransaction
from settings.transaction import (
    TRANSACTION_GAS,
    TRANSACTION_GAS_PRICE,
    TRANSACTION_VALUE,
)
from reporter.tool import (
    ReporterTool,
    TX_STATUS_CONFIRMED,
    TX_STATUS_FAILED,
)
from runner.logger import log
from settings.timeouts import TRANSACTION_BLOCK_CHECK_TIMEOUT_SECONDS


TRANSACTIONS_MAP = CycleList()
TRANSACTIONS_MAP.extend([
    (0, 2),
    (1, 0),
    (2, 0),
    (3, 0),
    (4, 0),
])

RUNS_BEFORE_TRANSACTION_QUEUE_CHECK = 64

# from settings.nodes import NODES_HOSTS
# test_accounts = {}
#
# for node_host in NODES_HOSTS:
#     accs_list = CycleList()
#     accs_list.extend([str(i + 1) + ':' + node_host for i in range(5)])
#     test_accounts.update({
#         node_host: accs_list
#     })


class SingleNodeRun:
    """
    Single node run implementation.
    """

    def __init__(self, node_index, gess_nodes, load_factor):
        self.node_index = node_index
        self.load_factor = load_factor
        self.gess_nodes = gess_nodes
        self.overall_transactions_count = self.load_factor
        # self.overall_transactions_count = self.load_factor * len(self.gess_nodes)
        self.transactions_performed = 0

        self.accounts = AccountsData().accounts
        self.reporter = ReporterTool(logger=log)
        # self.accounts = test_accounts

    def _single_run(self):
        """
        A. SINGLE RUN
          Send transaction from this_node.address1 to next_node.address2
          for i in range(next_node_inex + 1, next_node_inex + 4):
              Send transaction from this_node.address1 to next_node+1.address1
        """
        source_node = self.gess_nodes[self.node_index]

        source_address = self.accounts.get(source_node.host)[0].get('address')

        random_transaction_index = random.randint(0, self.transactions_performed)
        target_node_index = self.node_index + TRANSACTIONS_MAP[random_transaction_index][0]
        target_address_index = TRANSACTIONS_MAP[random_transaction_index][1]

        target_node = self.gess_nodes[target_node_index]
        target_address = self.accounts.get(target_node.host)[target_address_index].get('address')

        try:
            transaction_hash = source_node.wallet_transaction.create(
                from_=source_address,
                to=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
            )
            self.reporter.transaction(
                node_number=self.node_index + 1,
                number=self.transactions_performed + 1,
                host=source_node.host,
                hash_=transaction_hash,
                from_address=source_address,
                to_address=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
            )

            while not source_node.wallet_transaction.is_mined(node_number=self.node_index + 1, tx_number=self.transactions_performed, hash_=transaction_hash):
                sleep(TRANSACTION_BLOCK_CHECK_TIMEOUT_SECONDS)

            self.reporter.transaction(
                node_number=self.node_index + 1,
                number=self.transactions_performed + 1,
                host=source_node.host,
                hash_=transaction_hash,
                from_address=source_address,
                to_address=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
                status=TX_STATUS_CONFIRMED
            )

        except FailedToCreateTransaction as exception:
            self.reporter.transaction(
                node_number=self.node_index + 1,
                number=self.transactions_performed + 1,
                host=source_node.host,
                hash_='(no hash)',
                from_address=source_address,
                to_address=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
                status=TX_STATUS_FAILED,
            )
            self.reporter.error(text=f'N{self.node_index + 1}:#{self.transactions_performed + 1}{str(exception)}')

        self.transactions_performed += 1

    def run(self):
        """
        Perform all transactions.
        """
        try:
            log.info(
                f'STARTED TRANSACTION SWARMING ({self.overall_transactions_count}) '
                f'FOR {self.gess_nodes[self.node_index].host}'
            )

            for i in range(self.overall_transactions_count):
                self._single_run()

            self.reporter.run_ended(
                node_index=self.node_index,
                transactions_performed=self.transactions_performed,
                transactions_expected=self.overall_transactions_count,
            )

        except Exception as exception:
            self.reporter.run_failed(node_index=self.node_index, error_message=str(exception))
            # log.debug(f'Process #{self.node_index + 1} HAS FAILED WITH ERROR:\{str(exception)}n')
