"""
Provide implementation of block.
"""
from request.wrapper import RequestWrapper
from request.templates import get_request_json


class FailedToGetBlockInfo(Exception):
    """
    Failed to get block info exception.
    """
    pass


class Block(RequestWrapper):
    """
    Block implementation.
    """

    def get_by(self, tag=None, number=None, full_objects=True):
        """
        Return block by identifier
        """
        param = None

        if tag is not None:
            param = tag

        elif number is not None:
            param = number

        block_info, error = self.send(json=get_request_json('eth_getBlockByNumber', param, full_objects))

        if error:
            raise FailedToGetBlockInfo(error)

        if block_info is None:
            raise FailedToGetBlockInfo

        return block_info

    def get_latest(self):
        """
        Return latest block.
        """
        pass

    def get_latest_block_gas_limit(self):
        """
        Return latest block gasLimit value
        """
        pass
