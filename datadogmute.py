from datadog import initialize,  api
import json
import boto3
import base64
from botocore.exceptions import ClientError
import requests
import request
secret_name = "mydemosecret"
region_name = "us-west-2"
# Create a Secrets Manager client
session = boto3.session.Session()
client = session.client(service_name='secretsmanager',region_name=region_name)
try:
        def stgmuted(event):
                get_secret_value_response = client.get_secret_value(SecretId=secret_name)
                options = get_secret_value_response['SecretString']
                res = json.loads(options)
                initialize(**res)
                all_monitor = api.Monitor.get_all()
                for monitor in all_monitor:
                        if monitor['tags']:
                                if monitor['tags'][0] == 'env:mani':
                                        continue
        def prdmuted():
                get_secret_value_response = client.get_secret_value(SecretId=secret_name)
                options = get_secret_value_response['SecretString']
                res = json.loads(options)
                initialize(**res)
                all_monitor = api.Monitor.get_all()
                for monitor in all_monitor:
                        if monitor['tags']:
                                 if monitor['tags'][0] == 'team:ravi':
                                       continue
except:
        print ("test")

        # TODO: write code...
def lambda_handler(event,context):
        #hangout_input = event["message"]["argumentText"]
        thread_id = event["requestAttributes"]["thread_id"]
        #print (hangout_input)
        #url=(event["configCompleteRedirectUrl"])
        #requests.get(url)
        print (event)
        if(event["currentIntent"]["slots"]["slotTwo"]) == "mute" and (event["currentIntent"]["slots"]["slotThree"]) == "stage" and (event["currentIntent"]["slots"]["slotFour"]) == "sms":
                stgmuted(event)
                api.Monitor.mute_all()
                resText1 = {"text": "Alert muted","thread":{ "name" : thread_id}}
                resText1_load1 = json.dumps(resText1)
    
                payload = resText1_load1
                headers = {
                        'Cache-Control': "no-cache",
                        'Content-Type': "application/json"
                        }
    
                response = requests.request("POST", url, data=payload, headers=headers, params=querystring)
                print (response.text)
        elif(event["currentIntent"]["slots"]["slotTwo"]) == "unmute" and (event["currentIntent"]["slots"]["slotThree"]) == "stage" and (event["currentIntent"]["slots"]["slotFour"]) == "sms":
                stgmuted()
                api.Monitor.unmute_all()
                resText1 = {"text": "Alert unmuted","thread":{ "name" : thread_id}}
                resText1_load1 = json.dumps(resText1)
    
                payload = resText1_load1
                headers = {
                        'Cache-Control': "no-cache",
                        'Content-Type': "application/json"
                        }
    
                response = requests.request("POST", url, data=payload, headers=headers, params=querystring)
                print (response.text)
        if(event["currentIntent"]["slots"]["slotTwo"]) == "mute" and (event["currentIntent"]["slots"]["slotThree"]) == "prod" and (event["currentIntent"]["slots"]["slotFour"]) == "sms":
                prdmuted()
                api.Monitor.mute_all()
                resText1 = {"text": "Alert muted","thread":{ "name" : thread_id}}
                resText1_load1 = json.dumps(resText1)
    
                payload = resText1_load1
                headers = {
                        'Cache-Control': "no-cache",
                        'Content-Type': "application/json"
                        }
    
                response = requests.request("POST", url, data=payload, headers=headers, params=querystring)
                print (response.text)
                
        elif(event["currentIntent"]["slots"]["slotTwo"]) == "unmute" and (event["currentIntent"]["slots"]["slotThree"]) == "prod" and (event["currentIntent"]["slots"]["slotFour"]) == "sms":
                prdmuted()
                api.Monitor.unmute_all()
                resText1 = {"text": "Alert unmuted","thread":{ "name" : thread_id}}
                resText1_load1 = json.dumps(resText1)
    
                payload = resText1_load1
                headers = {
                        'Cache-Control': "no-cache",
                        'Content-Type': "application/json"
                        }
    
                response = requests.request("POST", url, data=payload, headers=headers, params=querystring)
                print (response.text)
#print "enter the value in double quotes";
#varaiable = "mute"
#sms(event,context







