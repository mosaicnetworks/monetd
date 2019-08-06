#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"jointest"}
PORT=${2:-8080}


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"


node $mydir/index.js --datadir="$CONFIG_DIR" --action="join"
ret=$?


giverny network push $NET node3 -v


sleep 10

lastinfo=""
lastip=""

exitcode=0

for n in $( giverny network dump $NET | awk -F "|" '{print $2}')  
do
   url="http://$n:8080/info"
   thisinfo=$(curl -s $url  | json_pp | grep  \"last_block_index\" | sed -e 's/,//g;s/[\t ]//g' )

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
