import requests

#Local endpoint
url = "http://localhost:8080/getSec"

#Get the response
response = requests.get(url)

if response.status_code == 200:
    #Print the reponse
    print(response.json())
    

else:
    #Request failed
    print(f"Failed API Request, status code: {response.status_code}")