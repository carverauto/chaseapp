import * as admin from 'firebase-admin';
import * as turf from '@turf/turf';
// eslint-disable-next-line no-unused-vars
import {GeojsonFeature, ADSB} from './types';

let adsbPromise: Promise<void>;
let updatePromise: Promise<void>;

const serviceAccount = require('./google.json');
admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://chaseapp-8459b.firebaseio.com',
});

const db = admin.database();
const features: GeojsonFeature[] = [];

exports.showBofs = (req: any, res: { send: (arg0: any) => void }) => {
  const ref = db.ref('bof/');
  ref.once('value', function(snapshot) {
    const data = snapshot.val();
    console.log(data);
    res.send(data);
  });
};

exports.findClusterz = (req: any, res: { send: (arg0: string) => void; }) => {
  const ref = db.ref('adsb');
  ref.once('value', function(snapshot) {
    const birdData = snapshot.val();
    // @ts-ignore
    // eslint-disable-next-line no-unused-vars
    for (const [key, value] of Object.entries(birdData)) {
      const item = value as ADSB;
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
          track: item.track,
          type: item.type,
          emergency: item.emergency,
        },
      } as GeojsonFeature;
      features.push(feature);
    }
  });

  const airships = {} as any;
  airships.type = 'geojson';
  airships.data = {
    type: 'FeatureCollection',
    features,
  };
  // @ts-ignore
  const clusters = turf.clustersDbscan(airships.data, 75);
  let seenMedia = 0;
  clusters.features.forEach((node) => {
    if (node.properties.group === 'media')
      seenMedia++;
  });

  if (seenMedia > 0) {
    console.log('Updating BOF');
    res.send('Updating BOF');
    const bofRef = ref.child('bof');
    bofRef.set(clusters);
  } else {
    // Upload an empty set
    console.log('no adsb pings from known media');
    // res.send('No ADSB pings from media recently')
    const data = {};
    const bofRef = ref.child('bof');
    bofRef.set(data);
    console.log('Updated BOF with empty set ');
  }
  // Wait for all transactions to complete and then exit
  return Promise.all([updatePromise, adsbPromise])
      .then(function() {
        process.exit(0);
      })
      .catch(function(error) {
        console.log('Transactions failed:', error);
        process.exit(1);
      });
};
