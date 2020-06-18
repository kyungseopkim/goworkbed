#!/bin/bash

version=v0.2
img=411026478373.dkr.ecr.us-east-1.amazonaws.com/kafka-filter:${version}

docker build -t $img .
docker-ecr-login
docker push $img
