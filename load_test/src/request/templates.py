DEFAULT_REQUEST_JSON_TEMPLATE = {
    'jsonrpc': '2.0',
    'id': 1,
}


def get_request_json(method, *params, **kwparams):
    """
    Get correct JSON request valid to be send via RPC Api.
    """
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
    })

    return request_json

# get_request_json('test', from_='from')

# WalletTransaction(host).create('0xeb1a4381a68a91ab71c2680e4d18b90ce6e6bdc6', '0x72659ae17432c3f96c5ea896336276795cbde1bd', 1, 1, 1)
