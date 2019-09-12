#!/bin/bash



# CLI Params section. These will become parameters
VERBOSE="-v"
ACCTCNT=3
TRANSCNT=12
FAUCET="Faucet"
PREFIX="Test"
NODENAME="Node0"
NODEHOST="172.77.5.10"
NODEPORT="8080"
CONFIGDIR="$HOME/.monet"
OUTDIRSTEM="/tmp"

# Derived globals section

OUTDIR="$OUTDIRSTEM/Trans.$$"
SUFFIX=".json"
TRANSFILE=$OUTDIR/trans$SUFFIX

mkdir -p $OUTDIR


mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"




# Generate Accounts to use for testing
giverny keys generate --prefix $PREFIX --min-suffix 1 --max-suffix $ACCTCNT $VERBOSE


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
giverny transactions solo -v --faucet $FAUCET --accounts $ACCTS   \
--count $TRANSCNT $VERBOSE --output $TRANSFILE


# Process Faucet
node index.js --account=$FAUCET --nodename=$NODENAME --nodehost=$NODEHOST \
 --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR --outfile=$OUTDIR/$FAUCET$SUFFIX


for i in $(seq 1 $ACCTCNT)
do
    node index.js --account=$PREFIX$i --nodename=$NODENAME --nodehost=$NODEHOST \
    --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR --outfile=$OUTDIR/$PREFIX$i$SUFFIX
done



for i in $(seq 1 $ACCTCNT)
do
   $mydir/run-trans.sh $OUTDIR/$PREFIX$i$SUFFIX
done


