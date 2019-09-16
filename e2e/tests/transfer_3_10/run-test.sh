#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"transfer_3_10"}

$mydir/../../tools/build-trans.sh -v --accounts=3 --transactions=10 --faucet="Faucet" \
  --prefix=Test --node-name=Node --node-host=172.77.5.10 --node-port=8080 \
  --config-dir=$HOME/.monettest --temp-dir=/tmp --faucet-config-dir=$HOME/.giverny/networks/$NET/keystore

exit $?
