import json
import requests
import boto3
import base64
from botocore.exceptions import ClientError
region_name = "us-west-2"
session = boto3.session.Session()
client = session.client(service_name='secretsmanager',region_name=region_name)
#########################################################prod auth ########################################################
secret_name = "prodkeypostman"
get_secret_value_response = client.get_secret_value(SecretId=secret_name)
#print (get_secret_value_response)
options = get_secret_value_response['SecretString']
#########################################################stage auth ########################################################
secret_name1 = "stagekeypostman"
get_secret_value_response1 = client.get_secret_value(SecretId=secret_name1)
options1 = get_secret_value_response1['SecretString']
#########################################################username########################################################
secret_name2 = "tpaas"
get_secret_value_response2 = client.get_secret_value(SecretId=secret_name2)
options2 = get_secret_value_response2['SecretString']
#print (options2)
#########################################################password########################################################
secret_name3 = "pass2"
get_secret_value_response3 = client.get_secret_value(SecretId=secret_name3)
options3 = get_secret_value_response3['SecretString']
#print (options3)
def prodtoken():
    global x;
    urlprod = "https://identity.trimble.com/token?grant_type=password&username="+options2+"&password="+options3+"&scope=openid"
    payload  = {}
    headers2 = {
            'Authorization': 'Basic '+options,
            'Content-Type': 'application/x-www-form-urlencoded'
    }


    response = requests.request("POST", urlprod, headers=headers2, data = payload)

    mani = response.text
    x = mani[17:49]
def stagetoken():
    global x;
    url = "https://identity-stg.trimble.com/token?grant_type=password&username="+options2+"&password="+options3+"&scope=openid"
    payload  = {}
    headers2 = {
            'Authorization': 'Basic '+options1,
            'Content-Type': 'application/x-www-form-urlencoded'
    }


    response = requests.request("POST", url, headers=headers2, data = payload)

    mani = response.text
    x = mani[17:49]    
def prod(output2):
    print (options)
    global username;
    global domain;
    global x;
    urlprod = "https://identity.trimble.com/token?grant_type=password&username="+options2+"&password="+options3+"&scope=openid"
    payload  = {}
    headers2 = {
            'Authorization': 'Basic '+options,
            'Content-Type': 'application/x-www-form-urlencoded'
    }


    response = requests.request("POST", urlprod, headers=headers2, data = payload)

    mani = response.text
    x = mani[17:49]
    key = output2
    username = key.split("@")[0]
    domain = key.split("@")[1]
def stage(output2):
    global username;
    global domain;
    global x;
    url = "https://identity-stg.trimble.com/token?grant_type=password&username="+options2+"&password="+options3+"&scope=openid"
    payload  = {}
    headers2 = {
            'Authorization': 'Basic '+options1,
            'Content-Type': 'application/x-www-form-urlencoded'
    }

    response = requests.request("POST", url, headers=headers2, data = payload)

    mani = response.text
    x = mani[17:49]
    key = output2
    username = key.split("@")[0]
    domain = key.split("@")[1]
