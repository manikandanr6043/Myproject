"""This module consists of functions that help in making http calls using request library."""
import inspect

import requests
from requests.models import Response

from common.logging.custom_logger import CustomLogger

log = CustomLogger()


# Custom Exception class for raising exceptions in ApiClient class
class ApiError(Exception):
    """Class that represents the exception thrown on API call error."""

    def __init__(self, response_code, message="API call failed", *, response=None):
        self.response_code = response_code
        self.response_body = response.text if response is not None else response
        self.message = message
        self.response = response
        super().__init__(self.message)

    def __str__(self):
        return self.message


class ApiClient(object):
    """Performs http(s) calls and throws ApiError in case response code is not < 300."""

    @staticmethod
    def __make_request(method, endpoint, headers=None, json=None, *, data=None) -> Response:
        """
        Perform HTTP calls.
        :param method: HTTP method
        :param endpoint: API endpoint
        :param headers: Request headers
        :param json: request json
        :param data: request body
        :return: Response
        """
        log.info(f'API ======== {method} {endpoint}')
        response: Response = requests.request(method, endpoint, headers=headers, json=json, data=data)
        response_code = response.status_code
        # raise ApiError if response code is not less than 300
        if response_code > 299:
            raise ApiError(response_code, f"API call failed with response code: {response_code}", response=response)
        return response

    @staticmethod
    def make_request(method, endpoint, headers=None, json=None, *, data=None) -> Response:
        """
        Overloaded method to log the key value pair based logs for downstream http calls.
        Calls the __make_request method to perform HTTP calls.
        :param method: HTTP method
        :param endpoint: API endpoint
        :param headers: Request headers
        :param json: request json
        :param data: request body
        :return: Response
        """
        try:
            response = ApiClient.__make_request(method, endpoint, headers, json, data=data)
            log_downstream_call(method, endpoint, response)
        except ApiError as ae:
            response = ae.response
            log_downstream_call(method, endpoint, response)
            raise ae
        return response


def log_downstream_call(http_method, api_url, response: Response) -> None:
    """
    Return fields to be appended to downstream call log as dict.
    :param http_method: HTTP method
    :param api_url: API URL
    :param response: Response of the API call
    """
    # Get Caller Information
    curr_frame_f_back = inspect.currentframe().f_back.f_back
    caller_method_name = curr_frame_f_back.f_code.co_name
    fields = {
        "api": f'{convert_snake_case_to_camel_case(caller_method_name)}',
        "httpMethod": http_method,
        "url": api_url,
        "httpStatus": response.status_code,
        "elapsedTime": response.elapsed.total_seconds(),
        "header": response.headers if response.status_code > 299 else None,
        "responseBody": response.text if response.status_code > 299 else None
    }
    log.info(f"{caller_method_name} Response", extra=fields)


def convert_snake_case_to_camel_case(input_str: str) -> str:
    """
    Convert Snake Case String to Camel Case.
    Using split() + join() + title() + generator expression.
    :param input_str: input string in snake case
    :return: string in camel case.
    """
    # split underscore using split
    temp_str = input_str.split('_')
    # joining result
    return temp_str[0] + ''.join(ele.title() for ele in temp_str[1:])
