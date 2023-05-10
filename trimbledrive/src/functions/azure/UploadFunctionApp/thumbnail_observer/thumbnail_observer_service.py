"""This module serves function to generate the message for thumbnail processor."""
import json
from datetime import datetime, timedelta
from typing import Optional

from common.content_storage_util import storage_path_util, storage_service
from common.content_storage_util.storage_service import StorageService
from common.logging.custom_logger import CustomLogger
from common.repository.file_upload_repository import FileUploadRepository
from common.util.constants import UPLOAD_STATUS_UPLOADABLE, CONTENT_STATUS_ERROR, UPLOAD_STATUS_ERROR, DATETIME_FORMAT

log = CustomLogger()
# Maximum supported thumbnail image size in bytes
MAX_THUMB_SIZE = 2097152  # 2MB


class ThumbnailObserverService(object):
    """A class with functions to generate message for thumbnail processor."""

    def __init__(self):
        self.__storage_service = StorageService()
        self.__file_upload_repository = FileUploadRepository()

    def generate_thumb_processor_msg(self, event_id: str, event_subject: str) -> Optional[str]:
        """
        Generated message to be sent to thumbnail processor topic
        :param event_id: event identifier
        :param event_subject: azure blob trigger event grid subject
        :return: message for thumb processor topic
        """
        log.debug("Validating the uploaded object properties")
        storage_path = storage_path_util.get_storage_path_from_event_subject(event_subject)
        log.append_key("uploadId", storage_path.upload_id)
        output_message = None
        object_properties = self.__storage_service.get_blob_properties(storage_path.tmp_object_path)
        if not self.__storage_service.is_first_version_of_object(storage_path.tmp_object_path,
                                                                 object_properties.version_id):
            log.warn("Uploaded object not first version")
            return None
        # Get upload entry for the uploadId
        file_upload_entry = self.__file_upload_repository.get_upload_by_id(storage_path.upload_id)
        # Skip processing in cases like status is ERROR, DONE
        if file_upload_entry.status != UPLOAD_STATUS_UPLOADABLE:
            log.warn(f"Overall status: {file_upload_entry.status} not {UPLOAD_STATUS_UPLOADABLE} for "
                     f"upload_id: {storage_path.upload_id}")
        elif object_properties.size > MAX_THUMB_SIZE:
            log.warn("Object size greater than max supported size")
            updated_time = datetime.utcnow().strftime(DATETIME_FORMAT)
            updated_value = file_upload_entry.input.contents[storage_path.tmp_object_path]
            updated_value['status'] = CONTENT_STATUS_ERROR
            updated_value['updatedAt'] = updated_time
            self.__file_upload_repository. \
                update_content_and_status(file_upload_entry, content_key_to_update=storage_path.tmp_object_path,
                                          updated_value=updated_value, upload_status=UPLOAD_STATUS_ERROR,
                                          error_reason="ThumbMaxSizeExceeded", modified_time=updated_time)
        else:
            log.debug("Generating Thumb Processor Message")
            url_expiry = datetime.utcnow() + timedelta(days=4)
            message_dict = {
                "specversion": "1.0",
                "id": event_id,
                "type": "com.trimble.tdrive.file_thumb_uploaded.v1",
                "subject": f"{storage_path.space_id}#{storage_path.file_id}#{storage_path.upload_id}",
                "time": datetime.utcnow().isoformat(),
                "datacontenttype": "application/json",
                "data": {
                    "spaceId": storage_path.space_id,
                    "fileId": storage_path.file_id,
                    "uploadId": storage_path.upload_id,
                    "downloadUrl": self.__storage_service.generate_pre_signed_url(storage_service.CONTAINER,
                                                                                  storage_path.tmp_object_path,
                                                                                  url_expiry, False),
                    "uploadUrl": self.__storage_service.generate_pre_signed_url(storage_service.CONTAINER,
                                                                                storage_path.default_object_path,
                                                                                url_expiry, True)
                }
            }
            output_message = json.dumps(message_dict)
            log.debug(f"Generated Thumb Processor Message with id: {message_dict['id']}")
        return output_message
