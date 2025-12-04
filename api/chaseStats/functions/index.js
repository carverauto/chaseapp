const Sentry = require("@sentry/serverless");
const functions = require("firebase-functions");
const admin = require("firebase-admin")
admin.initializeApp()

Sentry.GCPFunction.init({
    dsn: "https://261aa9bd2bd24afc9c0ca0ce161fd09a@o362496.ingest.sentry.io/4217236",
    tracesSampleRate: 1.0,
});

exports.helloHttp = Sentry.GCPFunction.wrapHttpFunction((req, res) => {
    throw new Error('oh, hello there!');
});

exports.chaseStats = functions.firestore
    .document("chases/{docId}").onWrite(change => {
        const after = change.after.data()
        const previous = change.before.data()


        /*
        if (after.Live !== previous.Live)
            if (!after.EndedAt) {
                if (previous.Live && (after.Live === false))
                    // If we were live and now we're not, that
                    // signals the end of the chase.
                    functions.logger.info("End of chase")
                return change.after.ref.set({
                    EndedAt: Date.now(),
                }, {merge: true});
            }

         */
        return 1;
    })
