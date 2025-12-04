# Certificates

https://sid04naik.medium.com/safari-push-notifications-spn-the-complete-setup-guide-6aa49889e8a1 

## The Apple WWDRCA certificate

Apple Worldwide Developer Relations Intermediate Certificate

This is also needed later on when sending notifications. Download it from https://developer.apple.com/support/certificates/expiration/

Convert the “AppleWWDRCA.cer” to “AppleWWDRCA.pem” by using the about steps of converting the .cer file to .pem file.

```
openssl x509 -inform der -in AppleWWDRCA.cer -out cert.pem
```


Not sure if this command below will work, it was taken out of the documentation from the first link. It did
not work and I had to use the above command instead.

```
openssl pkcs12 -in certificate.p12 -out certificate.pem -nodes
```
