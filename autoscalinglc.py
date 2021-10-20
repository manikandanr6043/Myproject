########### importing Modules ###################
import boto3
import base64
import datetime
clientec2 = boto3.client('ec2')
client = boto3.client('autoscaling')
e = datetime.datetime.now()
time = e.strftime("%d_%m_%Y")+e.strftime("%I_%M_%S")
autogroup = str(input("autoscaling group name:"))
lanuchname = autogroup+time
#################Describe existing autoscaling group#######################################
response = client.describe_auto_scaling_groups(
    AutoScalingGroupNames=[
        autogroup,
    ],
)
for test in response["AutoScalingGroups"]:
    for instance_id in test["Instances"]:
        #print(instance_id["InstanceId"])
        instance_id= instance_id["InstanceId"]
launch=test["LaunchConfigurationName"]
######################Describe the existing LC#######################################
describelaunch = client.describe_launch_configurations(LaunchConfigurationNames=[launch])
for lauchconfiguration in describelaunch["LaunchConfigurations"]:
     imageid = lauchconfiguration["ImageId"]
     keyname = lauchconfiguration["KeyName"]
     securityGroups = lauchconfiguration["SecurityGroups"]
     userdatas = lauchconfiguration["UserData"]
     base64_bytes = userdatas.encode('ascii')
     message_bytes = base64.b64decode(base64_bytes)
     userdata = message_bytes.decode('ascii')
     instancetype = lauchconfiguration["InstanceType"]
     InstanceMonitor=lauchconfiguration["InstanceMonitoring"]
     try:
         Spotprice = lauchconfiguration["SpotPrice"]
         spot=Spotprice
     except KeyError:
         spot=None
     try:
         Iam=lauchconfiguration["IamInstanceProfile"]
     except KeyError:
         print("No Iam Profile")
     if (lauchconfiguration["EbsOptimized"]) == True:
         EBSopt=lauchconfiguration["EbsOptimized"]
     elif (lauchconfiguration["EbsOptimized"]) == False:
         EBSopt = lauchconfiguration["EbsOptimized"]
     response9 = clientec2.describe_images(ImageIds=[imageid])
     for testing in response9["Images"]:
         amiid=testing["ImageId"]
         print(testing["BlockDeviceMappings"])
     BlockDeviceMappings=testing["BlockDeviceMappings"]
     count1 = len(BlockDeviceMappings)
     count = count1 - 1
