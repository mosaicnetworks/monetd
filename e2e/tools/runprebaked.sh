#!/bin/bash


TRANSSET=$1
STATSCODE=$TRANSSET


STATSFILE="$HOME/prebakedstats.txt"
NET="prebaked"
GIVDIR="$HOME/.giverny/networks/$NET"
TRANS="$GIVDIR/trans"
PREFIX="Test"
SUFFIX=".json"
FAUCET="Faucet"
NODEHOST="172.77.5.10"
NODEPORT="8080"
NODENAME="Node0"
OUTDIR=$TRANS/$TRANSSET

STEM=$(basename $TRANSSET)
STEM=${STEM##trans_}
ACCTCNT=${STEM%%_*}
TRANSCNT=${STEM##*_}


TRANSFILE=""
CONFIGDIR=""
PRE=""

# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


if [ ! -d "$OUTDIR" ] ; then
    >&2 echo "Cannot find $OUTDIR"
    exit 1
fi

# Start Trans Timestamp
res2=$(date +%s.%N)


# Launch signed transactions processing as a background process
PIDS=""
for i in $(seq 1 $ACCTCNT)
do
    ( $mydir/run-trans.sh $OUTDIR/$PREFIX$i$SUFFIX  ) & PIDS="$PIDS $!"
done

# Wait for background tasks to finish
FAIL=0
for job in $PIDS
do
    wait $job || let "FAIL+=1"
done

echo ""

# Timings
# Finish timer
res3=$(date +%s.%N)
dt=$(echo "$res2 - $res1" | bc)
dt2=$(echo "$res3 - $res2" | bc)


# Check values of accounts as expected
echo node $mydir/index.js --account=$FAUCET --nodename=$NODENAME --nodehost=$NODEHOST \
 --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR  --pre=$PRE
exitcode=$?

echo "Preparing $TRANSCNT transactions took $dt seconds"
echo "$TRANSCNT transactions applying took $dt2 seconds"
rate=$(echo "scale=4;$TRANSCNT / $dt2" | bc)
echo "$rate transactions per second"


if [ $exitcode -ne 0 ] ; then
    echo "Balance checks failed."
    exit $exitcode
fi


if [ "$FAIL" == "0" ];
then
    echo "PASSED"
        
    if [ ! -z "$STATSCODE" ] ; then
      echo "$TRANSCNT $ACCTCNT $dt2 $STATSCODE" >> $STATSFILE
    fi  


else
    echo "FAIL! ($FAIL)"
    exit 5
fi