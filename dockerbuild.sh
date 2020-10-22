#!/bin/bash -ev

docker-compose run --rm build
docker-compose up -d envoy

sleep 5
curl http://localhost:8000/get

docker-compose down