#!/bin/bash

cd docker/tinode/

echo "Stop and remove containers..."
docker-compose down
docker-compose up -d

echo "Containers success."
cd ../..
