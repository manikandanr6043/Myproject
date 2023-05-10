"""
This module consists of the main function.
The main method is the method in your function code that processes events.
"""
import json

import azure.functions as func

from thumbnail_observer.thumbnail_observer_service import ThumbnailObserverService
from common.logging.custom_logger import CustomLogger

log = CustomLogger()
thumbnail_observer_service = ThumbnailObserverService()


def main(event: func.EventGridEvent, message: func.Out[str], context: func.Context):
    """
       :param event: Thumbnail Blob created event
       :param message: message to be sent to the thumb processor topic
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
        thumb_processor_msg = thumbnail_observer_service.generate_thumb_processor_msg(event.id, event_subject)
        # Post message to thumb processor kafka topic
        message.set(thumb_processor_msg)
    except Exception as e:
        log.exception(f"Thumbnail Observer failed! {e}")
        raise e
