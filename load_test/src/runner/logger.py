"""
Provide logger for test runner.
"""
import logging

from settings.logger import (
    LOGGER_OUTPUT_FILE,
    LOGGER_TITLE,
)


def __create_logger():
    """
    Create test runner logger.
    """
    logger = logging.getLogger(LOGGER_TITLE)
    logger.setLevel(logging.DEBUG)

    file_handler = logging.FileHandler(LOGGER_OUTPUT_FILE)
    file_handler.setLevel(logging.INFO)

    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.DEBUG)

    logger.addHandler(file_handler)
    logger.addHandler(console_handler)

    return logger


log = __create_logger()
