"""
Provide settings for transaction.
"""
from utils.wei import gwei_to_wei

TRANSACTION_GAS = 200000
TRANSACTION_GAS_PRICE = gwei_to_wei(200)
TRANSACTION_VALUE = gwei_to_wei(1)

MAX_SEND_RAW_RETRIES_COUNT = 5
MAX_SEND_RETRIES_COUNT = 5