#################Assigning the BlockDeviceMappings to LC########################
     if count==0:
            device = BlockDeviceMappings[0]
            try:
                if device["Ebs"]["Iops"] <= 3000:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
     elif count==1:
            device = BlockDeviceMappings[0]
            try:
                if device["Ebs"]["Iops"] <= 3000:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
            device1 = BlockDeviceMappings[1]
            try:
                if device1["Ebs"]["Iops"] <= 3000:
                    device1["Ebs"]["Iops"] = 3000
                    device1["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device1["Ebs"]["Iops"] = 3000
                    device1["Ebs"]["VolumeType"] = "gp3"
     elif count==2:
            device = BlockDeviceMappings[0]
            try:
                if device["Ebs"]["Iops"] <= 3000:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device["Ebs"]["Iops"] = 3000
                    device["Ebs"]["VolumeType"] = "gp3"
            device1 = BlockDeviceMappings[1]
            try:
                if device1["Ebs"]["Iops"] <= 3000:
                    device1["Ebs"]["Iops"] = 3000
                    device1["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device1["Ebs"]["Iops"] = 3000
                    device1["Ebs"]["VolumeType"] = "gp3"
            device2 = BlockDeviceMappings[1]
            try:
                if device2["Ebs"]["Iops"] <= 3000:
                    device2["Ebs"]["Iops"] = 3000
                    device2["Ebs"]["VolumeType"] = "gp3"
            except KeyError:
                    device2["Ebs"]["Iops"] = 3000
                    device2["Ebs"]["VolumeType"] = "gp3"
     elif count==3:
         device = BlockDeviceMappings[0]
         try:
             if device["Ebs"]["Iops"] <= 3000:
                 device["Ebs"]["Iops"] = 3000
                 device["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device["Ebs"]["Iops"] = 3000
             device["Ebs"]["VolumeType"] = "gp3"
         device1 = BlockDeviceMappings[1]
         try:
             if device1["Ebs"]["Iops"] <= 3000:
                 device1["Ebs"]["Iops"] = 3000
                 device1["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device1["Ebs"]["Iops"] = 3000
             device1["Ebs"]["VolumeType"] = "gp3"
         device2 = BlockDeviceMappings[2]
         try:
             if device2["Ebs"]["Iops"] <= 3000:
                 device2["Ebs"]["Iops"] = 3000
                 device2["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device2["Ebs"]["Iops"] = 3000
             device2["Ebs"]["VolumeType"] = "gp3"
         device3 = BlockDeviceMappings[3]
         try:
             if device3["Ebs"]["Iops"] <= 3000:
                device3["Ebs"]["Iops"] = 3000
                device3["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device3["Ebs"]["Iops"] = 3000
            device3["Ebs"]["VolumeType"] = "gp3"
     elif count==4:
         device = BlockDeviceMappings[0]
         try:
             if device["Ebs"]["Iops"] <= 3000:
                 device["Ebs"]["Iops"] = 3000
                 device["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device["Ebs"]["Iops"] = 3000
             device["Ebs"]["VolumeType"] = "gp3"
         device1 = BlockDeviceMappings[1]
         try:
             if device1["Ebs"]["Iops"] <= 3000:
                 device1["Ebs"]["Iops"] = 3000
                 device1["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device1["Ebs"]["Iops"] = 3000
             device1["Ebs"]["VolumeType"] = "gp3"
         device2 = BlockDeviceMappings[2]
         try:
             if device2["Ebs"]["Iops"] <= 3000:
                 device2["Ebs"]["Iops"] = 3000
                 device2["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device2["Ebs"]["Iops"] = 3000
             device2["Ebs"]["VolumeType"] = "gp3"
         device3 = BlockDeviceMappings[3]
         try:
             if device3["Ebs"]["Iops"] <= 3000:
                device3["Ebs"]["Iops"] = 3000
                device3["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device3["Ebs"]["Iops"] = 3000
            device3["Ebs"]["VolumeType"] = "gp3"
         device4 = BlockDeviceMappings[4]
         try:
             if device4["Ebs"]["Iops"] <= 3000:
                device4["Ebs"]["Iops"] = 3000
                device4["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device4["Ebs"]["Iops"] = 3000
            device4["Ebs"]["VolumeType"] = "gp3"
     elif count ==5:
         device = BlockDeviceMappings[0]
         try:
             if device["Ebs"]["Iops"] <= 3000:
                 device["Ebs"]["Iops"] = 3000
                 device["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device["Ebs"]["Iops"] = 3000
             device["Ebs"]["VolumeType"] = "gp3"
         device1 = BlockDeviceMappings[1]
         try:
             if device1["Ebs"]["Iops"] <= 3000:
                 device1["Ebs"]["Iops"] = 3000
                 device1["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device1["Ebs"]["Iops"] = 3000
             device1["Ebs"]["VolumeType"] = "gp3"
         device2 = BlockDeviceMappings[2]
         try:
             if device2["Ebs"]["Iops"] <= 3000:
                 device2["Ebs"]["Iops"] = 3000
                 device2["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
             device2["Ebs"]["Iops"] = 3000
             device2["Ebs"]["VolumeType"] = "gp3"
         device3 = BlockDeviceMappings[3]
         try:
             if device3["Ebs"]["Iops"] <= 3000:
                device3["Ebs"]["Iops"] = 3000
                device3["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device3["Ebs"]["Iops"] = 3000
            device3["Ebs"]["VolumeType"] = "gp3"
         device4 = BlockDeviceMappings[4]
         try:
             if device4["Ebs"]["Iops"] <= 3000:
                device4["Ebs"]["Iops"] = 3000
                device4["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device4["Ebs"]["Iops"] = 3000
            device4["Ebs"]["VolumeType"] = "gp3"
         device5 = BlockDeviceMappings[5]
         try:
            if device5["Ebs"]["Iops"] <= 3000:
                device5["Ebs"]["Iops"] = 3000
                device5["Ebs"]["VolumeType"] = "gp3"
         except KeyError:
            device5["Ebs"]["Iops"] = 3000
            device5["Ebs"]["VolumeType"] = "gp3"
#########################Creating new LC##################################################3
     if spot:
         try:
            Iam = lauchconfiguration["IamInstanceProfile"]
            response4 = client.create_launch_configuration(
                LaunchConfigurationName=lanuchname,
                ImageId=amiid,
                KeyName=keyname,
                SecurityGroups=securityGroups,
                InstanceType=instancetype,
                UserData=userdata,
                BlockDeviceMappings=BlockDeviceMappings,
                IamInstanceProfile=Iam,
                InstanceMonitoring=InstanceMonitor,
                EbsOptimized=EBSopt,
                SpotPrice=spot,
                AssociatePublicIpAddress=True
            )
         except KeyError:
            response4 = client.create_launch_configuration(
                LaunchConfigurationName=lanuchname,
                ImageId=amiid,
                KeyName=keyname,
                SecurityGroups=securityGroups,
                InstanceType=instancetype,
                UserData=userdata,
                BlockDeviceMappings=BlockDeviceMappings,
                InstanceMonitoring=InstanceMonitor,
                EbsOptimized=EBSopt,
                SpotPrice=spot,
                AssociatePublicIpAddress=True
                )
     else:
        try:
            Iam = lauchconfiguration["IamInstanceProfile"]
            print(Iam)
            response4 = client.create_launch_configuration(
                LaunchConfigurationName=lanuchname,
                ImageId=amiid,
                KeyName=keyname,
                SecurityGroups=securityGroups,
                InstanceType=instancetype,
                UserData=userdata,
                BlockDeviceMappings=BlockDeviceMappings,
                IamInstanceProfile=Iam,
                InstanceMonitoring=InstanceMonitor,
                EbsOptimized=EBSopt,
                AssociatePublicIpAddress=True
            )
        except KeyError:
            response4 = client.create_launch_configuration(
                LaunchConfigurationName=lanuchname,
                ImageId=amiid,
                KeyName=keyname,
                SecurityGroups=securityGroups,
                InstanceType=instancetype,
                UserData=userdata,
                BlockDeviceMappings=BlockDeviceMappings,
                InstanceMonitoring=InstanceMonitor,
                EbsOptimized=EBSopt,
                AssociatePublicIpAddress=True
                )
updateautoscalinggroup = client.update_auto_scaling_group(
        AutoScalingGroupName=autogroup,
        LaunchConfigurationName=lanuchname)
print(updateautoscalinggroup)