def lambda_handler(event, context):
    global username;
    global domain;
    global x;
    global response1;
    global thread_id;
    input = (event["message"]["argumentText"])
    output = input.split()[0]    ############# postman or delete or help
    if output == "help":
        resText1 = {"text":  "*quickbot operation for view and activate the user in Prod and stage* \n *Format* - @quickbot postman <prod/stage> <user_email ID> <view/activate> \n @quickbot postman prod mohanasundaram_padmanaban@trimble.com view \n @quickbot postman prod mohanasundaram_padmanaban@trimble.com activate \n *quickbot operation to get application details* \n *Format* - @quickbot postman <prod/stage> <application_id> <get>  \n @quickbot postman prod 881da4dd-7608-4d40-9ffc-7cdbbdc29c2e get \n *quickbot operation to search with uuid* \n *Format* - @quickbot postman <prod/stage> <domain> <uuid> \n @quickbot postman stage trimble.com 79d6743d-6ce6-4ab7-91a0-d4df3de56bbb"} 
        return resText1
    else:
        output1 = input.split()[1]   ############# stage or prod
        print (output1)
        output2 = input.split()[2]   ############# mail id
        print (output2)
        output3 = input.split()[3]   ############# view activate token
        print (output3)
        if output == "postman":
            if output1 == "prod" and output3 == "view":
                prod(output2)
                getprod = "https://api.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users/"+username
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("GET", getprod, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1
            elif output1 == "stage" and output3 == "view":
                stage(output2)
                getstage = "https://api-stg.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users/"+username

                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }

                response1 = requests.request("GET", getstage, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1
            elif output1 == "prod" and output3 == "activate":
                prod(output2)
                prodact ="https://api.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users/"+username+"/status/true?notify=false"
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("PUT", prodact, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1            
            elif output1 == "stage" and output3 == "activate":
                stage(output2)
                stageact = "https://api-stg.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users/"+username+"/status/true?notify=false"
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("PUT", stageact, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1
            elif output1 == "prod" and output3 == "get":
                prodtoken()
                application = output2
                url = "https://api.trimble.com/t/trimble.com/identity-admin-service/1.0/oauth2application/"+application
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("GET", url, headers=headers1, data = payload)
                contojson = json.loads(response1.text)
                del contojson['consumerSecret'],contojson['createdAt'],contojson['grantTypes']
                str1 = json.dumps(contojson)
                uuidgiven = contojson["userUUID"]
                url = "https://api.trimble.com/t/trimble.com/identity/2.0/accounts/trimble.com/users?uuid="+uuidgiven
                payload = {}
                headers = {
                    'Authorization': 'Bearer '+x,
                    'Content-Type': 'application/json'
                }   
                response1 = requests.request("GET", url, headers=headers, data = payload)
                str2 = response1.text
                str3 = str1 + str2 
                payload = json.dumps(str3);
                bad_chars = ["\\"]
                for i in bad_chars :
                    test_string = payload.replace(i, '')
                resText1 = {"text":  test_string}
                return resText1
            elif output1 == "stage" and output3 == "get":
                stagetoken()
                application = output2
                url = "https://api-stg.trimble.com/t/trimble.com/identity-admin-service/1.0/oauth2application/"+application
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("GET", url, headers=headers1, data = payload)
                contojson = json.loads(response1.text)
                del contojson['consumerSecret'],contojson['createdAt'],contojson['grantTypes']
                str1 = json.dumps(contojson)
                uuidgiven = contojson["userUUID"]
                url = "https://api-stg.trimble.com/t/trimble.com/identity/2.0/accounts/trimble.com/users?uuid="+uuidgiven
                payload = {}
                headers = {
                    'Authorization': 'Bearer '+x,
                    'Content-Type': 'application/json'
                }   
                response1 = requests.request("GET", url, headers=headers, data = payload)
                str2 = response1.text
                str3 = str1 + str2 
                payload = json.dumps(str3);
                bad_chars = ["\\"]
                for i in bad_chars :
                    test_string = payload.replace(i, '')
                resText1 = {"text":  test_string}
                return resText1
                
            elif output1 == "prod":
                prodtoken()
                domain = output2
                uuid = output3
                uuidprod = "https://api.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users?uuid="+uuid
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("GET", uuidprod, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1
            elif output1 == "stage":
                stagetoken()
                domain = output2
                uuid = output3
                uuidstage = "https://api-stg.trimble.com/t/trimble.com/identity/2.0/accounts/"+domain+"/users?uuid="+uuid
                payload = {}
                headers1 = {
                    'Authorization': "Bearer "+x,
                    'Content-Type': 'application/json'
                }
                response1 = requests.request("GET", uuidstage, headers=headers1, data = payload)
                payload = json.dumps(response1.json());
                resText1 = {"text":  payload}
                return resText1
    