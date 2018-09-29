"""
Provide settings for accounts.
"""
import requests
from settings.nodes import (
    GESS_01,
    GESS_02,
    GESS_03,
    GESS_04,
    # GESS_05,
    MINER_01,
    MINER_02,
    MINER_03,
    MINER_04,
)

ACCOUNTS_PER_NODE = 5

BANK_ACCOUNTS = [
    {
        'host': MINER_01,
        'address': '0x14997ad5fbe8e4752e40d850394e35370428f108',
        'password': '123',
    },
    {
        'host': MINER_02,
        'address': '0xf098c7dfa65f1551296f15228fa65beb9c9db1d9',
        'password': '123',
    },
    {
        'host': MINER_03,
        'address': '0x5d24f0d5ad805ced2b5abf9acec8c30c02ee7d26',
        'password': '123',
    },
    {
        'host': MINER_04,
        'address': '0xd3ae3e941114226e0193ca2867f3f4a285084a42',
        'password': '123',
    },
]

BANK_ACCOUNT = BANK_ACCOUNTS[0]