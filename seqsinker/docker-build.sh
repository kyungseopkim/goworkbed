#!/bin/bash

VERSION=v0.0.3
IMG=411026478373.dkr.ecr.us-east-1.amazonaws.com/mqtt-payload-metadata-sinker:${VERSION}
ecr-login
docker build -t $IMG .
docker push $IMG
