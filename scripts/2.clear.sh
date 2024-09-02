#!/bin/bash

echo "clear cl-postgres and chainlink"

docker rm -f chainlink
docker rm -f cl-postgres