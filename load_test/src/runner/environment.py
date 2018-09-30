"""
Provide runner environment functionality.
"""
from datetime import datetime
from time import sleep

from runner.data.accounts import AccountsData
from services.node import Node
from settings.accounts import (
    ACCOUNTS_PER_NODE,
    BANK_ACCOUNT,
    DEFAULT_ACCOUNT_PASSWORD,
)
from settings.nodes import get_node_url
from runner.gess import GessNodes
from settings.transaction import (
    TRANSACTION_GAS,
    TRANSACTION_GAS_PRICE,
    TRANSACTION_VALUE,
)
from runner.logger import log
# from utils.log import log, log_in_line
# from utils.values import hex_to_int

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
        self.gess_nodes = GessNodes()

    def _count_single_account_needed_funds(self, transaction_price):
        """
        Count funds needed for a single account to perform test run.
        """
        return transaction_price * self.load_factor

    def _create_accounts(self, count, node):
        """
        Create accounts on a given node and store them.
        """
        # TODO: check count of saved accounts, check if they exist, create necessary amount
        known_accounts = self.accounts_data.accounts.get(node.host)

        for i in range(count):
            account_exists = True
            if known_accounts:
                try:
                    account_address = known_accounts[i].get('address')
                    account_password = known_accounts[i].get('password')
                except IndexError:
                    account_exists = False
            else:
                account_exists = False

            if not account_exists:
                account_address = node.account.create()
                self.accounts_data.add_account(node_host=node.host, address=account_address)
                account_password = DEFAULT_ACCOUNT_PASSWORD

            node.account.unlock(address=account_address, password=account_password)

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
            # for node in self.gess_nodes:
            for i in range(self.nodes_count):
                node = self.gess_nodes[i]

                # import pdb; pdb.set_trace()
                target_address = self.accounts_data.accounts.get(node.host)[0].get('address')

                wallet_balance = node.wallet_balance.get(address=target_address)

                if wallet_balance >= single_node_funds:
                    addresses_with_funds += 1

            sleep(1)

    def _bank_has_enough_funds(self):
        funds_to_run = self.load_factor * self.nodes_count * \
                       (TRANSACTION_GAS * TRANSACTION_GAS_PRICE + TRANSACTION_VALUE)

        bank_balance = self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address'))

        return bank_balance >= funds_to_run

    def setup_accounts(self):
        """
        1. Check if bank account has enough funds
        2. Unlock bank account

        3. Create accounts on gess nodes
        4. Top up every first account on the node
        5. Wait until funds are received to the accounts

        Return bool as status of success.
        """
        log.info(f'--- New session started at {datetime.strftime(datetime.now(), "%d %b %y %H:%M:%S")}')
        log.debug('Setup.')
        log.debug('Check if bank account has enough funds ', )
        if not self._bank_has_enough_funds():
            log.debug(FAILED_MESSAGE)
            return False
            log.debug(SUCCESS_MESSAGE)

            log.debug('Unlock bank account ')
        if self.bank_node.account.unlock(address=BANK_ACCOUNT.get('address'), password=BANK_ACCOUNT.get('password')):
            log.debug(SUCCESS_MESSAGE)
        else:
            log(FAILED_MESSAGE)
            return False

        single_node_funds = self._count_single_account_needed_funds(
            transaction_price=TRANSACTION_VALUE
        )

        # for node in GessNodes():
        for i in range(self.nodes_count):
            log.debug(f'Top up account of node #{i + 1}.')
            node = self.gess_nodes[i]

            self._create_accounts(count=ACCOUNTS_PER_NODE, node=node)

            target_address = self.accounts_data.accounts.get(node.host)[0].get('address')
            target_address_funds = node.wallet_balance.get(address=target_address)

            if target_address_funds < single_node_funds:
                self._top_up_account(address=target_address, value=single_node_funds)

        log.debug('Wait for funds to appear on the topped up accounts...')
        self._wait_for_funds_to_appear(single_node_funds=single_node_funds)

        return True

    def cleanup(self):
        """
        1. Go through every address on every node and send all funds on it bacj to the bank account.
        """
        log.debug('Cleanup.')

        self.accounts_data.save()

        value_sum = 0

        log.log.debug('Send all available funds back to the bank account...')
        # for node in GessNodes():
        for i in range(self.nodes_count):
            node = self.gess_nodes[i]

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

        log.debug('Wait while funds appear at bank account.')
        bank_account_funds = 0
        while bank_account_funds < value_sum:
            bank_account_funds = self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address'))
            sleep(1)
