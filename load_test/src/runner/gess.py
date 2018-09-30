"""
Provide nodes.
"""
from services.node import Node
from settings.nodes import (
    get_node_url,
    NODES_HOSTS,
    MINERS_HOSTS,
)
from utils.cycle_list import CycleList


class GenericNodes(CycleList):
    """
    Generic nodes list implementation.
    """

    __instance = None

    def __new__(cls):
        if not cls.__instance:
            cls.__instance = CycleList.__new__(cls)

        return cls.__instance

    def __init__(self, *args, nodes=[], **kwargs):
        super(GenericNodes, self).__init__(*args, **kwargs)

        self.extend([Node(host=get_node_url(node_host)) for node_host in nodes])


class GessNodes(GenericNodes):
    """
    Gess nodes list implementation.
    """

    def __init__(self, *args, nodes=NODES_HOSTS, **kwargs):
        super(GessNodes, self).__init__(*args, nodes=nodes, **kwargs)


class Miners(GenericNodes):
    """
    Miners list implementation.
    """

    def __init__(self, *args, nodes=MINERS_HOSTS, **kwargs):
        super(Miners, self).__init__(*args, nodes=nodes, **kwargs)
