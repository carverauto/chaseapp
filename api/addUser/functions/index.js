// functions/index.js
const functions = require("firebase-functions");
const admin = require('firebase-admin')
admin.initializeApp()
const db = admin.firestore()

// discord webhook stuff
const { Webhook } = require('discord-webhook-node')
const hook = new Webhook("https://discord.com/api/webhooks/894594619776585728/djM8VDxWndDKJksV8HFxlvQMlwisDhOT0yUJRUgHjHbdHyecCfYMX7WVvKq4Sv9_v5Pa")
const IMAGE_URL = 'https://chaseapp.tv/icon.png'

/*
exports.CheckFCMResponses = functions.pubsub.topic('chases').onPublish( async (message) => {
    // check through FCM logs and see if we have any messages that weren't delivered
    // Remove the FCM tokens for the bad ones
    // admin.messaging().
})
 */

// trigger function on new user creation.
exports.AddUserRole = functions.auth.user().onCreate(async (authUser) => {
    if (authUser.uid)
        hook.setUsername('ChaseApp')
        hook.setAvatar(IMAGE_URL)
        await hook.send(`New user added ${authUser.uid}`)
});
