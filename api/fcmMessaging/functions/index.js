const functions = require("firebase-functions");
const admin = require('firebase-admin')
admin.initializeApp()
const db = admin.firestore()
const hash = require('crypto').createHash
const topic = 'chases'

// // Create and Deploy Your First Cloud Functions
// // https://firebase.google.com/docs/functions/write-firebase-functions
//
// exports.helloWorld = functions.https.onRequest((request, response) => {
//   functions.logger.info("Hello logs!", {structuredData: true});
//   response.send("Hello from Firebase!");
// });
//const notifyUsers = require('./notify-users');

exports.fcmToken = functions.firestore.document('/users/{documentId}/messaging/{token}')
  .onUpdate(async (change, context) => {

    admin.messaging().subscribeToTopic(registrationTokens, topic)
        .then(function(response) {
          // See the MessagingTopicManagementResponse reference documentation
          // for the contents of response.
          console.log('Successfully subscribed to topic:', response);
        })
        .catch(function(error) {
          console.log('Error subscribing to topic:', error);
        });
  });
