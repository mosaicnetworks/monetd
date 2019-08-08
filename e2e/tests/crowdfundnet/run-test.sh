#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"crowdfundnet"}
PORT=${2:-8080}
OFFLINE=${3:-"true"}

if [ "$OFFLINE" != "" ] ; then
    OFFLINE="--offline true"
fi

SOL_FILE="$mydir/../../smart-contracts/CrowdFunding.sol"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"

ips=($(giverny network dump $NET | awk -F "|" '{print $2}' | paste -sd "," -))

set -x
node $mydir/../../crowd-funding/demo.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE  \
    $OFFLINE
