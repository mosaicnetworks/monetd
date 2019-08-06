#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"transfers"}
PORT=${2:-8080}


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"


(node $mydir/index.js --datadir="$CONFIG_DIR" --nodeno=3) &
PIDS="$!"

( node $mydir/index.js --datadir="$CONFIG_DIR" --nodeno=2) &
PIDS="$PIDS $!"

( node $mydir/index.js --datadir="$CONFIG_DIR" --nodeno=1 ) &
PIDS="$PIDS $!"

( node $mydir/index.js --datadir="$CONFIG_DIR" --nodeno=0 ) &
PIDS="$PIDS $!"



FAIL=0
for job in $PIDS
do
    wait $job || let "FAIL+=1"
    echo $job $FAIL
done

echo $FAIL

if [ "$FAIL" == "0" ];
then
    echo "PASSED"
else
    echo "FAIL! ($FAIL)"
    exit 5
fi


$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')

# Test Balances
