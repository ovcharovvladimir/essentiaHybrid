"""
Provide nodes settings.
"""
from utils.cycle_list import CycleList

# '18.224.218.68 18.224.50.75 52.15.181.235 18.188.111.198 18.188.240.197 18.221.62.255 18.224.11.186 18.224.106.72 18.224.121.61 18.224.159.84 18.224.168.178 18.224.198.158'

BOOT_01 = '18.224.218.68'
BOOT_02 = '18.224.50.75'
BOOT_03 = '52.15.181.235'

GESS_01 = '18.188.111.198'
GESS_02 = '18.188.240.197'
GESS_03 = '18.221.62.255'
GESS_04 = '18.224.11.186'
GESS_05 = '18.224.106.72'

MINER_01 = '18.224.121.61'
MINER_02 = '18.224.159.84'
MINER_03 = '18.224.168.178'
MINER_04 = '18.224.198.158'

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
