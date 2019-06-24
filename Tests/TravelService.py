import requests
import Context
import Utils
import json

def save_trip():
    print("http://localhost:8082/api/v1/travel/savetrip")
    url = 'http://localhost:8082/api/v1/travel/savetrip'
    trips = Utils.generate_trip()
    r = requests.Session()
    r.headers.update({"Authorization": "Bearer " + Context.token})
    response = r.post(url, data=json.dumps(trips))
    print(str(response.status_code))
    if response.status_code == 200:
        print(str(response.content))

def get_all_trips(username):
    print("http://localhost:8082/api/v1/travel/getalltrips")
    print(username)
    url = 'http://localhost:8082/api/v1/travel/getalltrips/' + username
    r = requests.Session()
    r.headers.update({"Authorization": "Bearer " + Context.token})
    response = r.get(url)
    print(str(response.status_code))
    if response.status_code == 200:
        print(str(response.content))