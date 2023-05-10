"""This module serves utility function for azure blob storage operations."""
import os
from datetime import datetime
from azure.storage.blob import BlobServiceClient, BlobProperties
from azure.storage.blob import BlobSasPermissions, generate_blob_sas

from common.logging.custom_logger import CustomLogger

CONNECTION_STRING = os.getenv("AzureWebJobsStorage")
STORAGE_ACCOUNT_URL = os.getenv("StorageAccountUrl")
CONTAINER = os.getenv("BlobContainer")

log = CustomLogger()


class StorageService(object):
    """A class with functions to perform operation on azure blob storage."""

    def __init__(self):
        self.__blob_service_client = BlobServiceClient.from_connection_string(CONNECTION_STRING)
        self.__container_client = self.__blob_service_client.get_container_client(CONTAINER)

    def get_blob_properties(self, object_key) -> BlobProperties:
        """ Fetch the given blob from blob properties.
        :param object_key: blob name in blob storage
        :return: blob properties """

        blob_client = self.__container_client.get_blob_client(object_key)
        return blob_client.get_blob_properties()

    def __list_blob_versions(self, object_key):
        """ Fetch the given blob versions from blob storage.
        :param object_key: blob name in blob storage
        :return: blob object versions response """

        log.info(f"Fetching Object versions for:  {object_key}")
        blob_list = self.__container_client.list_blobs(name_starts_with=object_key, include=['versions'])
        versions = []
        for blob_property in blob_list:
            versions.append(blob_property.version_id)

        return versions

    def is_first_version_of_object(self, object_key: str, current_version_id: str) -> bool:
        """
        Return true if given version is the first version else false.
        :param object_key: blob storage blob name
        :return: bool
        """
        versions_resp = self.__list_blob_versions(object_key)
        # validate latest versionsId is the first version of the blob
        return versions_resp[0] == current_version_id

    def generate_pre_signed_url(self, container: str, blob: str, expiry: datetime, upload: bool) -> str:
        """
       Generate pre-signed url for given blob with expiry
       :param container: container name
       :param blob: blob name along with path
       :param expiry: sas token expiry
       :param upload: where the token is for upload
       :return: pre-signed url
        """
        log.debug(f"Generating SAS token upload: {upload}")
        permissions = (
            BlobSasPermissions(create=True, write=True, tag=False)
            if upload
            else BlobSasPermissions(read=True, tag=False)
        )
        sas_token = generate_blob_sas(
            self.__blob_service_client.account_name,
            container,
            blob,
            account_key=self.__blob_service_client.credential.account_key,
            permission=permissions,
            expiry=expiry
        )
        blob_signed_url = f"{STORAGE_ACCOUNT_URL}/{container}/{blob}?{sas_token}"
        log.debug(f"Generated pre-signed url for blob {blob}")
        return blob_signed_url
