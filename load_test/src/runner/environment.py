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
from runner.nodes import GessNodes
from settings.transaction import (
    TRANSACTION_GAS,
    TRANSACTION_GAS_PRICE,
    TRANSACTION_VALUE,
)
from runner.logger import log
from utils.wei import wei_to_gwei
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
        1. Get nodes accounts
        2. Get known accounts
        3. Actualize the list of known accounts
        (Removes absent accounts, adds missing accounts)
        """
        # TODO: check if saved account exists on the node, if not, then discard it from the list and create a new one
        known_accounts = self.accounts_data.accounts.get(node.host)
        node_accounts = node.account.get_all()
        existing_accounts = []

        for i in range(count):
            account_exists = False
            account_address = ''
            account_password = ''

            if known_accounts:
                for known_account in known_accounts:
                    if known_account.get('address') in node_accounts:
                        account_address = known_account.get('address')
                        account_password = known_account.get('password')
                        existing_accounts.append({
                            'address': account_address,
                            'password': account_password,
                        })

                        log.info(f'Account exists: {node.host}:{account_address}')
                        known_accounts.remove(known_account)
                        account_exists = True
                        break

            if not account_exists:
                account_address = node.account.create()
                account_password = DEFAULT_ACCOUNT_PASSWORD
                existing_accounts.append({
                    'address': account_address,
                    'password': account_password,
                })
                log.info(f'Created account: {node.host}:{account_address}')

            log.info(f'Unlock account: {node.host}:{account_address}')
            node.account.unlock(address=account_address, password=account_password)

        self.accounts_data.set_actual_accounts_for_node(node_host=node.host, accounts_list=existing_accounts)

    def _top_up_account(self, address, value):
        """
        Top up account with funds from bank account.
        """
        log.info(f'Top up account {address} with {value}.')
        self.bank_node.wallet_transaction.create_raw(
            from_=BANK_ACCOUNT.get('address'),
            private_key=BANK_ACCOUNT.get('pk'),
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

                for account in self.accounts_data.accounts.get(node.host):
                    target_address = account.get('address')

                    wallet_balance = node.wallet_balance.get(address=target_address)

                    print(f'Single node balance: {single_node_funds}; Wallet balance: {wallet_balance}')
                    if wallet_balance >= single_node_funds:
                        log.info(f'Account: {target_address} has enough funds.')
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
        log.info(f'-------------------------------------------------------------')
        log.info(f'--- New session started on {datetime.strftime(datetime.now(), "%d %b %y at %H:%M:%S")}')
        log.info(f'-------------------------------------------------------------')

        log.debug('Setup.')
        log.debug('Check if bank account has enough funds ', )
        if not self._bank_has_enough_funds():
            log.debug(FAILED_MESSAGE)
            return False
        log.debug(SUCCESS_MESSAGE)

        # log.debug('Unlock bank account ')
        # if self.bank_node.account.unlock(address=BANK_ACCOUNT.get('address'), password=BANK_ACCOUNT.get('password')):
        #     log.debug(SUCCESS_MESSAGE)
        # else:
        #     log(FAILED_MESSAGE)
        #     return False

        single_node_funds = self._count_single_account_needed_funds(
            transaction_price=TRANSACTION_VALUE + (TRANSACTION_GAS_PRICE * TRANSACTION_GAS)
        )

        # for node in GessNodes():
        for i in range(self.nodes_count):
            node = self.gess_nodes[i]

            log.info(f'Create accounts for node #{i + 1}...')
            self._create_accounts(count=ACCOUNTS_PER_NODE, node=node)

            log.debug(f'Top up account of node #{i + 1}:{node.host}.')
            for account in self.accounts_data.accounts.get(node.host):
                target_address = account.get('address')
                target_address_funds = node.wallet_balance.get(address=target_address)
                log.info(f'Account: {target_address}; Funds: {target_address_funds}.')

                if target_address_funds < single_node_funds:
                    self._top_up_account(address=target_address, value=single_node_funds)

        log.debug('Wait for funds to appear on the topped up accounts...')
        self._wait_for_funds_to_appear(single_node_funds=single_node_funds)

        return True

    def save_accounts(self):
        """
        Save current accounts data to a file for future use.
        """
        self.accounts_data.save()

    def cleanup(self):
        """
        1. Go through every address on every node and send all funds on it bacj to the bank account.
        """
        log.debug('Cleanup.')

        self.accounts_data.save()

        bank_funds_before_refund = self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address'))
        value_sum = 0

        log.debug('Send all available funds back to the bank account...')
        # for node in GessNodes():
        for i in range(self.nodes_count):
            node = self.gess_nodes[i]

            for account in self.accounts_data.accounts.get(node.host):
                account_address = account.get('address')

                balance_value = node.wallet_balance.get(address=account_address)
                value = balance_value - (TRANSACTION_GAS * TRANSACTION_GAS_PRICE)
                if value <= 0:
                    log.warn(
                        f'Cannot send funds back from address {node.host}::{account_address}. '
                        f'Balance: {wei_to_gwei(balance_value)} (Need {wei_to_gwei(abs(value))} more).')
                    continue

                node.wallet_transaction.create(
                    from_=account_address,
                    to=BANK_ACCOUNT.get('address'),
                    gas=TRANSACTION_GAS,
                    gas_price=TRANSACTION_GAS_PRICE,
                    value=value,
                )

                value_sum += value
        print(f'Value sum: {value_sum}')

        log.debug('Wait while funds appear at bank account.')
        bank_account_funds = 0
        while (bank_account_funds - bank_funds_before_refund) < value_sum:
            bank_account_funds = self.bank_node.wallet_balance.get(address=BANK_ACCOUNT.get('address'))
            sleep(1)
