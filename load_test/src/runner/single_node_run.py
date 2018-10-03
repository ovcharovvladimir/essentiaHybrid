"""
Provide single test run functionality.
"""
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
from runner.logger import log


TRANSACTIONS_MAP = CycleList()
TRANSACTIONS_MAP.extend([
    (0, 2),
    (1, 0),
    (2, 0),
    (3, 0),
    (4, 0),
])

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
        self.overall_transactions_count = self.load_factor * len(self.gess_nodes)
        self.transactions_performed = 0

        self.accounts = AccountsData().accounts
        # self.accounts = test_accounts

    def _single_run(self):
        """
        A. SINGLE RUN
          Send transaction from this_node.address1 to next_node.address2
          for i in range(next_node_inex + 1, next_node_inex + 4):
              Send transaction from this_node.address1 to next_node+1.address1
        """
        source_node = self.gess_nodes[self.node_index]

        # import pdb; pdb.set_trace()
        source_address = self.accounts.get(source_node.host)[0].get('address')

        # random_transaction_index = self.transactions_performed
        # random_transaction_index = random.randint(0, len(TRANSACTIONS_MAP) - 1)
        random_transaction_index = random.randint(0, self.transactions_performed)
        target_node_index = self.node_index + TRANSACTIONS_MAP[random_transaction_index][0]
        target_address_index = TRANSACTIONS_MAP[random_transaction_index][1]

        target_node = self.gess_nodes[target_node_index]
        target_address = self.accounts.get(target_node.host)[target_address_index].get('address')

        log.info(
            f'N{self.node_index + 1}#{self.transactions_performed + 1} '
            f'Performing operation for: {source_address} -> {target_address}'
        )

        try:
            # log.info(
            #     f'Sending to send transaction:\nfrom {source_address}'
            #     f'\nto {target_address}'
            #     f'\ngas {TRANSACTION_GAS}'
            #     f'\ngas_price {TRANSACTION_GAS_PRICE}'
            #     f'\nvalue {TRANSACTION_VALUE}'
            # )
            source_node.wallet_transaction.create(
                from_=source_address,
                to=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
            )

        except FailedToCreateTransaction:
            # log.info(
            #     f'Failed to send transaction:\nfrom {source_address}'
            #     f'\nto {target_address}'
            #     f'\ngas {TRANSACTION_GAS}'
            #     f'\ngas_price {TRANSACTION_GAS_PRICE}'
            #     f'\nvalue {TRANSACTION_VALUE}'
            # )
            log.error(
                f'FAILED! N{self.node_index + 1}#{self.transactions_performed + 1} '
                f'operation for: {source_address} -> {target_address}'
            )

        self.transactions_performed += 1

    def run(self):
        """
        Perform all transactions.
        """
        log.info(
            f'STARTED TRANSACTION SWARMING ({self.overall_transactions_count}) '
            f'FOR {self.gess_nodes[self.node_index].host}'
        )

        for i in range(self.overall_transactions_count):
            print(f'<<<<<<<<<<< TX PERFORMED: {self.transactions_performed}')
            self._single_run()

        # while self.transactions_performed < self.overall_transactions_count:
        #     print(f'<<<<<<<<<<< TX PERFORMED: {self.transactions_performed}')
        #     self._single_run()

        # pool = Pool()
        # for i in range(self.overall_transactions_count):
        # pool.map(self._single_run, [i for i in range(self.overall_transactions_count)])
