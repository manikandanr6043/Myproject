import boto3
import json
import requests
import time

####################
ec2 = boto3.client('ec2')
global user_instance;
global user_command;
global output;
client = boto3.client('ssm')
def sendcommand():
    global user_instance;
    global user_command;
    global output;
    global restext1;
    response = client.send_command(
        InstanceIds=[
            user_instance,
            ],
        DocumentName='AWS-RunPowerShellScript',
        Comment='IP',
        Parameters={
            'commands': [
                user_command,
            ]
        },
    )
    command_id = response['Command']['CommandId']
    print(response)
    print(command_id)
    time.sleep(6)
    output = client.get_command_invocation(
        CommandId = command_id,
        InstanceId = user_instance,
    )
    print(output)
def lambda_handler(event, context):
    input = (event["message"]["argumentText"])
    #output = "44.232.209.100"
    output3 = input.split()[0]
    output1 = input.split()[1]
    if output3 == "instance":
        response = ec2.describe_regions()
        for i in response['Regions']:
            region = i['RegionName']
            session_res=boto3.client(service_name="ec2",region_name=region)
            f1={"Name":"network-interface.addresses.private-ip-address","Values":[output1]}
            f2={"Name":"network-interface.addresses.association.public-ip","Values":[output1]}
            f3={"Name":"instance-id","Values":[output1]}
            my_instance=session_res.describe_instances(Filters=[f1]) ######Filter for Private IP #############
            my_instance1=session_res.describe_instances(Filters=[f2]) #####Filter for Public IP  #############
            my_instance2=session_res.describe_instances(Filters=[f3]) #####Filter for Instance ID#############
            for Tags in my_instance['Reservations']:
                for tag in Tags['Instances']:
                    for taged in tag['Tags']:
                        if taged['Key'] == 'Name':
                            response1 = (tag['InstanceId'],tag['PrivateIpAddress'],tag['State']['Name'],taged['Value'],region,tag['Platform'])
                            payload = json.dumps(response1);
                            resText1 = {"text":  payload}
                            return resText1
            for Tags in my_instance1['Reservations']:
                for tag in Tags['Instances']:
                    for taged in tag['Tags']:
                        if taged['Key'] == 'Name':
                            response1 = (tag['InstanceId'],tag['PrivateIpAddress'],tag['State']['Name'],taged['Value'],region,tag['Platform'])
                            payload = json.dumps(response1);
                            resText1 = {"text":  payload}
                            return resText1
            for Tags in my_instance2['Reservations']:
                for tag in Tags['Instances']:
                    for taged in tag['Tags']:
                        if taged['Key'] == 'Name':
                            response1 = (tag['PrivateIpAddress'],tag['State']['Name'],taged['Value'],region,tag['Platform'])
                            payload = json.dumps(response1);
                            resText1 = {"text":  payload}
                            return resText1
    else:
        global user_instance;
        global user_command;
        global output;
        print(event)
        mail = (event["message"]["sender"]["email"])
        print(mail)
        print(type(mail))
        mailid = ["premkumar_jawahar@trimble.com","sathishkumar_m@trimble.com","manikandan_ravichandran@trimble.com",
        "prasaanth_sriinivasan@trimble.com"] #############List of person have a permission to access############
        adminmail =["manikandan_ravichandran@trimbl.com"] #####Admin User access############################
# Using loop for constructing output list
        if mail in mailid :
            event_Message = (event["message"]["argumentText"])
            msg = event_Message.split()
            user_instance = msg[0]
            msg.pop(0)
            str_space=" "
            user_command=str_space.join(msg)
            print(user_instance)
            print(user_command)
            commandwin = ["Get-PSDrive","Get-Service AmazonSSMAgent","get-service winraa",
            "get-service winrm","start-service winrm","stop-service winrm","aws --version"] #########List of command to execute##############
            if user_command in commandwin:
                sendcommand()
                if output["StandardOutputContent"] !="":
                    mani = output["StandardOutputContent"]
                    restext1 = {"text" : mani}
                    return restext1
                elif output["StandardErrorContent"] !="":
                    manik = output["StandardErrorContent"] 
                    restext2 = {"text" : manik}
                    print(type(restext2))
                    return restext2
                elif output["StatusDetails"]:
                    manika = output["StatusDetails"] 
                    restext2 = {"text" : manika}
                    print(type(restext2))
                    return restext2
                
            else:
                restext1={"text": user_command + " is not authorized to execute command"}
                return restext1
        elif mail in adminmail:
            event_Message = (event["message"]["argumentText"])
            msg = event_Message.split()
            user_instance = msg[0]
            msg.pop(0)
            str_space=" "
            user_command=str_space.join(msg)
            print(user_instance)
            print(user_command)
            sendcommand()
            if output["StandardOutputContent"] !="":
                mani = output["StandardOutputContent"]
                restext1 = {"text" : mani}
                return restext1
            elif output["StandardErrorContent"] !="":
                manik = output["StandardErrorContent"] 
                restext1 = {"text" : manik}
                print(type(restext1))
                return restext1
            elif output["StatusDetails"]:
                manika = output["StatusDetails"] 
                restext1 = {"text" : manika}
                print(type(restext1))
                return restext1  
        elif mail:
            restext1 = {"text" :"Access Denied"}
            return restext1