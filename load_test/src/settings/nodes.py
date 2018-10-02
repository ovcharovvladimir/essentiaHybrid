"""
Provide nodes settings.
"""
from utils.cycle_list import CycleList

GESS_01 = '18.224.0.169'
GESS_02 = '52.14.180.128'   # connection refused
GESS_03 = '18.221.62.255'
GESS_04 = '52.14.5.83'
GESS_05 = '18.219.132.34'

MINER_01 = '18.222.107.145'
MINER_02 = '18.219.184.139'
MINER_03 = '18.218.220.164'
MINER_04 = '18.220.24.83'

DEFAULT_PORT = 8545


NODES_HOSTS = CycleList()
NODES_HOSTS.extend([
    GESS_01,
    GESS_02,
    GESS_03,
    GESS_04,
    GESS_05,
])

MINERS_HOSTS = CycleList()
MINERS_HOSTS.extend([
    MINER_01,
    MINER_02,
    MINER_03,
    MINER_04,
])


def get_node_url(node_host=None, index=None):
    """
    Return formatted node url.
    """
    if index:
        node_host = NODES_HOSTS[index]

    return f'http://{node_host}:{DEFAULT_PORT}'
