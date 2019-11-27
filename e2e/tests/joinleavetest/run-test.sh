#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"joinleavetest"}
PORT=${2:-8080}


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"

# The network initially contains nodes 0, 1, 2, and 4
# The following script causes node3 to be added to the whitelist, and node4 to
# be removed
node $mydir/index.js --datadir="$CONFIG_DIR" --action="join"
ret=$?

# node3 was added to the whitelist, so should be allowed to join
giverny network push $NET node3 -v

sleep 10

lastinfo=""
lastip=""

exitcode=0

# nodes 0, 1, 2, and 3 should have the same block index. Node4 should be 
# Suspended.
for n in $( giverny network dump $NET | awk -F "|" '{print $2}')  
do
   url="http://$n:8080/info"
   thisinfo=$(curl -s $url  | jq .last_block_index)
   node=$(curl -s $url  | jq .moniker)
   state=$(curl -s $url  | jq .state)

   if [ "$node" == "\"node4\"" ] ; then
      if [ "$state" != "\"Suspended\"" ] ; then
         exitcode=201
         echo "node4 should be suspended"
         echo "node4 state: $state"
      fi
      continue
   fi

   if [ "$lastip" = "" ] ; then 
      lastip=$n
      lastinfo="$thisinfo"
      continue
   fi

   if [ "$lastinfo" != "$thisinfo" ] ; then
      exitcode=201
      echo "$lastip info does not match $n info"
      echo "$lastip: $lastinfo"
      echo "$n: $thisinfo"
   fi
   lastip=$n
   lastinfo="$thisinfo"   
done

if [ "$exitcode" != "0" ] ; then
   exit $exitcode
fi

$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')
