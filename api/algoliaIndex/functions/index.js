const functions = require("firebase-functions");
const algoliasearch = require("algoliasearch");

// Initialize Algolia, requires installing Algolia dependencies:
// https://www.algolia.com/doc/api-client/javascript/getting-started/#install
//
// App ID and API Key are stored in functions config variables
const ALGOLIA_ID = functions.config().algolia.app_id;
const ALGOLIA_ADMIN_KEY = functions.config().algolia.api_key;
// const ALGOLIA_SEARCH_KEY = functions.config().algolia.search_key;

const ALGOLIA_INDEX_NAME = "chases";
const client = algoliasearch(ALGOLIA_ID, ALGOLIA_ADMIN_KEY);

// Update the search index every time a blog post is written.
exports.updateAlgolia = functions.firestore.document("chases/{docId}")
    .onCreate((snap, context) => {
      // Get the chase document
        const chase = snap.data();
          // Add an 'objectID' field which Algolia requires
          // chase.objectID = context.params.chaseId;
        chase.objectID = snap.id
        // Write to the algolia index
        console.log(chase)
        const index = client.initIndex(ALGOLIA_INDEX_NAME);
        return index.saveObject(chase);
    });
