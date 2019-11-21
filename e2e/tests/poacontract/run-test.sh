#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"poacontract"}
PORT=${2:-8080}


CONFIG_DIR="$HOME/.giverny/networks/$NET/"
KEY_DIR="$HOME/.giverny/networks/$NET/keystore/"
PWD_FILE="$mydir/../../networks/pwd.txt"
MONETCLI="monetcli --datadir $CONFIG_DIR "


nodes=()
for i in $HOME/.giverny/networks/poacontract/keystore/*.json
do  
   ADDR=$(cat $i | jq ".address" | sed -e 's/"//g') 
   NODENO=$(basename $i .json | sed -e 's/^node//g')
   nodes[$NODENO]=$ADDR
done


function statusCheck() {
   EXPECTED=$1

   WHITE=$($MONETCLI  poa whitelist --silent --json | jq length)  # --silent
   NOM=$($MONETCLI  poa nominee list --silent --json    | jq length)
   EVICT=$($MONETCLI  poa evictee list --silent --json    | jq length)

   ACTUAL="${WHITE}_${NOM}_${EVICT}"

   if [ "$ACTUAL" != "$EXPECTED" ] ; then
      echo "Expected $EXPECTED, got $ACTUAL"
      return 9
   fi
   return 0
}



# First we set up monetcli
$MONETCLI config set --host 172.77.5.10 --from node0
statusCheck 3_0_0 || exit 11

# - Node 3 self nominates

$MONETCLI poa nominee new --pwd $mydir/../../networks/pwd.txt  --silent --moniker "node3" --from "node3" ${nodes[3]}
$MONETCLI poa nominee list

statusCheck 3_1_0 || exit 12



#    - Node 0 votes Yes
$MONETCLI poa nominee vote --verdict yes --from node0  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 3_1_0 || exit 13


#    - Node 0 tries to vote again - should fail
$MONETCLI poa nominee vote --verdict yes --from node0  --silent --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 3_1_0 || exit 14


#    - Node 3 attempt to self-nominate - should fail
$MONETCLI poa nominee new --pwd $mydir/../../networks/pwd.txt  --silent  --moniker "node3" --from "node3" ${nodes[3]}

statusCheck 3_1_0 || exit 15


#    - Node 1 votes No
$MONETCLI poa nominee vote --verdict no --from node1 --pwd $mydir/../../networks/pwd.txt ${nodes[3]}


$MONETCLI poa nominee list

statusCheck 3_0_0 || exit 16


# - Node 3 self nominates 
$MONETCLI poa nominee new --pwd $mydir/../../networks/pwd.txt  --silent  --moniker "node3" --from "node3" ${nodes[3]}

statusCheck 3_1_0 || exit 17


#    - Node 0 votes yes - this is the current point of failure
$MONETCLI poa nominee vote --verdict yes --from node0  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 3_1_0 || exit 18


#    - Node 1 votes yes
$MONETCLI poa nominee vote --verdict yes --from node1  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 3_1_0 || exit 19


#    - Node 2 votes yes - should be approved
$MONETCLI poa nominee vote --verdict yes --from node2  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 4_0_0 || exit 20


$MONETCLI poa whitelist
$MONETCLI poa nominee list
$MONETCLI poa evictee list

#    
# - Node 0 nominates Node 4 for eviction

$MONETCLI poa evictee new --pwd $mydir/../../networks/pwd.txt  --silent  --from "node0" ${nodes[3]}

statusCheck 4_0_1 || exit 21

#    - Node 0 votes to evict Node 3
$MONETCLI poa evictee vote --verdict yes --from node0  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 4_0_1 || exit 22

#    - Node 0 attempts to re-nominate node 3 for eviction - should fail
$MONETCLI poa evictee new --pwd $mydir/../../networks/pwd.txt  --silent  --from "node0" ${nodes[3]}

statusCheck 4_0_1 || exit 23

#    - Node 1 votes No

$MONETCLI poa evictee list

$MONETCLI poa evictee vote --verdict no --from node1  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

$MONETCLI poa evictee list

statusCheck 4_0_0 || exit 24

#
# - Node 0 nominates Node 4 for eviction
$MONETCLI poa evictee new --pwd $mydir/../../networks/pwd.txt  --silent  --from "node0" ${nodes[3]}

statusCheck 4_0_1 || exit 25

#    - Node 0 votes yes - this is the current point of failure
$MONETCLI poa evictee vote --verdict yes --from node0  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}

statusCheck 4_0_1 || exit 26

#    - Node 1 votes yes
$MONETCLI poa evictee vote --verdict yes --from node1  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}


statusCheck 4_0_1 || exit 27

#    - Node 2 votes yes - should be approved
$MONETCLI poa evictee vote --verdict yes --from node2  --silent  --pwd $mydir/../../networks/pwd.txt ${nodes[3]}


statusCheck 3_0_0 || exit 28

exit 0




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
