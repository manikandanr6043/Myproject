"""
This module consists of the main function.
The main method is the method in your function code that processes events.
This function validates if the upload is ready for commit and publishes the message to commit processor topic.
"""
import json
from typing import List

import azure.functions as func

from common.logging.custom_logger import CustomLogger
from package_aggregator.package_aggregator_service import PackageAggregatorService

package_aggregator_service = PackageAggregatorService()
log = CustomLogger()


def main(documents: func.DocumentList, message: func.Out[str], context: func.Context):
    """
    Changes in document feed will trigger this function with the changed documents.
    This will publish a message in the CommitProcessor Topic.
    :param documents: documents from CosmosDB
    :param message: messages to be published to commit processor topic
    :param context: consists of azure function context information like invocationId and function name
    """
    log.add_keys_from_context(context)
    messages_to_publish: List[str] = []
    try:
        number_of_documents = 0
        for doc in documents:
            number_of_documents += 1
            # process document
            commit_processor_message = package_aggregator_service.process_doc(doc)
            if commit_processor_message is not None:
                messages_to_publish.append(commit_processor_message)
        log.info(f"Number of Documents Processed = {number_of_documents}")
        message.set(json.dumps(messages_to_publish))
    except Exception as e:
        log.exception(f"Exception while processing record in Package Aggregator {e}")
        raise e
