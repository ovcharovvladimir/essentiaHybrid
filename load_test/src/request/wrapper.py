"""
Provide base request.
"""
import json
import requests


class RequestWrapper:
    """
    Request wrapper implementation.
    """

    def __init__(self, url):
        self.url = url

    @staticmethod
    def _wrap_response(response):
        """
        Extract JSON object from responce.
        """
        response_json = json.loads(response.content.decode())
        return response_json.get('result'), response_json.get('error')

    def _send(self, method, **kwargs):
        request_method = getattr(requests, method) if hasattr(requests, method) else None

        if request_method is not None:
            json_string = kwargs.get("json")
            if json_string:
                json_string = json.dumps(json_string).replace('\'', '"')

            print(f'→ Sent {method.upper()} to {self.url} json: {json_string}')
            response = request_method(self.url, timeout=(60, 60), **kwargs)
            print(f'← Received: {response.content}')

            return self._wrap_response(response=response)

        return None

    def send(self, **kwargs):

        return self._send('post', **kwargs)
