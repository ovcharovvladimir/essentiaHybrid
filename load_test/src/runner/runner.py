"""
Provide main runner functionality.
"""
from multiprocessing import Pool

from runner.nodes import GessNodes
from runner.single_node_run import SingleNodeRun
from utils.cycle_list import CycleList
from utils.values import clamp

# _load_factor = 0
# _nodes_count = 0
# _gess_nodes_for_the_run = GessNodes()


def _spawn_single_run(node_index, load_factor, gess_nodes_for_the_run):
    """
    Spawn a single node load test runner.
    """
    # global _load_factor
    # global _gess_nodes_for_the_run

    # print(f'#{node_index}: {gess_nodes_for_the_run}')

    SingleNodeRun(node_index=node_index, gess_nodes=gess_nodes_for_the_run, load_factor=load_factor).run()


def run(nodes_count, load_factor):
    """
    Spawn load tests on given nodes count with a given load factor.
    """
    # global _load_factor
    # global _nodes_count
    # global _gess_nodes_for_the_run

    _load_factor = load_factor
    _nodes_count = nodes_count
    pool = Pool(nodes_count)

    minimal_nodes_count = 1
    maximal_nodes_count = len(GessNodes()) - 1

    nodes_count = clamp(nodes_count, minimal_nodes_count, maximal_nodes_count)
    gess_nodes_for_the_run = CycleList(GessNodes()[:_nodes_count])

    # pool.map(func=_spawn_single_run, iterable=[node_index for node_index in range(nodes_count)])
    for i in range(nodes_count):
        pool.apply_async(
            func=_spawn_single_run,
            args=(i, load_factor, gess_nodes_for_the_run),
        )

    pool.close()
    pool.join()


# def test_run(nodes_count, load_factor):
#     _spawn_single_run(node_index=0, load_factor=load_factor, gess_nodes_for_the_run=CycleList(GessNodes()[:2]))