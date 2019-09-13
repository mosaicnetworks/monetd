#!/bin/bash



# CLI Params section. These will become parameters
VERBOSE="-v"             # EIther "" or "-v"
ACCTCNT=20                # Number of Accounts to transfer between       
TRANSCNT=200              # Total number of transactions 
FAUCET="Faucet"          # Faucet Account Moniker
PREFIX="Test"            # Prefix of the Moniker for transfer monikers   
NODENAME="Node0"         # Node Name
NODEHOST="172.77.5.10"   # Node IP
NODEPORT="8080"          # Node Port
CONFIGDIR="$HOME/.monettest" # Monet Config Dir
OUTDIRSTEM="/tmp"        # Output Directory

# Derived globals section

OUTDIR="$OUTDIRSTEM/Trans.$$"
SUFFIX=".json"
TRANSFILE=$OUTDIR/trans$SUFFIX

mkdir -p $OUTDIR


mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


res1=$(date +%s.%N)


# Generate Accounts to use for testing
giverny --monet-data-dir $CONFIGDIR keys generate \
    --prefix $PREFIX \ 
    --min-suffix 1 \
    --max-suffix $ACCTCNT \
    $VERBOSE



# Create expanded account list
ACCTS=""
for i in $(seq 1 $ACCTCNT)
do
    if [ "$ACCTS" != "" ] ; then
        ACCTS="$ACCTS,"
    fi
    ACCTS=$ACCTS$PREFIX$i
done



# Generate Transactions
giverny --monet-data-dir $CONFIGDIR transactions solo -v \
    --faucet $FAUCET \
    --accounts $ACCTS \
    --count $TRANSCNT \
    --output $TRANSFILE \
    $VERBOSE





# Process Faucet
node index.js --account=$FAUCET --nodename=$NODENAME --nodehost=$NODEHOST \
 --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR --outfile=$OUTDIR/$FAUCET$SUFFIX


for i in $(seq 1 $ACCTCNT)
do
    node index.js --account=$PREFIX$i --nodename=$NODENAME --nodehost=$NODEHOST \
    --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR --outfile=$OUTDIR/$PREFIX$i$SUFFIX
done


res2=$(date +%s.%N)


PIDS=""

for i in $(seq 1 $ACCTCNT)
do
    ( $mydir/run-trans.sh $OUTDIR/$PREFIX$i$SUFFIX  ) & PIDS="$PIDS $!"
done


FAIL=0
for job in $PIDS
do
    wait $job || let "FAIL+=1"
    echo $job $FAIL
done


# TImings

res3=$(date +%s.%N)
dt=$(echo "$res2 - $res1" | bc)
dt2=$(echo "$res3 - $res2" | bc)
echo "Preparing $TRANSCNT transactions took $dt seconds"
echo "$TRANSCNT transactions applying took $dt2 seconds"
rate=$(echo "scale=4;$TRANSCNT / $dt2" | bc)
echo "$rate transactions per second"


if [ "$FAIL" == "0" ];
then
    echo "PASSED"
else
    echo "FAIL! ($FAIL)"
    exit 5
fi

