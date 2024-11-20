import requests

#Local endpoint
url = "http://localhost:8080/getIndex"

#Get the response
response = requests.get(url)

if response.status_code == 200:
    #Print the reponse
    securities = response.json()
    print(securities["securities"])

else:
    #Request failed
    print(f"Failed API Request, status code: {response.status_code}")