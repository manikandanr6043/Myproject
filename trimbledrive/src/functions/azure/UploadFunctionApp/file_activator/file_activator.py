"""
This module consists of the main function.
The main method is the method in your function code that processes events.
The Activator function received file creation events from our content bucket.
"""
import json
from file_activator.file_activator_service import FileActivatorService
from common.logging.custom_logger import CustomLogger
import azure.functions as func

log = CustomLogger()
file_activator = FileActivatorService()


def main(event: func.EventGridEvent, context: func.Context):
    """
       :param event: Blob created event
       :param context: consists of azure function context information like invocationId and function name
    """
    log.add_keys_from_context(context)
    event_subject = event.subject
    event_dict = {
        'id': event.id,
        'data': json.dumps(event.get_json()),
        'topic': event.topic,
        'subject': event_subject,
        'event_type': event.event_type,
    }
    log.info(f'Processing event: {event_subject}', extra=event_dict)
    try:
        file_activator.activate_file(event_subject)
    except Exception as e:
        log.exception(f"File Activator failed! {e}")
        raise e
