#!/bin/bash
while getopts v: flag
do
  case "${flag}" in
    v) version=${OPTARG};;
  esac
done
CONTAINERNAME=icedoor_go
echo "Building $CONTAINERNAME:$version"
docker buildx build --platform linux/amd64 --load -t cantimages.azurecr.io/cantitlan/$CONTAINERNAME:$version .
echo "Pushed $CONTAINERNAME:$version"
