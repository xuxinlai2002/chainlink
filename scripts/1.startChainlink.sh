#!/bin/bash

echo "start run chainlink node"

docker run --platform linux/x86_64/v8 --name chainlink \
-v /Users/xuxinlai/carv/work/chainlink/scripts/.chainlink-arbitrum-sepolia:/chainlink -it -p 6688:6688 \
--add-host=host.docker.internal:192.168.20.172 \
public.ecr.aws/chainlink/chainlink:2.12.0 node \
-config /chainlink/config.toml \
-secrets /chainlink/secrets.toml start -a /chainlink/.api
