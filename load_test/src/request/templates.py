"""
Provide request JSON generation from template.
"""
from time import time
from random import randint
from sys import maxsize

ID_PRECISION_FACTOR = 100000

DEFAULT_REQUEST_JSON_TEMPLATE = {
    'jsonrpc': '2.0',
    'id': 1,
}

id_counter = 0


def _get_random_int():
    """
    Get a random int id from 0 to max int size
    """
    return randint(0, maxsize)


def _get_next_id():
    """
    Get next int id based on current elapsed time.
    """
    return int(time() * ID_PRECISION_FACTOR)


def get_request_json(method, *params, **kwparams):
    """
    Get correct JSON request valid to be send via RPC Api.
    """
    # global id_counter

    request_json = DEFAULT_REQUEST_JSON_TEMPLATE

    parameters = []
    parameters.extend(params)

    # print(f'keys: {kwparams.keys()}')

    if len(kwparams) > 0:
        reformatted_kwparams = {}

        for key in kwparams.keys():
            # print(f'Key: {key}; pos: {key.find("_")}')
            if key.find('_') > 0:
                new_key = key.strip('_')
            else:
                new_key = key

            value = kwparams.get(key)

            reformatted_kwparams.update({
                new_key: value,
            })

        parameters.append(reformatted_kwparams)

    request_json.update({
        'method': method,
        'params': parameters,
        'id': _get_next_id(),
    })

    # id_counter += 1

    return request_json

# get_request_json('test', from_='from')

# WalletTransaction(host).create('0xeb1a4381a68a91ab71c2680e4d18b90ce6e6bdc6', '0x72659ae17432c3f96c5ea896336276795cbde1bd', 1, 1, 1)
