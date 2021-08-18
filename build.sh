#!/bin/sh

TAG_NAME=$(cat version.txt)

echo Building IPGeoLocation docker
docker build -t evgenis/ip-geo-location-app:${TAG_NAME} .