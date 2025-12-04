const functions = require("firebase-functions");
const admin = require("firebase-admin")
admin.initializeApp()

exports.chaseEnded = functions.firestore
    .document("chases/{docId}").onUpdate(change => {
        const after = change.after.data()
        const before = change.before.data()

        if (before.Live !== after.Live)
          if (before.Live)
            if (!after.Live)
              functions.logger.info("End of chase")
              return change.after.ref.set({
                EndedAt: admin.firestore.Timestamp.now(),
              }, {merge: true});
    })
  return 1;
