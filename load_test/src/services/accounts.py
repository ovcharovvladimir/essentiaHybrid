"""
Provide accounts functionality.
"""
from request.wrapper import RequestWrapper
from request.templates import get_request_json
from settings.accounts import DEFAULT_ACCOUNT_PASSWORD

DEFAULT_UNLOCK_TIME = 1 * 60 * 60   # 1 HOUR (seconds)


class FailedToCreateAccount(Exception):
    """
    Failed to create account error.
    """

    pass


class FailedToGetAccountsList(Exception):
    """
    Failed to get accounts list error.
    """

    pass


class FailedToUnlockAccount(Exception):
    """
    Failed to unlock account error.
    """

    pass


class Account(RequestWrapper):
    """
    Accounts implementations.
    """

    def create(self, password=DEFAULT_ACCOUNT_PASSWORD):
        """
        Create a new account and return it's address.
        """
        new_address, error = self.send(json=get_request_json('personal_newAccount', password))

        if error:
            raise FailedToCreateAccount(error)

        if new_address is None:
            raise FailedToCreateAccount

        return new_address

    def get_all(self):
        """
        Get list of all host accounts.
        """
        accounts_list, error = self.send(json=get_request_json('eth_accounts'))

        if error:
            raise FailedToGetAccountsList(error)

        if accounts_list is None:
            raise FailedToGetAccountsList

        return accounts_list

    def unlock(self, address, password=DEFAULT_ACCOUNT_PASSWORD, time=DEFAULT_UNLOCK_TIME):
        unlock_status, error = self.send(json=get_request_json('personal_unlockAccount', address, password, time))

        if error:
            raise FailedToUnlockAccount(error)

        if unlock_status is None:
            raise FailedToUnlockAccount

        elif not unlock_status:
            raise FailedToUnlockAccount('Node returned "False"!')

        return unlock_status
