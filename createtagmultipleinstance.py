import boto3
ec2 = boto3.client('ec2')
Resour=["i-0592759166a48eb1a","i-04f954b830170d82f"]
print(type(Resour))
def lambda_handler(event,context):
    for i in Resour:
          ec2.create_tags(Resources=[i], Tags=[{'Key':'ssm', 'Value':'yes'}])








