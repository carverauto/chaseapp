# Deploying google cloud functions:

## golang runtime 1.13

```bash
gcloud functions deploy ListChases --runtime go116 --trigger-http --allow-unauthenticated --env-vars-file .env.yaml
gcloud functions deploy AddChase --runtime go116 --trigger-http --allow-unauthenticated --env-vars-file .env.yaml
gcloud functions deploy UpdateChase --runtime go116 --trigger-http --allow-unauthenticated --env-vars-file .env.yaml
gcloud functions deploy DeleteChase --runtime go116 --trigger-http --allow-unauthenticated --env-vars-file .env.yaml

gcloud functions deploy ListAirships --runtime go113 --trigger-http --set-env-vars APIKEY=<apikey>
gcloud functions deploy AddAirship --runtime go113 --trigger-http --set-env-vars APIKEY=<apikey>
gcloud functions deploy UpdateAirship --runtime go113 --trigger-http --set-env-vars APIKEY=<apikey>
gcloud functions deploy DeleteAirship --runtime go113 --trigger-http --set-env-vars APIKEY=<apikey>
gcloud functions deploy GetBoats --runtime go113 --trigger-http --allow-unauthenticated --set-env-vars APIKEY=<apikey> AISHUB_USERNAME=<AISHUB_USERNAME>
gcloud functions deploy GetLaunches --runtime go113 --trigger-http --allow-unauthenticated --set-env-vars APIKEY=<apikey> ROCKETLAUNCHAPI=<ROCKETLAUNCHAPI>

```

## Endpoints

* /AddUser ("/users")
* /UpdateToken (FCM)
* /DeleteToken (FCM)
* /GetChase
* /UpdateChase
* /AddChase
* /ListChases
* /DeleteChase
* /ListAirships
* /AddAirship
* /UpdateAirship
* /DeleteAirship
* /GetWeatherAlerts
* /GetBoats
* /GetLaunches

# Testing with CURL:

## Add Chase -

```sh
curl --header "X-ApiKey: <ApiKey>" --request POST -d '{"name":"Pursuit in Azusa area","desc":"Police are in Pursuit of Vehicle in Azusa Area","url":"https://www.facebook.com/CBSLA/videos/642113759585815/","live":true}' -H 'Content-Type: application/json' https://us-central1-chaseapp-8459b.cloudfunctions.net/AddChase
```

## List Chases -

```sh
curl -H 'Content-Type: application/json' https://us-central1-chaseapp-8459b.cloudfunctions.net/ListChases
```
