"""
Provide log functionality.
"""


def log(*args, end='\n'):
    """
    Log wrapper.
    """
    print(*args, end=end)


def log_in_line(*args):
    """
    Log on the same line, no new line break.
    """
    print(*args, end='')

