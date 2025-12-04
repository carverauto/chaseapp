#curl --location --request POST 'https://air.aloft.ai/airspace-api/airspace' --http1.1 --header 'Authorization: Bearer a5nLCj7B0ClmBssGInW9QQDaz5HVfFgU0YPj0yxlTmLO09Yzse0XvsLZNd81prdji9np5avQ4Xh' --header 'Content-Type: application/json' --data-raw '{"geometry": {"type":"Point","coordinates":[-118.951721,32.75004]}}'

# curl --location --request POST 'https://air.aloft.ai/airspace-api/airspace' --http1.1 --header 'Authorization: Bearer a5nLCj7B0ClmBssGInW9QQDaz5HVfFgU0YPj0yxlTmLO09Yzse0XvsLZNd81prdji9np5avQ4Xh' --header 'Content-Type: application/json' --data-raw '{"geometry": {"type":"Polygon","coordinates":[[124.409591,32.534156],[-114.131211,42.009518]}}'

# curl --location --request POST 'https://air.aloft.ai/airspace-api/airspace' --http1.1 --header 'Authorization: Bearer a5nLCj7B0ClmBssGInW9QQDaz5HVfFgU0YPj0yxlTmLO09Yzse0XvsLZNd81prdji9np5avQ4Xh' --header 'Content-Type: application/json' -d '{"geometry": {"type":"Polygon","coordinates":[[124.409591,32.534156],[-114.131211,42.009518]]}}'

# curl --location --request POST 'https://air.aloft.ai/airspace-api/airspace' --http1.1 --header 'Authorization: Bearer a5nLCj7B0ClmBssGInW9QQDaz5HVfFgU0YPj0yxlTmLO09Yzse0XvsLZNd81prdji9np5avQ4Xh' --header 'Content-Type: application/json' -d '{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[124.409591,32.534156],[-114.131211,42.009518]]]}}'

curl --location --request POST 'https://air.aloft.ai/airspace-api/airspace' --http1.1 --header 'Authorization: Bearer a5nLCj7B0ClmBssGInW9QQDaz5HVfFgU0YPj0yxlTmLO09Yzse0XvsLZNd81prdji9np5avQ4Xh' --header 'Content-Type: application/json' -d @test.json
