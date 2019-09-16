#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"transfer_10_200"}

$mydir/../../tools/build-trans.sh -v --accounts=10 --transactions=200 --faucet="Faucet" \
  --prefix=Test --node-name=Node --node-host=172.77.5.10 --node-port=8080 \
  --config-dir=$HOME/.monettest --temp-dir=/tmp --faucet-config-dir=$HOME/.giverny/networks/$NET/keystore

ex=$?

if [ $ex -ne 0 ] ; then
    exit $?
fi

# Sleep to allow the blocks to be committed.
  sleep 2

  $mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')

  exit $?  