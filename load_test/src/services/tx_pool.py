"""
Provide transactions pool implementation.
"""
from request.templates import get_request_json
from request.wrapper import RequestWrapper


class FailedToGetTxPoolContent(Exception):
    """
    Failed to get tx pool content exception.
    """
    pass


class FailedToGetTxPoolInspectInfo(Exception):
    """
    Failed to get tx pool inspect exception.
    """
    pass


class FailedToGetTxPoolStatus(Exception):
    """
    Failed to get tx pool content exception.
    """
    pass


class TransactionPool(RequestWrapper):
    """
    Transaction pool implementation.
    """

    def get_content(self):
        """
        Get content of tx pool.
        """
        tx_pool_content, error = self.send(json=get_request_json('txpool_content'))

        if error:
            raise FailedToGetTxPoolContent(error)

        if tx_pool_content is None:
            raise FailedToGetTxPoolContent

        return tx_pool_content

    def get_inspect_info(self):
        """
        Get inspect info of tx pool.
        """
        tx_pool_inspect_info, error = self.send(json=get_request_json('txpool_inspect'))

        if error:
            raise FailedToGetTxPoolInspectInfo(error)

        if tx_pool_inspect_info is None:
            raise FailedToGetTxPoolInspectInfo

        return tx_pool_inspect_info

    def get_status(self):
        """
        Get status of tx pool.
        """
        tx_pool_status, error = self.send(json=get_request_json('txpool_inspect'))

        if error:
            raise FailedToGetTxPoolStatus(error)

        if tx_pool_status is None:
            raise FailedToGetTxPoolStatus

        return tx_pool_status
