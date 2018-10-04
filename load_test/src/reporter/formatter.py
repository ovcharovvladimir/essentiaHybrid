"""
Provide report formatter.
"""
from datetime import datetime

TX_STATUS_CREATED = '~'
TX_STATUS_CONFIRMED = '+'
TX_STATUS_FAILED = '-'

SCREEN_WIDTH = 80
REPORT_DATETIME_FORMAT = '%d %b %y at %H:%M:%S'
TRANSACTION_DATETIME_FORMAT = '%H:%M:%S.%f'
HEADER_SYMBOL = '='
SUB_HEADER_SYMBOL = '*'
ERROR_SYMBOL = '!'
LEFT_COLUMN_SIZE_FACTOR = 0.375
RIGHT_COLUMN_SIZE_FACTOR = 1 - LEFT_COLUMN_SIZE_FACTOR


class ReportFormatter:
    """
    Report formatter implementation.
    """

    @staticmethod
    def _get_screen_wide_symbol_string(symbol):
        """
        Fill screen with a given symbol and return as a string.
        """
        return str(symbol) * SCREEN_WIDTH

    @staticmethod
    def _get_centered(text, symbol='', width=SCREEN_WIDTH):
        return '{:{}^{}}'.format(text, symbol, width)

    def get_start_header(self, text, decor_symbol=HEADER_SYMBOL):
        """
        Create header.
        """
        header = self._get_screen_wide_symbol_string(symbol=decor_symbol)
        header += f'\n{self._get_centered(text)}'
        header += f'\n{self._get_centered(datetime.strftime(datetime.now(), REPORT_DATETIME_FORMAT))}\n'
        header += self._get_screen_wide_symbol_string(symbol=decor_symbol)

        return header

    def get_sub_header(self, text):

        return self._get_centered(text=text, symbol=SUB_HEADER_SYMBOL)

    def get_sub_footer(self):
        return self._get_screen_wide_symbol_string(symbol=SUB_HEADER_SYMBOL)

    def get_error(self, text):
        error_text = self._get_screen_wide_symbol_string(symbol=ERROR_SYMBOL)
        error_text += f'\n{self._get_centered(text=text,symbol=ERROR_SYMBOL)}'
        error_text += f'\n{self._get_screen_wide_symbol_string(symbol=ERROR_SYMBOL)}'

        return error_text

    def get_transactions_header(self):
        return self.get_sub_header(text='TRANSACTIONS LOG')

    @staticmethod
    def get_transaction(
            node_number, number, host, hash_, from_address, to_address, gas, gas_price, value, status=TX_STATUS_CREATED):
        time = datetime.strftime(datetime.now(), TRANSACTION_DATETIME_FORMAT)

        return f'[{time}] {status} N{node_number}:#{number} HOST: {host} ' \
               f'HASH: {hash_} TX: {from_address} > {to_address} ' \
               f'{{G: {gas}; GP: {gas_price} V: {value}}}'

    @staticmethod
    def get_run_end(node_index, transactions_performed, transactions_expected):
        time = datetime.strftime(datetime.now(), TRANSACTION_DATETIME_FORMAT)

        return f'[{time}] N{node_index + 1} RUN ENDED. PERFORMED {transactions_performed}/{transactions_expected} TXs.'

    def get_run_failed(self, node_index, error_message):
        run_failed_title = self.get_error(f'N{node_index + 1} NODE RUN HAS FAILED! WITH ERROR:')
        run_failed_error = f'\n{error_message}'
        run_failed_error += f'\nself._get_screen_wide_symbol_string(symbol=ERROR_SYMBOL)'

        return f'\n{run_failed_title}{run_failed_error}\n'

    def get_nodes_and_accounts_table(self, accounts_data):
        table = self.get_sub_header('NODES & ACCOUNTS')
        left_column_header = self._get_centered(
            text="NODE", symbol=SUB_HEADER_SYMBOL, width=round(SCREEN_WIDTH * LEFT_COLUMN_SIZE_FACTOR)
        )

        table += f'\n{left_column_header}'

        right_column_header = self._get_centered(
            text="ACCOUNTS", symbol=SUB_HEADER_SYMBOL, width=round(SCREEN_WIDTH * RIGHT_COLUMN_SIZE_FACTOR)
        )

        table += f'{right_column_header}'

        for node_host, node_accounts in accounts_data.items():
            table += f'\n{node_host}\n'

            for account in node_accounts:
                table += ' ' * round(SCREEN_WIDTH * LEFT_COLUMN_SIZE_FACTOR)
                table += '| '
                table += f'{account.get("address")}\n'

        table += self.get_sub_footer()

        return table

    def get_bank_account_info(self, address, balance, stage):
        bank_account_info = self.get_sub_header(f'BANK ACCOUNT STATS ({stage} OF THE RUN)')

        bank_account_info += f'\nADDRESS: {address}'
        bank_account_info += f'\nBALANCE: {balance}\n'

        bank_account_info += self.get_sub_footer()

        return bank_account_info
