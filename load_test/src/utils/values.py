"""
Provide values utils.
"""


def hex_to_int(hex_):
    """
    Convert hex string to int
    """
    return int(hex_, 0)


def clamp(x, mn, mx):
    """
    Clamp x value between minimal and maximal values.
    """
    return min(max(x, mn), mx)


def add_hex_prefix(data):
    """
    Add simple '0x' prefix to a given data.
    """
    if data[:2] == '0x':
        return data

    return f'0x{data}'


def strip_hex_prefix(data):
    """
    Remove '0x' prefix from a given data.
    """
    if data[:2] == '0x':
        return data[2:]

    return data
