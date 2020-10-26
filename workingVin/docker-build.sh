#!/bin/bash

VERSION=v5.4
IMG=411026478373.dkr.ecr.us-east-1.amazonaws.com/influxdb-sinker:$VERSION
ecr-login
docker build -t $IMG .
docker push $IMG
