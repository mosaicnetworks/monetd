#!/bin/bash

# CLI Params default section
VERBOSE=""                    # EIther "" or "-v"
ACCTCNT=10                      # Number of Accounts to transfer between       
TRANSCNT=200                    # Total number of transactions 
FAUCET="Faucet"                 # Faucet Account Moniker
FAUCETCONFIG=""                 # Keystore to copy faucet key from
PREFIX="Test"                   # Prefix of the Moniker for transfer monikers   
NODENAME="Node"                 # Node Name
NODEHOST="172.77.5.11"          # Node IP
NODEPORT="8080"                 # Node Port
CONFIGDIR="$HOME/.monettest"    # Monet Config Dir used for this test
OUTDIRSTEM="/tmp"               # Output Directory
ROUNDROBIN=""                   # Round Robin Transaction generation
STATSCODE=""                    # Log Stats Under This Code
STATSFILE="$HOME/monetstats.txt"  # Stats File
MINACCT=1001



# Populate the helptext including the default values
HELPTEXT="$0 [-v] [--accounts=$ACCTCNT] [--transactions=$TRANSCNT] [--faucet=\"$FAUCET\"] \
 [--faucet-config-dir=$FAUCETCONFIG] [--prefix=$PREFIX] [--node-name=$NODENAME] [--node-host=$NODEHOST]\
 [--node-port=$NODEPORT] [--config-dir=$CONFIGDIR] [--temp-dir=$OUTDIRSTEM] \
 [--stats-code=$STATSCODE] [--stats-file=$STATSFILE] [--round-robin] [-h|--help]"


# Parse all the command line options
while [ $# -gt 0 ]; do
  case "$1" in
    -v)
      VERBOSE="-v"
      ;;
    --accounts=*)
      ACCTCNT="${1#*=}"
      ;;
    --transactions=*)
      TRANSCNT="${1#*=}"
      ;;
    --faucet=*)
      FAUCET="${1#*=}"
      ;;
    --faucet-config-dir=*)
      FAUCETCONFIG="${1#*=}"
      ;;
    --prefix=*)
      PREFIX="${1#*=}"
      ;;
    --node-name=*)
      NODENAME="${1#*=}"
      ;;
    --node-host=*)
      NODEHOST="${1#*=}"
      ;;
    --node-port=*)
      NODEPORT="${1#*=}"
      ;;
    --config-dir=*)
      CONFIGDIR="${1#*=}"
      ;;
    --temp-dir=*)
      OUTDIRSTEM="${1#*=}"
      ;;
    --stats-file=*)
      STATSFILE="${1#*=}"
      ;;
    --stats-code=*)
      STATSCODE="${1#*=}"
      ;;
    --round-robin)
      ROUNDROBIN="--round-robin"
      ;;  
    -h|--help)
      echo $HELPTEXT
      exit 0
      ;;


    *)
      printf "***************************\n"
      printf "* Error: Invalid argument.*\n"
      printf "***************************\n"
      echo $HELPTEXT
      exit 1
  esac
  shift
done

# Output the options used in command line format to aid debugging.

echo "$0 $VERBOSE --accounts=$ACCTCNT --transactions=$TRANSCNT --faucet=\"$FAUCET\" \
 --faucet-config-dir=$FAUCETCONFIG --prefix=$PREFIX --node-name=$NODENAME \
 --node-host=$NODEHOST --node-port=$NODEPORT $ROUNDROBIN \
 --config-dir=$CONFIGDIR --temp-dir=$OUTDIRSTEM \
 --stats-code=\"$STATSCODE\" --stats-file=$STATSFILE "


# Derived globals section

OUTDIR="$OUTDIRSTEM/trans.$$"
SUFFIX=".json"
TRANSFILE=$OUTDIR/trans$SUFFIX
PRE=$OUTDIR/pre.json

# Make directories we will write files to if they do not already exist
mkdir -p $OUTDIR
mkdir -p $CONFIGDIR/keystore

# If FAUCETCONFIG is set, then we copy the faucet account pair to CONFIGDIR

if [ "$FAUCETCONFIG" != "" ] ; then
    if [ ! -f "$FAUCETCONFIG/keystore/$FAUCET.json" ] ; then
        cp $FAUCETCONFIG/$FAUCET.json  $CONFIGDIR/keystore

    fi

    if [ ! -f "$FAUCETCONFIG/keystore/$FAUCET.txt" ] ; then
        cp $FAUCETCONFIG/$FAUCET.txt  $CONFIGDIR/keystore
    fi
fi

# Verify that we have a key pair for the Faucet. 
if [ ! -f "$CONFIGDIR/keystore/$FAUCET.json" ] ; then
    echo "You need to copy $FAUCET.json into $CONFIGDIR/keystore"
    exit 1
fi

if [ ! -f "$CONFIGDIR/keystore/$FAUCET.txt" ] ; then
    echo "You need to copy $FAUCET.txt into $CONFIGDIR/keystore"
    exit 1
fi


# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# Time Stamp 1
res1=$(date +%s.%N)


# Generate Accounts to use for testing
giverny --monet-data-dir $CONFIGDIR keys generate \
    --prefix $PREFIX \
    --min-suffix $MINACCT \
    --max-suffix $ACCTCNT \
    $VERBOSE 

# Create expanded account list
ACCTS=""
for i in $(seq $MINACCT $ACCTCNT)
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
    $VERBOSE \
    $ROUNDROBIN




# Get Peers List
peers=()
for peer in $(curl -s http://$NODEHOST:$NODEPORT/peers | jq ".[] | .NetAddr" | sed -e 's/"//g;s/:1337//g')
do
   peers+=( $peer )
done
numpeers=${#peers[@]}


# Process Faucet
node $mydir/index.js --account=$FAUCET --nodename=${NODENAME}0 --nodehost=$NODEHOST \
 --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR  --total=$PRE
     # --outfile=$OUTDIR/$FAUCET$SUFFIX

exitcode=$?

if [ $exitcode -ne 0 ] ; then
    echo "Faucet Allocation failed."
    exit $exitcode
fi


if [ ! -f "$PRE" ] ; then
    echo "Funding from $FAUCET failed. Aborting."
    exit 6
fi

# Generate Signed Transactions
for i in $(seq 1 $ACCTCNT)
do
    node $mydir/index.js --account=$PREFIX$i --nodename=${NODENAME}$(($i % $numpeers))  --nodehost=${peers[$(($i % $numpeers))]} \
    --nodeport=$NODEPORT --transfile=$TRANSFILE --configdir=$CONFIGDIR --outfile=$OUTDIR/$PREFIX$i$SUFFIX

    exitcode=$?

    if [ $exitcode -ne 0 ] ; then
        echo "Transaction signing for $PREFIX$i failed."
        exit $exitcode
    fi
done





echo ""
echo "Starting timed section ($TRANSCNT transactions)"
echo ""

printf '%*s' $TRANSCNT | tr ' ' '*' 
printf "\e[${TRANSCNT}D\e[$(( $TRANSCNT / $(tput cols) ))A"

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
node $mydir/index.js --account=$FAUCET --nodename=$NODENAME --nodehost=$NODEHOST \
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

echo "Temporary files in $OUTDIR have not been deleted."