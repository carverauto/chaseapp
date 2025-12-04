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

// firebase firestore
const db = firebase.firestore()
const adsbRef = db.collection('adsb')

const MS_PER_MIN = 60000;
const startDate = new Date(Date.now() - 4 * MS_PER_MIN)
const query = adsbRef.where('updated', '>=', startDate.getTime())

// firebase realtime db
const rtDb = firebase.database()
const rtDbRef = rtDb.ref()

// globals
let adsbObj = {}
let snapshot

exports.updateRTDB = (req, res) => {

    getADSB().then( foo => {
        const finalData = {}
        for (const [key, value] of Object.entries(adsbObj)) {
            finalData[key] = {
                geohash: value.geohash,
                group: value.group,
                imageUrl: value.imageUrl ? value.imageUrl : null,
                lat: value.lat,
                lon: value.lon,
                icao: value.hex,
                postime: value.updated,
                reg: value.tailno,
                tailno: value.tailno,
                trak: value.track,
                type: value.type
            }
            console.log(`Updating ${value.tailno}`)
        }

        const setPromise = rtDb.ref('adsb3/').set(finalData).then(r => console.log('Updated RTDB'))
        res.send('Updated RTDB')

        // Wait for all transactions to complete and then exit
        return Promise.all([foo, setPromise])
            .then(function() {
                process.exit(0);
            })
            .catch(function(error) {
                console.log("Transactions failed:", error);
                res.send('RTDB Update FAILED')
                process.exit(1);
            });
    })
};

async function getADSB () {
    snapshot = await query.get()
    snapshot.forEach(doc => {
        // writeRTDB(doc.data())
        const data = doc.data()
        const reg = data.tailno
        //console.log(reg)
        adsbObj[reg] = data
    })
}
