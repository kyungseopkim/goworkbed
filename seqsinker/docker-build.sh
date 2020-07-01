#!/bin/bash

VERSION=v0.0.2
IMG=411026478373.dkr.ecr.us-east-1.amazonaws.com/mqtt-payload-metadata-sinker:${VERSION}
docker build -t $IMG .
docker-ecr-login
docker push $IMG
