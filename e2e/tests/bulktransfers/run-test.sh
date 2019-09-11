#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"bulktransfers"}
# CNT=${2:-30}
CNT=200

giverny transactions generate -n "$NET" --count "$CNT"


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
GIVDIR="$HOME/.giverny"

PRETOT=/tmp/pre.$$.json
POSTOT=/tmp/post.$$.json

node $mydir/index.js --network=$NET --account=faucet --totals=$PRETOT  --givdir="$GIVDIR" 


# Only time from here - as it is the actual test rather than preconfig.
res1=$(date +%s.%N)

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

res2=$(date +%s.%N)
dt=$(echo "$res2 - $res1" | bc)



echo $FAIL

echo "\n\n\n"

node $mydir/index.js --network=$NET --pretotals=$PRETOT  --givdir="$GIVDIR" 


if [ "$FAIL" == "0" ];
then
    echo "PASSED"
else
    echo "FAIL! ($FAIL)"
    exit 5
fi

echo "$CNT transactions took $dt seconds"
rate=$(echo "scale=4;$CNT / $dt" | bc)
echo "$rate transactions per second"
