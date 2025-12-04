#!/bin/sh

DESC="LA Motorcycle chase"
NAME="LA motorcycle chase"
URL="https://abc7.com/watch/23340/"
LIVE="true"

curl --request POST -d '{"name":$NAME,"url":$URL,"live":$LIVE}' -H 'Content-Type: application/json' https://us-central1-chaseapp-8459b.cloudfunctions.net/AddChase
