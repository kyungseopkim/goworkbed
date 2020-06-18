#!/bin/bash

VERSION=v0.7
IMG=411026478373.dkr.ecr.us-east-1.amazonaws.com/influxdb-sinker:$VERSION
docker build -t $IMG .
docker-ecr-login
docker push $IMG
