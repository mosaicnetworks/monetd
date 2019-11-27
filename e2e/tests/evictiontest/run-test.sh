#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"evictiontest"}
PORT=${2:-8080}

CONFIG_DIR="$HOME/.giverny/networks/$NET/"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"

# The network starts with two nodes, A and B, both whitelisted.

# The script tells A to evict B from the whitelist. The core of this test is to
# check that B was also automatically evicted from the Babble validator-set and
# suspended.
node $mydir/index.js --datadir="$CONFIG_DIR"
ret=$?
$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')
