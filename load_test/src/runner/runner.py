"""
Provide main runner functionality.
"""
from multiprocessing import Pool

from runner.single_node_run import SingleNodeRun

_load_factor = 0


def _spawn_single_run(node_index):
    """
    Spawn a single node load test runner.
    """
    global _load_factor

    SingleNodeRun(node_index=node_index, load_factor=_load_factor).run()


def run(nodes_count, load_factor):
    """
    Spawn load tests on given nodes count with a given load factor.
    """
    global _load_factor

    _load_factor = load_factor
    pool = Pool()

    pool.map(func=_spawn_single_run, iterable=[node_index for node_index in range(nodes_count)])
