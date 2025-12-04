# Safari push package stuff

https://stackoverflow.com/questions/64289662/create-safari-push-notification-signature-with-node-js
https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/NotificationProgrammingGuideForWebsites/PushNotifications/PushNotifications.html#//apple_ref/doc/uid/TP40013225-CH3-SW7

# Deploying

```shell
gcloud functions deploy pushPackage --runtime nodejs14 --trigger-http --allow-unauthenticated
```


# Certs

## Download the Apple WWDRCA intermediate cert and convert .cer to .pem

```shell
openssl x509 -inform der -in AppleWWDRCA.cer -out cert.pem
```

