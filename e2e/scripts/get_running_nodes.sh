#!/bin/bash

# Return the "[name] [ip]" list of running monetd nodes

NET=${1:-"net9"}

docker ps --filter network=$NET --format "{{.Names}}" | \
sort -u | \
xargs docker inspect -f "{{.Config.Hostname}} {{.NetworkSettings.Networks.$NET.IPAddress}}"