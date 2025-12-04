# updateRTDB

## Description

Runs out of the Google Cloud Scheduler every 4 minutes and retrieves the latest ADSB updates to the firebase
firestore, then writes them to our google realtime DB in adsb4/. 

### Schedule
Crontab/google cloud scheduler
```shell
*/4 * * * *
```

## Deployment

### Google cloud
```shell
gcloud functions deploy updateRTDB --runtime nodejs14 --trigge
r-http --allow-unauthenticated
```

### Google cloud function URL
https://us-central1-chaseapp-8459b.cloudfunctions.net/updateRTDB
