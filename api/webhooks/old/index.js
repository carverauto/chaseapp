// The goal of this code is to post some shit to this webhook. (discord)
// https://discord.com/api/webhooks/872966356620955718/mVYTmjIcP7_V_PjscznxQNcL4YJ9rf0Ul6DjmABzUc-XR869Kk0jSSnlNuNTuHRHB-rg

const { Webhook } = require('discord-webhook-node')
const hook = new Webhook("https://discord.com/api/webhooks/872966356620955718/mVYTmjIcP7_V_PjscznxQNcL4YJ9rf0Ul6DjmABzUc-XR869Kk0jSSnlNuNTuHRHB-rg")
const IMAGE_URL = 'https://homepages.cae.wisc.edu/~ece533/images/airplane.png'

const firebase = require('firebase')

firebase.initializeApp(( {
    apiKey: "AIzaSyDZVvCuh81AYFsNqNhdI5GUzwQC91na580",
    authDomain: "chaseapp-8459b.firebaseapp.com",
    databaseURL: "https://chaseapp-8459b.firebaseio.com",
    projectId: "chaseapp-8459b",
    storageBucket: "chaseapp-8459b.appspot.com",
    messagingSenderId: "1020122644146",
    appId: "1:1020122644146:web:68f163a80a77facbcc13ab",
    measurementId: "G-V87EKNP10J"
}))

/*
const db = firebase.firestore()
const webhooksRef = db.collection('webhooks')

const webPromise = webhooksRef.get().then((snapshot) => {
    snapshot.forEach((doc) => {
        if (doc.exists) {
            console.log(doc.data())
        }
    })
}).catch((error) => {
    console.error(error)
})

Promise.all([webPromise])
  .then(function() {
      process.exit(0);
  })
  .catch(function(error) {
      console.log("Transactions failed:", error);
      process.exit(1);
  });
let snapshot
let discords = []
getWebhooks()
console.log(discords)

async function getWebhooks () {
    snapshot = await webhooksRef.get()
    snapshot.forEach(doc => {
        if (doc.data().discord) {
            discords.push(doc.data().discord)
        }
    })
    return 1
}
 */

hook.setUsername('ChaseApp')
hook.setAvatar(IMAGE_URL)
hook.send("We have a Chase!")