const functions = require("firebase-functions");
const admin = require("firebase-admin");
admin.initializeApp();
const db = admin.firestore();
const topic = "chases";
const apn = require("@parse/node-apn");

// APN options
/*
const options = {
  token: {
    key: "./AuthKey_6DP4R2NA4X.p8",
    keyId: "4DU9DWL5WS",
    teamId: "432Q4W72Q7",
  },
  production: false,
};
*/
// const apnProvider = new apn.Provider(options);

// google cloud pubsub stuff
// const {PubSub} = require("@google-cloud/pubsub");
const {v1} = require("@google-cloud/pubsub");
const projectId = "chaseapp-8459b";
const topicName = "chases";

const pubsub = new v1.PublisherClient({
  projectId: projectId,
  keyFilename: "./account.json",
});
/*
const pubsub = new PubSub({
  projectId: psubConfig.projectId,
  keyFilename: "./account.json",
});
 */

const formattedTopic = pubsub.projectTopicPath(
    projectId,
    topicName,
);

const retrySettings = {
  retryCodes: [
    10, // 'ABORTED'
    1, // 'CANCELLED',
    4, // 'DEADLINE_EXCEEDED'
    13, // 'INTERNAL'
    8, // 'RESOURCE_EXHAUSTED'
    14, // 'UNAVAILABLE'
    2, // 'UNKNOWN'
  ],
  backoffSettings: {
    initialRetryDelayMillis: 100,
    retryDelayMultiplier: 1.3,
    maxRetryDelayMillis: 60000,
    initialRpcTimeoutMillis: 5000,
    rpcTimeoutMultiplier: 1.0,
    maxRpcTimeoutMillis: 600000,
    totalTimeoutMillis: 600000,
  },
};

const fcmMessage = {
  data: {
    update: "chase-update",
    createdAt: Date.now().toString(),
  },
  topic: topic,
};

const note = new apn.Notification();

note.expiry = Math.floor(Date.now() / 1000) + 3600; // Expires 1 hour from now.
note.badge = 3;
note.sound = "ping.aiff";
note.alert = "\uD83D\uDCE7 \u2709 ChaseApp - Alert";
note.payload = {"messageFrom": "ChaseApp"};
note.topic = "432Q4W72Q7";

exports.updateUI = functions.firestore
    .document("chases/{docId}")
    .onWrite(async (change, context) => {
      // Update pubsub
      const myData = change.after.data();
      myData.ID = change.after.ref.id;
      const dataBuffer = Buffer.from(JSON.stringify(myData));
      const messagesElement = {
        data: dataBuffer,
      };
      const messages = [messagesElement];
      // Build the request
      const request = {
        topic: formattedTopic,
        messages: messages,
      };

      const [response] = await pubsub.publish(request, {
        retry: retrySettings,
      });
      console.log(`Message ${response.messageIds} published.`);
      /*
        pubsub.topic("chases")
            .publish(Buffer.from(JSON.stringify(myData)), {retrySettings})
            .then((res) => {
              console.log(res);
            }).catch((error) => {
              console.error("Error publishing to topic chases: " + error);
            });
         */
      // FCM Messaging to everyone but safari..
      admin.messaging().send(fcmMessage).then((res) => {
        console.log("Successfully sent message ", res);
      }).catch((e) => {
        console.error("Error sending message", e);
      });
    /*
      apnProvider.send(note, [""]).then((r) => {
        if (r.failed) {
          // remove that token
        }
      });
         */
    });

exports.createBundle = functions.https.onRequest(async (request, response) => {
  // Query the 20 latest chases
  const latestChases = await db.collection("chases")
      .orderBy("CreatedAt", "desc")
      .limit(20)
      .get();

  // Build the bundle from the query results
  const bundleBuffer = db.bundle("latest-chases")
      .add("latest-chases-query", latestChases)
      .build();

  // Cache the response for up to 5 minutes;
  // see https://firebase.google.com/docs/hosting/manage-cache
  response.set("Cache-Control", "public, max-age=300, s-maxage=600");

  response.end(bundleBuffer);
});
