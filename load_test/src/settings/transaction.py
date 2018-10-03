"""
Provide settings for transaction.
"""
from utils.wei import gwei_to_wei

TRANSACTION_GAS = 100000
TRANSACTION_GAS_PRICE = gwei_to_wei(20)
TRANSACTION_VALUE = gwei_to_wei(1)