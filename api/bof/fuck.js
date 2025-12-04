import { createRequire } from "module"; // Bring in the ability to create the 'require' method
const require = createRequire(import.meta.url); // construct the require method
const serviceAccount = require("./google.json") // use the require method

import clustersDbscan from '@turf/clusters-dbscan';
import admin from 'firebase-admin';

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://chaseapp-8459b.firebaseio.com'
})

const db = admin.database();
const features = [];

let updatePromise, setPromise, fPromise;
const ref = db.ref('adsb');
const bofRef = db.ref('bof');

const adsbPromise = ref.once('value', function (snapshot) {
  const birdData = snapshot.val();
  // @ts-ignore
  // eslint-disable-next-line no-unused-vars
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
}).then(async (doc) => {
  if (features.length > 0) {
    const airships = {};
    airships.type = 'geojson';
    airships.data = {
      type: 'FeatureCollection',
      features,
    };
    // @ts-ignore
    const butthole = {
      fuck: "you",
      my: "butt"
    }
    const clusters = clustersDbscan(airships.data, 100, { minPoints: 2});
    // console.log(clusters.features)
    clusters.features.forEach((cluster) => {
      if (cluster.properties.dbscan !== 'noise') {
        console.log(`${cluster.properties.dbscan}`)
        const tailno = cluster.properties.title
        bofRef.child('clusters')

        const fPromise = bofRef.set({ tailno: {cluster} }).then(() => {
          console.log('fungali')
        }).catch((e) => {
          console.error(e)
        })

        bofRef.child(cluster.properties.title)
        bofRef.set(cluster.properties).then((foo) => {
          console.log(foo)
          console.log(`Updated ${cluster.properties.title} Cluster: ${cluster.properties.cluster}`)
        }).catch((error) => {
          console.error(error)
        })
      }
    })
  } else {
    console.log("Empty features list")
  }
})

// Wait for all transactions to complete and then exit
Promise.all([fPromise, updatePromise, setPromise, adsbPromise])
  .then(function () {
    process.exit(0);
  })
  .catch(function (error) {
    console.log('Transactions failed:', error);
    process.exit(1);
  });
