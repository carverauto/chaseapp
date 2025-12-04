# BoF (Birds of a Feather)

findBofs is scheduled to run out of google cloud functions every 4 minutes i believe
findBofs will tell you if an LEO bird is close to media
and write to google realtime database in bofs/

## Deployment
```shell
gcloud functions deploy findBofs --runtime nodejs14 --trigger-http --allow-unauthenticated
```
