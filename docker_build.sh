#!/bin/bash

docker container prune
docker image build -t forum .
docker container run -p 8080:8080 forum