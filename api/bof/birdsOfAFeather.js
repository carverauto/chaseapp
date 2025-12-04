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

let seenMedia = []
let seenLeo = []
let birdInfo = {}

const rtdb = firebase.database()
const dbRef = rtdb.ref()
let snapshot
let results

const adsbPromise = dbRef.child("adsb").get().then((snapshot) => {
    let birdData = snapshot.val()
    for (const [key, value] of Object.entries(birdData)) {
        if (value.group === 'media') {
            seenMedia.push(key)
            birdInfo[key] = value
        }
        if (value.group === 'leo') {
            // console.log(`LEO Key: ${key}`)
            seenLeo.push(key)
            birdInfo[key] = value
        }
    }

    if (seenMedia.length > 0) {
        let media = 0
        let leo = 0
        let myBof = {}
        console.log(`[${seenMedia}]`)
        seenMedia.forEach( mediaBird => {
            // console.log(`${mediaBird}`)
            seenLeo.forEach( leoBird => {
                if (birdInfo[leoBird]) {
                    // console.log(birdInfo[leoBird])
                    let d = checkDistance(birdInfo[seenMedia[media]].lat,birdInfo[seenMedia[media]].lon,birdInfo[leoBird].lat,birdInfo[leoBird].lon)
                    if (( d / 1609.34) <= 100) {
                        let mediaTailNo = birdInfo[seenMedia[media]].tailno
                        let leoTailNo = birdInfo[leoBird].tailno
                        // console.log(birdInfo[leoBird].tailno)
                        // results[birdInfo[seenMedia[media]].tailno] = { tailno: birdInfo[leoBird].tailno, distance: Math.round(d / 1609.34) }
                        const mediaNearAirport = nearAirport(birdInfo[seenMedia[media]].lat,birdInfo[seenMedia[media]].lon)
                        const leoNearAirport = nearAirport(birdInfo[leoBird].lat,birdInfo[leoBird].lon)
                        console.log(`Distance between ${birdInfo[seenMedia[media]].tailno} and ${birdInfo[leoBird].tailno} - ${Math.round(d / 1609.34)}/mi`)
                    }
                }
                leo++
            })
            media++
        })
        // console.log(myBof)
        //rtdb.ref('bof/').set(myBof)
    } else {
        console.log('no adsb pings from known media')
        return 0
    }
}).catch((err) => {
    console.error(err)
})

// Wait for all transactions to complete and then exit
return Promise.all([adsbPromise, snapshot])
    .then(function() {
        process.exit(0);
    })
    .catch(function(error) {
        console.log("Transactions failed:", error);
        process.exit(1);
    });


// nearAirport - check to see if a bird is near an airport
function nearAirport(lat,lon) {
    // console.log(`${lat} ${lon}`)
    return 1
}

// checkDistance will return the distance (in meters) between two sets of
// coordinates, using the haversine formula.
function checkDistance (lat1,lon1,lat2,lon2) {
    const R = 6371e3    // meters
    const φ1 = lat1 * Math.PI/180   // φ, λ in radians
    const φ2 = lat2 * Math.PI/180
    const Δφ = (lat2-lat1) * Math.PI/180
    const Δλ = (lon2-lon1) * Math.PI/180

    const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
        Math.cos(φ1) * Math.cos(φ2) *
        Math.sin(Δλ/2) * Math.sin(Δλ/2)
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))

    return R * c    // in meters
}


function convertTZ(date, tzString) {
    return new Date((typeof date === "string" ? new Date(date) : date).toLocaleString("en-US", {timeZone: tzString}));
}
