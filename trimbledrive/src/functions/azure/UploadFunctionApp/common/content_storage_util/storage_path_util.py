"""This module consists of storage path utility."""
import dataclasses

from common.util import constants


@dataclasses.dataclass(eq=True, frozen=True)
class StoragePath(object):
    """A class that represents storage path components."""
    space_id: str
    file_id: str
    upload_id: str
    tmp_object_path: str
    default_object_path: str


def get_storage_path_from_event_subject(event_subject: str) -> StoragePath:
    """
    Get StoragePath from azure blob created event grid event subject
    :param event_subject: azure blob created event grid event subject
    :return: StoragePath
    """
    # Subject: /blobServices/default/containers/{containerName}/blobs/tmp/{spaceId}/{fileId}/{uploadId}/{contentName}
    truncated_subject = event_subject.replace("/blobServices/default/containers/", "")
    storage_path_list = truncated_subject.split(constants.STORAGE_PATH_SEPARATOR)
    container_name = storage_path_list[0]
    blob_path_prefix = f"{container_name}/blobs/"
    if constants.THUMB_STORAGE_DIR + constants.STORAGE_PATH_SEPARATOR in truncated_subject:
        blob_path_prefix = blob_path_prefix + constants.THUMB_STORAGE_DIR + constants.STORAGE_PATH_SEPARATOR
    else:
        blob_path_prefix = blob_path_prefix + constants.DEFAULT_STORAGE_DIR + constants.STORAGE_PATH_SEPARATOR
    blob_path = truncated_subject.replace(blob_path_prefix, "")
    blob_path_list = blob_path.split(constants.STORAGE_PATH_SEPARATOR)
    space_id = blob_path_list[0]
    file_id = blob_path_list[1]
    upload_id = blob_path_list[2]
    object_tmp_path = constants.THUMB_STORAGE_DIR + constants.STORAGE_PATH_SEPARATOR + blob_path
    object_default_path = constants.DEFAULT_STORAGE_DIR + constants.STORAGE_PATH_SEPARATOR + blob_path
    return StoragePath(space_id=space_id, file_id=file_id, upload_id=upload_id, tmp_object_path=object_tmp_path,
                       default_object_path=object_default_path)
