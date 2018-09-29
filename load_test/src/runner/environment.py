"""
Provide runner environment functionality.
"""
from time import sleep

from runner.data.accounts import AccountsData
from services.node import Node
from settings.accounts import ACCOUNTS_PER_NODE, BANK_ACCOUNT
from settings.nodes import Nodes, get_node_url
from settings.transaction import (
    TRANSACTION_GAS,
    TRANSACTION_GAS_PRICE,
    TRANSACTION_VALUE,
)
from utils.log import log, log_in_line
from utils.values import hex_to_int

DEFAULT_LOAD_FACTOR=1000

SUCCESS_MESSAGE = 'SUCCESS'
FAILED_MESSAGE = 'FAILED'


class RunnerEnvironment:
    """
    Runner environment actions implementation.
    """

    def __init__(self, nodes_count=1, load_factor=DEFAULT_LOAD_FACTOR):
        self.nodes_count = nodes_count
        self.load_factor = load_factor
        self.accounts_data = AccountsData()
        self.bank_node = Node(host=get_node_url(BANK_ACCOUNT.get('host')))

    def _count_single_account_needed_funds(self, transaction_price):
        """
        Count funds needed for a single account to perform test run.
        """
        return transaction_price * self.load_factor

    def _create_accounts(self, count, node):
        """
        Create accounts on a given node and store them.
        """
        for i in range(count):
            account_address = node.account.create()

            self.accounts_data.add_account(node_host=node.host, address=account_address)

            node.account.unlock(address=account_address)

    def _top_up_account(self, address, value):
        """
        Top up account with funds from bank account.
        """
        self.bank_node.wallet_transaction.create(
            from_=BANK_ACCOUNT.get('address'),
            to=address,
            gas=TRANSACTION_GAS,
            gas_price=TRANSACTION_GAS_PRICE,
            value=value,
        )

    def _wait_for_funds_to_appear(self, single_node_funds):
        """
        Wait for funds to appear on all test nodes.
        """
        addresses_with_funds = 0

        while addresses_with_funds < self.nodes_count:
            for node in Nodes():
                import pdb; pdb.set_trace()
                target_address = self.accounts_data.accounts.get(node.host)[0]

                wallet_balance = hex_to_int(node.wallet_balance.get(address=target_address))

                if wallet_balance >= single_node_funds:
                    addresses_with_funds += 1

                sleep(1)

    def _bank_has_enough_funds(self):
        funds_to_run = self.load_factor * self.nodes_count * \
                       (TRANSACTION_GAS * TRANSACTION_GAS_PRICE + TRANSACTION_VALUE)

        bank_balance = hex_to_int(self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address')))

        return bank_balance >= funds_to_run

    def setup_accounts(self):
        """
        1. Create N accounts for every node.
        3. Unlock EVERY account.

        4. Unlock bank account.
        5. Count the amount of funds needed for SINGLE account.
        6. Send funds to EVERY account.

        7. Wait for funds to be received.

        Return bool as status of success.
        """
        log('Setup.')
        log('Check if bank acount has enough funds ', )
        if not self._bank_has_enough_funds():
            log(FAILED_MESSAGE)
            return False
        log(SUCCESS_MESSAGE)

        log('Unlock bank account ')
        if self.bank_node.account.unlock(address=BANK_ACCOUNT.get('address'), password=BANK_ACCOUNT.get('password')):
            log(SUCCESS_MESSAGE)
        else:
            log(FAILED_MESSAGE)
            return False

        single_node_funds = self._count_single_account_needed_funds(
            transaction_price=TRANSACTION_VALUE
        )

        # for node in Nodes():
        for i in range(self.nodes_count):
            log(f'Top up account of node #{i + 1}.')
            node = Nodes()[i]

            self._create_accounts(count=ACCOUNTS_PER_NODE, node=node)

            target_address = self.accounts_data.accounts.get(node.host)[0]

            self._top_up_account(address=target_address, value=single_node_funds)

        log('Wait for funds to appear on the topped up accounts...')
        self._wait_for_funds_to_appear(single_node_funds=single_node_funds)

        return True

    def cleanup(self):
        """
        1. Go through every address on every node and send all funds on it to the bank account.
        """
        log('Cleanup.')

        value_sum = 0

        log('Send all available funds back to the bank account...')
        # for node in Nodes():
        for i in range(self.nodes_count):
            node = Nodes()[i]

            for account_address in self.accounts_data.accounts.get(node.host):
                value = node.wallet_balance.get(address=account_address)

                node.wallet_transaction.create(
                    from_=account_address,
                    to=BANK_ACCOUNT.get('address'),
                    gas=TRANSACTION_GAS,
                    gas_price=TRANSACTION_GAS_PRICE,
                    value=value,
                )

                value_sum += value

        log('Wait while funds appear at bank account.')
        bank_account_funds = 0
        while bank_account_funds < value_sum:
            bank_account_funds = hex_to_int(self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address')))
            sleep(1)
