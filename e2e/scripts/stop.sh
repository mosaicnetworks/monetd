#!/bin/bash

NET=${1:-"net9"}

docker ps --filter network=$NET --format "{{.Names}}" -aq | xargs docker rm -f
docker network rm $NET
