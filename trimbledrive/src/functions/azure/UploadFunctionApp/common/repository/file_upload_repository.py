"""This module consists of functions to perform DB calls to CosmosDB."""
import os
import json
from datetime import datetime
from time import mktime
from wsgiref.handlers import format_date_time
import base64
from urllib.parse import quote
import hmac
from hashlib import sha256
from typing import Dict, Optional

from azure.cosmos import CosmosClient

from common.logging.custom_logger import CustomLogger
from common.util.api_client import ApiClient
from common.repository.entity.file_upload import FileUpload

from pyckson import set_defaults, explicit_nulls, parse

set_defaults(explicit_nulls)

log = CustomLogger()
MASTER_KEY = os.getenv("CosmosKey")
COSMOSDB_NAME = os.getenv("CosmosAccountName")
DATABASE_NAME = os.getenv("CosmosDbName")
CONTAINER_NAME = "file_upload"
ENDPOINT = f"https://{COSMOSDB_NAME}.documents.azure.com:443/"


class FileUploadRepository(object):
    """A class with functions to perform operations related to CosmosDB."""

    def __init__(self) -> None:
        self.__client = CosmosClient(ENDPOINT, credential=MASTER_KEY)
        self.__database = self.__client.get_database_client(DATABASE_NAME)
        self.__container = self.__database.get_container_client(CONTAINER_NAME)

    def get_upload_by_id(self, upload_id: str) -> FileUpload:
        """
        Get File Upload Details by Upload ID.
        :param upload_id: upload identifier.
        :return: FileUpload
        """
        log.debug(f"Fetching upload details from DB: {upload_id}")
        file_upload_dict = self.__container.read_item(item=upload_id, partition_key=upload_id)
        file_upload_entry = parse(FileUpload, file_upload_dict)
        return file_upload_entry

    @staticmethod
    def update_content_and_status(file_upload_entry: FileUpload, content_key_to_update: str, updated_value: Dict,
                                  *, upload_status: Optional[str] = None, error_reason: Optional[str] = None,
                                  modified_time: str) -> None:
        """
        Update File Upload details like content by key, overall status and error reason.
        :param: file_upload_entry entry obtained by get_upload_by_id
        :param: content_key_to_update in the content map attribute key of the content to be updated
        :param: updated_value value to be updated for the given content attribute key
        :param: upload_status optional, value of status attribute to be updated
        :param: error_reason optional, value for error_reason attribute
        :return: None
        """
        log.debug(f"Updating upload contents for the upload [uploadId={file_upload_entry.id},"
                  f"content_key_to_update={content_key_to_update}, updated_value={updated_value}, "
                  f"upload_status={upload_status}, error_reason={error_reason}]")
        partition_key = f'["{file_upload_entry.id}"]'

        # MAX 10 OPERATIONS PER REQUEST

        partial_update_operations = [
            {
                "op": "set",
                "path": f"/input/contents/{content_key_to_update.replace('/', '~1')}",
                "value": updated_value
            },
            {
                "op": "set",
                "path": "/modifiedOn",
                "value": modified_time
            }
        ]

        if upload_status is not None:
            status_op = {
                "op": "set",
                "path": "/status",
                "value": upload_status
            }
            partial_update_operations.append(status_op)
        if error_reason is not None:
            error_reason_op = {
                "op": "add",
                "path": "/errorReason",
                "value": error_reason
            }
            partial_update_operations.append(error_reason_op)

        payload = {
            "operations": partial_update_operations
        }

        FileUploadRepository.patch_document(file_upload_entry.id, partition_key, payload)

    @staticmethod
    def generate_master_key_signature(doc_id, verb, date) -> str:
        """
        Generate the master key signature for REST API calls to cosmos DB
        :param: doc_id document id
        :param: verb HTTP verb
        :param: date
        :return: url safe master key signature
        """
        log.debug("Generating master key signature")
        key_type = 'master'
        version = '1.0'

        resource_type = 'docs'
        resource_id = f'dbs/{DATABASE_NAME}/colls/{CONTAINER_NAME}/docs/{doc_id}'

        text = "{verb}\n{resource_type}\n{resource_id}\n{date}\n\n".format(
            verb=(verb.lower() or ''),
            resource_type=(resource_type.lower() or ""),
            resource_id=(resource_id or ""),
            date=date.lower())

        digest = hmac.new(base64.b64decode(MASTER_KEY), text.encode("utf-8"), sha256).digest()
        signature = base64.encodebytes(digest).decode("utf-8")
        uri = f'type={key_type}&ver={version}&sig={signature[:-1]}'
        uri_encoded = quote(uri)

        return uri_encoded

    @staticmethod
    def patch_document(doc_id, partition_key, payload) -> None:
        """
        Perform patch document on the given document
        :param: partition_key
        :param: payload
        :return: None
        """
        log.debug("Performing patch document")
        now = datetime.now()
        stamp = mktime(now.timetuple())
        date = format_date_time(stamp)

        rest_api_version = '2018-12-31'

        # Using REST API for patch document as cosmos python sdk does not support patch document currently

        url = f"{ENDPOINT}/dbs/{DATABASE_NAME}/colls/{CONTAINER_NAME}/docs/{doc_id}"

        headers = {
            'Content-Type': 'application/json_patch+json',
            'authorization': FileUploadRepository.generate_master_key_signature(doc_id, 'patch', date),
            'x-ms-date': date,
            'x-ms-version': rest_api_version,
            'x-ms-documentdb-partitionkey': partition_key
        }

        log.debug(f"patch document payload {payload}")

        ApiClient.make_request("PATCH", url, headers=headers, data=json.dumps(payload, default=str))
