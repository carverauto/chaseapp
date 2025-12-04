curl --location --request POST 'https://us-central1-chaseapp-8459b.cloudfunctions.net/smallestSurroundingRectangleByArea' --http1.1 \
        --header 'Content-Type: application/json' \
        -d '{ "type":"Feature", "geometry": { "type":"LineString", "coordinates": [ [-74.020514, 40.71041], [-74.01103, 40.71311], [-74.018712, 40.717079]]}}'
