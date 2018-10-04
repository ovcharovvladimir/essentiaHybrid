"""
Provide wei utilities.
"""
import math


def gwei_to_wei(gwei):
    """
    Convert Gwei to Wei.

    Arguments:
        gwei (float): Ethereum amount of Gwei to be converted.

    Returns:
        int: converted value in Wei.
    """
    return int(gwei * math.pow(10, 9))


def wei_to_gwei(wei):
    """
    Convert Wei to Gwei.

    Arguments:
        wei (float): Ethereum amount of Wei to be converted.

    Returns:
        int: converted value in Gwei.
    """
    return int(wei / math.pow(10, 9))


def gwei_to_ether(gwei):

    return int(gwei / math.pow(10, 9))


def ether_to_gwei(ether):

    return int(ether * math.pow(10, 9))
