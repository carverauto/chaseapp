curl --location --request POST 'http://127.0.0.1:5001/chaseapp-8459b/us-central1/smallestSurroundingRectangleByArea' --http1.1 \
        --header 'Content-Type: application/json' \
        -d '{ "type":"Feature", "geometry": { "type":"LineString", "coordinates": [ [-74.020514, 40.71041], [-74.01103, 40.71311], [-74.018712, 40.717079]]}}'
