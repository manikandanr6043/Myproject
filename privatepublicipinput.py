import json
import boto3
ec2 = boto3.client('ec2')
def lambda_handler(event, context):
    input = (event["message"]["argumentText"])
    #output = "44.232.209.100"
    output = input.split()[0]
    if output:
        response = ec2.describe_regions()
        for i in response['Regions']:
            region = i['RegionName']
            session_res=boto3.resource(service_name="ec2",region_name=region)
            f1={"Name":"network-interface.addresses.private-ip-address","Values":[output]}
            f2={"Name":"network-interface.addresses.association.public-ip","Values":[output]}
            my_instance=session_res.instances.filter(Filters=[f1])
            my_instance1=session_res.instances.filter(Filters=[f2])
            for each_in in my_instance:
                response1 = (each_in.id,each_in.state['Name'],region)
                payload = json.dumps(response1);
                resText1 = {"text":  payload}
                return resText1
                if each_in.id == i
            for each_in2 in my_instance1:
                response1 = (each_in2.id,each_in2.state['Name'],region)
                payload = json.dumps(response1);
                resText1 = {"text":  payload}
                return resText1
                
        
        
    