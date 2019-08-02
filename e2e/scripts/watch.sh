#!/bin/bash
set -eux
NET=${1:-"net9"}
PORT=${2:-8080}

watch -t -n 1 '
docker ps --filter network='$NET' --format "{{.Names}}" | sort -u | \
xargs docker inspect -f "{{.NetworkSettings.Networks.'$NET'.IPAddress}}" | \
xargs -I % curl -s -m 1 http://%:'${PORT}'/stats | \
tr -d "{}\"" | \
awk -F "," '"'"'{gsub (/[,]/," "); print;}'"'"'
'