#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"rebuild"}


echo "Generate some transaction history"
echo ""

$mydir/../../tools/build-trans.sh -v --accounts=10 --transactions=100 --faucet="Faucet" \
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


# Authorise node 3

echo "Starting node 3"

docker start node3

# We allow node3 time to join and sync with the other nodes
# We could replace this is a loop that monitors the completed rounds of all the 
# node and break when they all level out.
# But a simple sleep works too...

# We start another election - just to verify the state carries over the reset

sleep 5 

# We get the current state from :8080/export
# We also get the peers from :8080/peers

# We then write these files to the network configuration

# Generate account balances

# We destroy the current docker nodes
 
# We create a new network with the new config

# We bring it up

# We compare the account totals that we generated earlier. 


$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')
