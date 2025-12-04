const functions = require("firebase-functions");
const admin = require('firebase-admin')
admin.initializeApp()
const fieldValue = admin.firestore.FieldValue

const dayjs = require('dayjs')
const relativeTime = require('dayjs/plugin/relativeTime')
const duration = require('dayjs/plugin/duration')

const {FieldValue} = require("@google-cloud/firestore/build/src");

// Lets store stuff in a document in the stats/ collection
// each document is named YYYY-MM
const dt = new Date()
const timePeriod = dt.getFullYear() + "-" + (dt.getMonth() + 1)

const db = admin.firestore()
const docRef = db.collection('stats').doc(timePeriod)

function dateDiff (createdAt,endedAt) {
    const x = dayjs(createdAt.toDate())
    const y = dayjs(endedAt.toDate())
    // Return seconds
    if (y) {
        console.log('We got y')
        return y.diff(x, 'seconds').humanize
    } else {
        return false
    }
}

// https://stackoverflow.com/questions/46554091/cloud-firestore-collection-count
exports.chaseCountListener = functions.firestore
    .document('chases/{documentUid}')
    .onWrite(async (change, context) => {

        if (!change.before.exists) {
            // New document Created : add one to count
            const res = await docRef.set({lastChaseDate: Date.now(), numberOfChases: fieldValue.increment(1)}, { merge: true} )
        } else if (change.before.exists && change.after.exists) {
            // Updating existing document : Do nothing
        } else if (!change.after.exists) {
            // Deleting document : subtract one from count
            const res = await docRef.set({numberOfChases: fieldValue.increment(-1)}, { merge: true} )
        }
        return
    });

exports.chaseDurationListener = functions.firestore
    .document('chases/{documentUid}')
    .onUpdate(async (change, context) => {
        const after = change.after.data()
        // const before = change.before.data()

       if (after.EndedAt) {
           console.log('In EndedAt')
           // Find the time (duration) from CreatedAt and EndedAt, maybe in ms?
           // const seconds = dateDiff(after.CreatedAt,after.EndedAt)
           const x = dayjs(after.CreatedAt.toDate())
           const y = dayjs(after.EndedAt.toDate())

           const seconds = y.diff(x, 'seconds')

           if (seconds) {
               console.log(`Seconds: ${seconds}`)

               // Add this duration to 'durations' array field
               const statsRef = db.collection('stats').doc(timePeriod)
               const unionRes = await statsRef.set({
                   durations: FieldValue.arrayUnion(seconds)
               }, { merge: true })
           } else {
               console.log(`Didn't get seconds: ${seconds}`)
           }
       }
       return
    })
