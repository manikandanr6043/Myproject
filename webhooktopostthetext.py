import json
from json import dumps

from httplib2 import Http

def main(event,context):
	###########################################Webhook URL Below#############################################
    url = 'https://chat.googleapis.com/v1/spaces/AAAAZIWeu0E/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=vcicDPAZNi89EDRy2yVdR6sHBxO0-6xkG0rACQK7K1k%3D'
    ###########################################Bot Message###################################################
    bot_message = {
        'text' : 'Guys in shift please join the meeting https://meet.google.com/pxk-jgzk-sdg'}

    message_headers = {'Content-Type': 'application/json; charset=UTF-8'}

    http_obj = Http()
	###########################################Post to the Chat webhook#######################################
    response = http_obj.request(
        uri=url,
        method='POST',
        headers=message_headers,
        body=json.dumps(bot_message),
    )




