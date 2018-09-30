"""
Provide single test run functionality.
"""
from multiprocessing import Pool, Process

from utils.cycle_list import CycleList
from runner.data.accounts import AccountsData
from runner.gess import GessNodes
from settings.transaction import (
    TRANSACTION_GAS,
    TRANSACTION_GAS_PRICE,
    TRANSACTION_VALUE,
)
from utils.log import log


TRANSACTIONS_MAP = CycleList()
TRANSACTIONS_MAP.extend([
    (0, 2),
    (1, 0),
    (2, 0),
    (3, 0),
    (4, 0),
])

from settings.nodes import NODES_HOSTS
test_accounts = {}

for node_host in NODES_HOSTS:
    accs_list = CycleList()
    accs_list.extend([str(i + 1) + ':' + node_host for i in range(5)])
    test_accounts.update({
        node_host: accs_list
    })


class SingleNodeRun:
    """
    Single node run implementation.
    """

    def __init__(self, node_index, load_factor):
        self.node_index = node_index
        self.load_factor = load_factor
        self.gess_nodes = GessNodes()
        self.overall_transactions_count = self.load_factor * len(self.gess_nodes)
        self.transactions_performed = 0

        self.accounts = AccountsData().accounts
        self.accounts = test_accounts

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

        target_node_index = self.node_index + TRANSACTIONS_MAP[self.transactions_performed][0]
        target_address_index = TRANSACTIONS_MAP[self.transactions_performed][1]

        target_node = self.gess_nodes[target_node_index]
        target_address = self.accounts.get(target_node.host)[target_address_index].get('address')

        # print(f'N{self.node_index + 1}#{self.transactions_performed + 1} Performing operation for: {source_address} -> {target_address}')

        from services.wallet.transaction import FailedToCreateTransaction
        try:
            source_node.wallet_transaction.create(
                from_=source_address,
                to=target_address,
                gas=TRANSACTION_GAS,
                gas_price=TRANSACTION_GAS_PRICE,
                value=TRANSACTION_VALUE,
            )

        except FailedToCreateTransaction:
            log('Tried to send transaction:\nfrom', source_address,
                '\nto', target_address,
                '\ngas', TRANSACTION_GAS,
                '\ngas_price', TRANSACTION_GAS_PRICE,
                '\nvalue', TRANSACTION_VALUE
                )

        self.transactions_performed += 1

    def run(self):
        """
        Perform all transactions.
        """
        while self.transactions_performed < self.overall_transactions_count:
            self._single_run()
        # pool = Pool()
        # for i in range(self.overall_transactions_count):
        # pool.map(self._single_run, [i for i in range(self.overall_transactions_count)])
