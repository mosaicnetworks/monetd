#!/bin/bash

set -eu

# NET=${1:-"crowdfundnet"}
# INITIP=${2:-""}
# VERBOSE=${3:-""}

# CLI Defaults

NET="crowdfundnet"
INITIP=""
VERBOSE=""
BOOTSTRAP="false"
CACHESIZE=50000
HEARTBEAT="200ms"
MAXPOOL=2
SYNCLIMIT=1000
TIMEOUT="1s"
CACHE=128
MINGASPRICE=0

HELPTEXT="$0 [--network=\"$NET\"] [--init-ip=$INITIP] [--bootstrap=$BOOTSTRAP] \
[--cache-size=$CACHESIZE] [--heartbeat=$HEARTBEAT] [--max-pool=$MAXPOOL] \
[--sync-limit=$SYNCLIMIT] [--timeout=$TIMEOUT] [--cache=$CACHE] \
[--min-gas-price=$MINGASPRICE] [-v|--verbose] [-h|--help]"


# Parse all the command line options
while [ $# -gt 0 ]; do
  case "$1" in
    --network=*)
      NET="${1#*=}"
      ;;
    --init-ip=*)
      INITIP="${1#*=}"
      ;;

    --bootstrap=*)
      BOOTSTRAP="${1#*=}"
      ;;

    --cache-size=*)
      CACHESIZE="${1#*=}"
      ;;
    --heartbeat=*)
      HEARTBEAT="${1#*=}"
      ;;
    --max-pool=*)
      MAXPOOL="${1#*=}"
      ;;
    --sync-limit=*)
      SYNCLIMIT="${1#*=}"
      ;;
    --timeout=*)
      TIMEOUT="${1#*=}"
      ;;
    --cache=*)
      CACHE="${1#*=}"
      ;;
    --min-gas-price=*)
      MINGASPRICE="${1#*=}"
      ;;

    -v|--verbose)
      VERBOSE="true"
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

if [ "$INITIP" != "" ] ; then
    INITIP="--initial-ip $INITIP"
fi

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

rm -rf $HOME/.giverny/networks/$NET

giverny network new $NET \
    $INITIP \
    --names $mydir/../networks/$NET.txt \
    --pass $mydir/../networks/pwd.txt -v



if [ "$VERBOSE" != "" ] ; then
    sed --in-place "s/verbose = \"false\"/verbose = \"true\"/g" $HOME/.giverny/networks/$NET/monetd.toml
fi
 
sed --in-place "s/bootstrap = .*$/bootstrap = $BOOTSTRAP/g;s/cache-size = .*$/cache-size = $CACHESIZE/g; \
s/heartbeat = .*$/heartbeat = \"$HEARTBEAT\"/g;s/max-pool = .*$/max-pool = $MAXPOOL/g; \
s/sync-limit = .*$/sync-limit = $SYNCLIMIT/g;s/timeout = .*$/timeout = \"$TIMEOUT\"/g; \
s/cache = .*$/cache = $CACHE/g;s/min-gas-price = .*$/min-gas-price = $MINGASPRICE/g;" \
$HOME/.giverny/networks/$NET/monetd.toml 

giverny network build $NET -v

giverny network start $NET --use-existing -v

for node in $(giverny network dump $NET | grep "|true|" | cut -f1 -d'|')
do 
    giverny network push $NET $node -v
done 
