"""This module is responsible for validating the upload and marking the content as AVAILABLE in the database"""
import datetime

from common.content_storage_util import storage_path_util
from common.content_storage_util.storage_service import StorageService
from common.logging.custom_logger import CustomLogger
from common.repository.file_upload_repository import FileUploadRepository
from common.util.constants import UPLOAD_STATUS_UPLOADABLE, CONTENT_STATUS_AVAILABLE, DATETIME_FORMAT

log = CustomLogger()


class FileActivatorService(object):
    """A class with functions validating the upload and marking the content as AVAILABLE in the database."""

    def __init__(self):
        self.__storage_service = StorageService()
        self.__file_upload_repository = FileUploadRepository()

    def activate_file(self, event_subject: str) -> None:
        """
        Activate file based on incoming event
        :param event_subject: azure blob trigger event grid subject
        :return: None
        """
        log.debug("Validating the uploaded object properties")
        storage_path = storage_path_util.get_storage_path_from_event_subject(event_subject)
        log.append_key("uploadId", storage_path.upload_id)
        object_properties = self.__storage_service.get_blob_properties(storage_path.default_object_path)
        if not self.__storage_service.is_first_version_of_object(storage_path.default_object_path,
                                                                 object_properties.version_id):
            log.warn("Uploaded object not first version")
            return None
        # Get upload entry for the uploadId
        file_upload_entry = self.__file_upload_repository.get_upload_by_id(storage_path.upload_id)
        # Skip processing in cases like status is ERROR, DONE
        if file_upload_entry.status != UPLOAD_STATUS_UPLOADABLE:
            log.warn(f"Overall status: {file_upload_entry.status} not {UPLOAD_STATUS_UPLOADABLE} for "
                     f"upload_id: {storage_path.upload_id}")
            return None
        log.debug("Activating the file")
        updated_time = datetime.datetime.utcnow().strftime(DATETIME_FORMAT)
        updated_value = file_upload_entry.input.contents[storage_path.default_object_path]
        updated_value['status'] = CONTENT_STATUS_AVAILABLE
        updated_value['size'] = object_properties.size
        updated_value['updatedOn'] = updated_time
        self.__file_upload_repository.update_content_and_status(file_upload_entry,
                                                                content_key_to_update=storage_path.default_object_path,
                                                                updated_value=updated_value, modified_time=updated_time)
        log.debug("File Activated")
