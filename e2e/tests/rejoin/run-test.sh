#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"rejoin"}


echo "Generate some transaction history"
echo ""

$mydir/../../tools/build-trans.sh -v --accounts=20 --transactions=400 --faucet="Faucet" \
  --prefix=Test --node-name=Node --node-host=172.77.5.10 --node-port=8080 \
  --config-dir=$HOME/.monettest --temp-dir=/tmp --faucet-config-dir=$HOME/.giverny/networks/$NET/keystore

ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi

# Sleep to allow the blocks to be committed.
sleep 2

$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')

ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi


echo "Killing node 3"
# monetd is always process 1. 
docker exec node3 kill 1


sleep 10


echo "Restarting node 3"
docker start node3

docker logs node3 2>&1 | tail -20 
sleep 10
echo ""
docker logs node3 2>&1 | tail -20 


$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')
