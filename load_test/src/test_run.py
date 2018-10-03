import multiprocessing

from runner.logger import log
from runner.nodes import Miners, GessNodes


def get_balance_from(node):
    for i in range(1000):
        # log.debug(f'Request #{i}')
        node.wallet_balance.get(address='0x232d5195c74a9c8332b5d24b5e3f5b351099050e')
        # node.wallet_balance.get(address='0x5d24f0d5ad805ced2b5abf9acec8c30c02ee7d26')
        # balance = node.wallet_balance.get(address='0x5d24f0d5ad805ced2b5abf9acec8c30c02ee7d26')

        # print(f'Got balance from {node.host}: {balance}')


def separate_logger(number):
    for i in range(1000):
        log.info(f'Process #{number} :: Request #{i}')


if __name__ == '__main__':
    # pool = multiprocessing.Pool(4)
    #
    # pool.map(get_balance_from, [miner for miner in Miners()])

    # log.info(f'----- New session start -----')
    # pool = multiprocessing.Pool(4)
    # pool.map(separate_logger, [i + 1 for i in range(4)])
    #
    # exit(0)


    #
    # Miners
    #
    # miner = Miners()[0]
    #
    # log.debug(miner.host)
    #
    # for i in range(1000):
    #     log.debug(f'Request #{i}')
    #     miner.wallet_balance.get(address='0x5d24f0d5ad805ced2b5abf9acec8c30c02ee7d26')

    #
    # Gess
    #
    nodes = GessNodes()

    # print(nodes)
    # import pdb; pdb.set_trace()

    pool = multiprocessing.Pool(len(nodes))
    pool.map(get_balance_from, [gess_node for gess_node in nodes])

