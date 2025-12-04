import { createRequire } from "module"; // Bring in the ability to create the 'require' method
const require = createRequire(import.meta.url); // construct the require method
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

// Cluster airships using DBSCAN and haversine algo for distance check
import clustersDbscan from '@turf/clusters-dbscan';
const adsbRef = firebase.database().ref()

// globals
const features = [];
let fPromise, setPromise

var exports = {}
exports.findClusters = (req, res) => {
    fPromise = adsbRef.child('adsb').get().then((snapshot) => {
        if (snapshot.exists()) {
            const birdData = snapshot.val()
            for (const [key, value] of Object.entries(birdData)) {
                const item = value;
                const feature = {
                    type: 'Feature',
                    geometry: {
                        type: 'Point',
                        // Mapbox wants things in long,lat
                        coordinates: [item.lon, item.lat],
                    },
                    properties: {
                        title: item.tailno,
                        group: item.group,
                        imageUrl: item.imageUrl,
                        type: item.type,
                    },
                };
                features.push(feature);
            }
        }
        if (features.length > 0) {
            const airships = {};
            airships.type = 'geojson';
            airships.data = {
                type: 'FeatureCollection',
                features,
            };
            const clusters = clustersDbscan(airships.data, 2, {minPoints: 2});
            if (clusters.features.length > 0) {
                console.log(clusters.features)
                setPromise = firebase.database().ref('bof').set(clusters.features)
            } else {
                console.log('Error: clusters is empty and should never be')
            }
        } else {
            console.log('Features list empty')
        }
    })

    Promise.all([fPromise, setPromise])
      .then(function() {
          process.exit(0);
      })
      .catch(function(error) {
          console.log("Transactions failed:", error);
          process.exit(1);
      });
}
