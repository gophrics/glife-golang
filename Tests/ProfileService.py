import requests
import Utils
import Context
import json


# Public API
def username_exist():
    username = Utils.generate_username()
    print("http://localhost:8080/api/v1/profile/username_exist")
    url = 'http://localhost:8080/api/v1/profile/username_exist/' + username
    r = requests.Session()
    response = r.get(url)
    print(str(response.status_code))
    print(str(response.content))

def register():
    print("http://localhost:8080/api/v1/profile/register")
    url = 'http://localhost:8080/api/v1/profile/register'
    name = Utils.generate_name()
    password = Utils.generate_password()
    email = Utils.generate_email()
    country = Utils.generate_country()
    phone = Utils.generate_phone()
    username = Utils.generate_username()
    data = dict (
        name=name,
        password=password,
        email=email,
        phone=phone,
        country=country,
        username=username
    )

    print(data)
    r = requests.Session()
    r.headers.update({"Content-Type": "application/json"})
    response = r.post(url, data=json.dumps(data))
    print("Response status: " + str(response.status_code))
    if response.status_code == 200:
        print("Token: " + response.json()['Token'])
        Context.token = response.json()['Token']
        Context.username = username
        return True
    else:
        print("Response: " + str(response.content))
    return False

def get_me():
    print("http://localhost:8080/api/v1/profile/getme")
    url = 'http://localhost:8080/api/v1/profile/getme'
    r = requests.Session()
    r.headers.update({"Authorization": "Bearer " + Context.token})
    response = r.get(url)
    print("Response status: " + str(response.status_code))
    if response.status_code == 200:
        print("Response: " + str(response.json()))
    else:
        print("Response: " + str(response.content))

def get_user():
    print("http://localhost:8080/api/v1/profile/getuserwithusername/oynavzqhs")
    url = "http://localhost:8080/api/v1/profile/getuserwithusername/oynavzqhs"
    r = requests.Session()
    r.headers.update({"Authorization": "Bearer " + Context.token})
    response = r.get(url)
    print("Response status: " + str(response.status_code))
    if response.status_code == 200:
        print("Response: " + str(response.json()))
    else:
        print("Response: " + str(response.content))