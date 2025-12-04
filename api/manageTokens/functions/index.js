const apn = require('apn');
const functions = require("firebase-functions");
const admin = require('firebase-admin')
// const apnKey = require('./AuthKey_6DP4R2NA4X.p8')
admin.initializeApp()
// FCM messaging topic to subscribe to
const topic = 'chases'
let apnKey = ''

exports.addFCMToken = functions.firestore
  .document('tokens/{docId}').onCreate(handler => {
        const token = handler.token
        console.log(`Token: ${token}`)

        if (token) {
            admin.messaging().subscribeToTopic(token, topic)
            // const response = admin.messaging().subscribeToTopic([ token ], topic)
            admin.messaging().subscribeToTopic([ token ], topic).then((sub) => {
                console.log(sub)
                functions.logger.info('Subscribed to topic: ', topic)
            }).catch((e) => {
                functions.logger.error('Error subscribing to topic: ', e)
            })
            return 1
        }
    })

exports.notifyAPN = functions.firestore
    .document('chases/{docId}')
    .onUpdate(async (change, context) => {
        const newValue = change.after.data()
        const newFieldValue = newValue.Live

        const previousValue = change.before.data()
        const previousFieldValue = previousValue.Live

        // we want to only run when the 'Live' field changes to True
        if (previousFieldValue !== newFieldValue)
            // is the field set to True?
            if (newFieldValue) {
                console.log('Live is set to true')
                // get the APN key..
                console.log('Found the APN key..')
                // Set up apn with the APNs Auth Key
                const apnProvider = new apn.Provider({
                    token: {
                        key: "./AuthKey_949L536D8R.p8", // Path to the key p8 file
                        keyId: '4DU9DWL5WS', // The Key ID of the p8 file (available at https://developer.apple.com/account/ios/certificate/key)
                        teamId: '432Q4W72Q7', // The Team ID of your Apple Developer Account (available at https://developer.apple.com/account/#/membership/)
                    },
                    production: true, // Set to true if sending a notification to a production iOS app
                });

                // build a list of all the APN device tokens
                const db = admin.firestore()
                const tokensRef = db.collection('apntokens')
                const snapshot = await tokensRef.get()
                if (snapshot.empty) {
                    console.log('No matching documents.')
                    return
                }

                let tokensList = []

                snapshot.forEach(doc => {
                    console.log(doc.id, '=>', doc.data())
                    const tokenDoc = doc.data()
                    tokensList.push(tokenDoc.token)
                })
                console.log(tokensList)
                console.log('Sending notification')
                //const deviceToken = '5311839E985FA01B56E7AD74444C03133745D0xDEADBEEF';
                const notification = new apn.Notification();
                notification.topic = 'web.tv.chaseapp';
                notification.expiry = Math.floor(Date.now() / 1000) + 3600;
                notification.badge = 3;
                notification.sound = 'ping.aiff';
                notification.alert = 'ChaseApp - WE HAVE A CHASE!';
                notification.payload = {id: 123};
                apnProvider.send(notification, tokensList).then(function (result) {
                    // Check the result for any failed devices
                    console.log(result);
                });
            }
    })