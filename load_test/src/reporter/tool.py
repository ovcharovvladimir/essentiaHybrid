"""
Provide reporter tool implementation.
"""
from reporter.formatter import (
    ReportFormatter,
    TX_STATUS_CREATED,
    TX_STATUS_CONFIRMED,
    TX_STATUS_FAILED,
)


class ReporterTool:
    """
    Reporter tool implementation.
    """

    def __init__(self, logger):
        self.format = ReportFormatter()
        self.logger = logger

    def start(self):
        self.logger.info(self.format.get_start_header(text='ESS BLOCKCHAIN TX TEST REPORT'))

    def start_bank_account(self, bank_account_data):
        bank_account_info = self.format.get_bank_account_info(
            address=bank_account_data.get('address'), balance=bank_account_data.get('balance'), stage='BEGIN',
        )

        self.logger.info(bank_account_info)

    def start_accounts(self, accounts_data):
        self.logger.info(self.format.get_nodes_and_accounts_table(accounts_data=accounts_data))

    def start_transactions(self):
        self.logger.info(self.format.get_transactions_header())

    def transaction(self, node_number, number, host, hash_, from_address, to_address, gas, gas_price, value, status):
        transaction_line = self.format.get_transaction(
            node_number=node_number,
            number=number,
            host=host,
            hash_=hash_,
            from_address=from_address,
            to_address=to_address,
            gas=gas,
            gas_price=gas_price,
            value=value,
            status=status
        )

        self.logger.info(transaction_line)

    def run_ended(self, node_index, transactions_performed, transactions_expected):
        self.logger.info(
            self.format.get_run_end(
                node_index=node_index,
                transactions_performed=transactions_performed,
                transactions_expected=transactions_expected,
            ),
        )

    def run_failed(self, node_index, error_message):
        self.format.get_run_failed(node_index=node_index, error_message=error_message)

    def end(self, bank_account_data):
        self.logger.info(self.format.get_sub_footer())

        bank_account_info = self.format.get_bank_account_info(
            address=bank_account_data.get('address'), balance=bank_account_data.get('balance'), stage='END',
        )

        self.logger.info(bank_account_info)

        # TODO: report summary
        # self.logger.info(self.format.get_summary_text())

    def sub_header(self, text):
        self.logger.info(self.format.get_sub_header(text=text))

    def error(self, text):
        self.logger.error(self.format.get_error(text=text))
