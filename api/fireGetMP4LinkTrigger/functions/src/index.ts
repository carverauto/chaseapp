import * as functions from "firebase-functions";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const axios = require("axios");

// // Start writing Firebase Functions
// // https://firebase.google.com/docs/functions/typescript
//
// export const helloWorld = functions.https.onRequest((request, response) => {
//   functions.logger.info("Hello logs!", {structuredData: true});
//   response.send("Hello from Firebase!");
// });

const GetMP4LinkURL = "https://us-central1-chaseapp-8459b.cloudfunctions.net/GetMP4Link";

exports.FireGetMP4LinkTrigger = functions.firestore
    .document("chases/{docId}")
    .onWrite((change, context) => {
      const data = change.after.data();
      const previousData = change.before.data();
      if (data !== null) {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        data.sentiment = null;
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        previousData.sentiment = null;

        if (data == previousData) {
          return null;
        }
        const doc = change.after.exists ? change.after.data() : null;
        if (doc) {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          const requestBody = {
            chase_id: doc.ID,
          };
          return axios.post(GetMP4LinkURL, requestBody)
              .then((res: any) => {
                console.log(res.status);
              }).catch((err: any) => {
                return console.log(err);
              });
        } else {
          console.log("No data");
        }
      } else {
        console.log("No data, possibly deleted");
      }
    });
