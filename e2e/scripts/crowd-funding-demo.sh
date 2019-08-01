#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"net9"}
PORT=${2:-8080}

SOL_FILE="$mydir/../smart-contracts/CrowdFunding.sol"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$HOME/.giverny/networks/$NET/pwd.txt"

ips=($($mydir/get_running_nodes.sh | awk '{print $2}' | paste -sd "," -))

set -x
node crowd-funding/demo.js --ips=$ips \
    --port=$PORT \
    --contract=$SOL_FILE \
    --keystore=$KEY_DIR \
    --pwd=$PWD_FILE