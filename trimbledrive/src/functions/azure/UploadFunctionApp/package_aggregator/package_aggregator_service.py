"""This module consists of class with functions to perform PackageAggregator Business logic."""
import json
import uuid
from datetime import datetime
from typing import Optional

from common.logging.custom_logger import CustomLogger
from common.util.constants import UPLOAD_STATUS_UPLOADABLE, CONTENT_STATUS_AVAILABLE

log = CustomLogger()


class PackageAggregatorService(object):
    """A class with functions to perform PackageAggregator Business logic."""

    @staticmethod
    def process_doc(doc: dict) -> Optional[str]:
        """
        Process the upload document and return the commit message to publish if upload is ready for commit
        """
        # add uploadId to custom dimensions
        log.append_key("uploadId", doc['id'])
        commit_processor_msg = None
        if doc['status'] == UPLOAD_STATUS_UPLOADABLE and PackageAggregatorService.all_parts_available(doc):
            log.info("Package ready for commit")
            commit_processor_msg = PackageAggregatorService.generate_commit_processor_msg(doc)
        # remove uploadId from logging context
        log.remove_key("uploadId")
        return commit_processor_msg

    @staticmethod
    def all_parts_available(upload_doc: dict) -> bool:
        """
        Aggregates the overall status of the package content and determine if the upload is commit ready.
        :param upload_doc: upload document as dict
        :return True|False
        """
        # Check if all parts are available
        return all(v.get("status", None) == CONTENT_STATUS_AVAILABLE
                   for k, v in upload_doc['input']['contents'].items())

    @staticmethod
    def generate_commit_processor_msg(upload_doc: dict) -> str:
        """
        Generate message to publish for Commit Processor Topic
        :param upload_doc: upload document as dict
        """
        # Message body to be sent to commit processor topic for further processing
        upload_id = upload_doc['id']
        upload_input = upload_doc['input']
        space_id = upload_input['spaceId']
        file_id = upload_input['fileId']
        message_dict = {
            "specversion": "1.0",
            "id": str(uuid.uuid4()),
            "type": "com.trimble.tdrive.upload_commit_ready.v1",
            "subject": f"{space_id}#{file_id}#{upload_id}",
            "time": datetime.utcnow().isoformat(),
            "datacontenttype": "application/json",
            "data": {
                "spaceId": space_id,
                "fileId": file_id,
                "uploadId": upload_id
            }
        }
        log.debug(f"Message with id:{message_dict['id']} created for the commit processor topic")
        return json.dumps(message_dict)
