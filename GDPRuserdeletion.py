import requests
file1 = open("testing.txt","r")    #########need to put the mailid in that texting.txt file########################
f1 = file1.readlines()
key = str(input("Key :"))          ##############pass the access token key########################
for email in f1:
    print (email)
    url = "https://api.trimble.com/t/trimble.com/identity/3.0/users/purge-request?email="+email+"&notify=false"
    payload = {}
    headers = {
            'Authorization': "Bearer "+key,
            'Content-Type': 'application/json'
    }

    response = requests.request("DELETE", url, headers=headers, data = payload)   ###################user will deleted with help of delete method###########

    print(response.text.encode('utf8'))