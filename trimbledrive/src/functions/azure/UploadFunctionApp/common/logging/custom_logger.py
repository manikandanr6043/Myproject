import logging as log
import os
from typing import Dict, Optional

# for more information: https://pypi.org/project/opencensus-ext-azure/
from opencensus.ext.azure.log_exporter import AzureLogHandler

APPINSIGHTS_INSTRUMENTATION_KEY = os.getenv("APPINSIGHTS_INSTRUMENTATIONKEY")


# make sure we can have a single instance of CustomLogger across multiple files/modules (
# https://stackoverflow.com/questions/42237752/single-instance-of-class-in-python)
def singleton(cls, *args, **kwargs):
    """
    To Make sure we can have a single instance of CustomLogger across multiple files/modules (
    https://stackoverflow.com/questions/42237752/single-instance-of-class-in-python).
    """
    instances = {}

    def _singleton(*arg, **kw):
        """Private function to create class instance."""
        if cls not in instances:
            instances[cls] = cls(*arg, **kw)
        return instances[cls]

    return _singleton


@singleton
class CustomLogger(object):
    """This class consists of logging utility functions with azure app insights support"""
    def __init__(self):
        self.logger = None
        self.__context = {}
        self.set_logger("TDrive")
        self.set_level(os.getenv("LOG_LEVEL", "DEBUG"))

    # logging methods
    def info(self, message, *, extra: Optional[Dict] = None):
        properties = self.__build_context(extra)
        self.logger.info(message, extra=properties)

    def warn(self, message, *, extra: Optional[Dict] = None):
        properties = self.__build_context(extra)
        self.logger.warning(message, extra=properties)

    def debug(self, message, *, extra: Optional[Dict] = None):
        properties = self.__build_context(extra)
        self.logger.debug(message, extra=properties)

    def exception(self, message, *, extra: Optional[Dict] = None):
        properties = self.__build_context(extra)
        self.logger.exception(message, extra=properties)

    def error(self, message, *, extra: Optional[Dict] = None):
        properties = self.__build_context(extra)
        self.logger.error(message, extra=properties)

    # public helper methods
    def clear_context(self):
        self.__context.clear()

    def add_keys_from_context(self, context):
        self.clear_context()
        self.append_keys({"invocationId": context.invocation_id, "functionName": context.function_name})

    def append_key(self, key, value):
        self.__context[key] = value

    def append_keys(self, keys: Dict):
        self.__context.update(keys)

    def remove_key(self, key):
        self.__context.pop(key, None)

    # private helper methods
    def __build_context(self, extra):
        # WARNING: For the extra dimension to work, we need to pass a dictionary to the custom_dimensions field (
        # https://pypi.org/project/opencensus-ext-azure/)
        properties = {}
        if extra:
            for key in extra:
                properties[key] = extra[key]
        return {
            'custom_dimensions': {**properties, **self.__context}
        }

    def set_level(self, level):
        self.logger.setLevel(level=level)

    def set_logger(self, name):
        self.logger = log.getLogger(name=name)
        if len(self.logger.handlers) == 0:
            connection_string = 'InstrumentationKey=' + APPINSIGHTS_INSTRUMENTATION_KEY
            self.logger.addHandler(AzureLogHandler(connection_string=connection_string))
