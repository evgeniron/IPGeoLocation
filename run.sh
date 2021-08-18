#!/bin/sh

TAG=$(cat version.txt)
port=$(grep APP_PORT env | cut -d'=' -f2)

docker run --rm --env-file env -p $port:$port evgenis/ip-geo-location-app:$TAG