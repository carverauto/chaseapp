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
const db = firebase.firestore()
let docRef = db.collection("airships/")

docRef.get().then((doc) => {
    if (doc.exists) {
        console.log("Document data:", doc.data());
    } else {
        // doc.data() will be undefined in this case
        console.log("No such document!");
    }
}).catch((error) => {
    console.log("Error getting document:", error);
});

/*
dbRef.child("airships/").get().then((snapshot) => {
    let bofData = snapshot.val()
    console.log(bofData)
}).catch((err) => {
    console.error(err)
})


 */
/*
// Wait for all transactions to complete and then exit
return Promise.all([bofPromise])
    .then(function() {
        process.exit(0);
    })
    .catch(function(error) {
        console.log("Transactions failed:", error);
        process.exit(1);
    });


 */
