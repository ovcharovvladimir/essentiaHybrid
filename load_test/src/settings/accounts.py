"""
Provide settings for accounts.
"""
from settings.nodes import (
    GESS_01,
    GESS_02,
    MINER_01,
    MINER_02,
    MINER_03,
    MINER_04,
)

DEFAULT_ACCOUNT_PASSWORD = 'pass'

ACCOUNTS_PER_NODE = 5

BANK_ACCOUNTS = [
    {
        'host': GESS_01,
        'address': '0x369b45aB7795090aCa14209492A5E67868bEd2BF',
        'password': '0987654321',
        'pk': '35808ac1b34728e8e78803c2aef48818c621dee4ecafac207d0409b3044fd24e',
    },
    {
        'host': GESS_01,
        'address': '0x74e863e30899cd336B86EC4aa6903ff6C2fcacD7',
        'password': '0987654321',
        'pk': 'ce9d418df3bcd93bfb9ed7465b8952f78a5d8b7597b510e2d6b1528aa8f5b1d8',
    },
    {
        'host': GESS_01,
        'address': '0xea053c8ccbdf49371a01a95a7bf17f721da1e900',
        'password': '123',
        'pk': '075e40cfd8ca49d760b839fa07397c4c41550497d4058962fc00eaab847eea71',
    },
    {
        'host': GESS_01,
        'address': '0x816d3cf337ade88b1d7f6fcae924a5f71a58fe7c',
        'password': '0987654321',
        'pk': '',
    },
    # {
    #     'host': MINER_01,
    #     'address': '0x14997ad5fbe8e4752e40d850394e35370428f108',  # NO BALANCE
    #     'password': '123',
    # },
    # {
    #     'host': MINER_02,
    #     'address': '0xf098c7dfa65f1551296f15228fa65beb9c9db1d9',
    #     'password': '123',
    # },
    # {
    #     'host': MINER_03,
    #     'address': '0x5d24f0d5ad805ced2b5abf9acec8c30c02ee7d26',  # NO BALANCE
    #     'password': '123',
    # },
    # {
    #     'host': MINER_04,
    #     'address': '0xd3ae3e941114226e0193ca2867f3f4a285084a42',  # NO BALANCE
    #     'password': '123',
    # },
]

BANK_ACCOUNT = BANK_ACCOUNTS[2]
