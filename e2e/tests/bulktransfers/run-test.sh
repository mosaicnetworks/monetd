#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"bulktransfers"}
# CNT=${2:-30}
CNT=50

giverny transactions generate -n "$NET" --count "$CNT"


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
GIVDIR="$HOME/.giverny"

PRETOT=/tmp/pre.$$.json
POSTOT=/tmp/post.$$.json

node $mydir/index.js --network=$NET --account=faucet --totals=$PRETOT  --givdir="$GIVDIR" 



PIDS=""

for i in ${CONFIG_DIR}trans/*.json
do
    stub=$(basename $i .json)

    if [ "$stub" == "faucet" ] ; then
        continue
    fi
    if [ "$stub" == "delta" ] ; then
        continue
    fi
    if [ "$stub" == "trans" ] ; then
        continue
    fi


    echo $stub

    ( node $mydir/index.js --network=$NET --account=$stub  --givdir="$GIVDIR"  ) & PIDS="$PIDS $!"

done


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

# Test last balances...