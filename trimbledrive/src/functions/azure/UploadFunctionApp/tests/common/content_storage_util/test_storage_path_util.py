import uuid
from unittest import TestCase

from common.content_storage_util.storage_path_util import get_storage_path_from_event_subject, StoragePath

TEST_SPACE_ID = str(uuid.uuid4())
TEST_FILE_ID = str(uuid.uuid4())
TEST_UPLOAD_ID = str(uuid.uuid4())
TEST_CONTENT = "testContent"
TEST_TMP_PATH = f"thumb/{TEST_SPACE_ID}/{TEST_FILE_ID}/{TEST_UPLOAD_ID}/{TEST_CONTENT}"
TEST_DEFAULT_PATH = f"orig/{TEST_SPACE_ID}/{TEST_FILE_ID}/{TEST_UPLOAD_ID}/{TEST_CONTENT}"


class Test(TestCase):
    """Class for testing storage_path_util module."""
    def test_get_storage_path_from_event_subject_with_tmp_path(self):
        test_tmp_path = f"/blobServices/default/containers/test-container/blobs/thumb/{TEST_SPACE_ID}/{TEST_FILE_ID}" \
                        f"/{TEST_UPLOAD_ID}/{TEST_CONTENT}"
        storage_path = get_storage_path_from_event_subject(test_tmp_path)
        self.validate_storage_path(storage_path)

    def test_get_storage_path_from_event_subject_with_default_path(self):
        test_tmp_path = f"/blobServices/default/containers/test-container/blobs/orig/{TEST_SPACE_ID}/{TEST_FILE_ID}" \
                        f"/{TEST_UPLOAD_ID}/{TEST_CONTENT}"
        storage_path = get_storage_path_from_event_subject(test_tmp_path)
        self.validate_storage_path(storage_path)

    def validate_storage_path(self, storage_path: StoragePath):
        self.assertTrue(storage_path.space_id == TEST_SPACE_ID)
        self.assertTrue(storage_path.file_id == TEST_FILE_ID)
        self.assertTrue(storage_path.upload_id == TEST_UPLOAD_ID)
        self.assertTrue(storage_path.tmp_object_path == TEST_TMP_PATH)
        self.assertTrue(storage_path.default_object_path == TEST_DEFAULT_PATH)
