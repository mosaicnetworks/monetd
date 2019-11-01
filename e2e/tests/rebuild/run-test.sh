#!/bin/bash

set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

NET=${1:-"rebuild"}
NETDIR=$HOME/.giverny/networks/$NET
TMP_DIR=/tmp/rebuild.$$
HOST="172.77.5.10"
PORT="8080"
ACCTS=10
TRANS=100

# Requires monetcli version 1.2.6 or later

mkdir -p $TMP_DIR

echo "Generate some transaction history"
echo ""

$mydir/../../tools/build-trans.sh -v --accounts=$ACCTS --transactions=$TRANS --faucet="Faucet" \
  --prefix=Test --node-name=Node --node-host=$HOST --node-port=$PORT \
  --config-dir=$HOME/.monettest --temp-dir=/tmp --faucet-config-dir=$NETDIR/keystore

ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi

addr0=$(sed -e 's/",.*$//g;s/^.*":"//g'  $NETDIR/keystore/node0.json)
addr1=$(sed -e 's/",.*$//g;s/^.*":"//g'  $NETDIR/keystore/node1.json)
addr2=$(sed -e 's/",.*$//g;s/^.*":"//g'  $NETDIR/keystore/node2.json)
addr3=$(sed -e 's/",.*$//g;s/^.*":"//g'  $NETDIR/keystore/node3.json)
addr4="6666666666666666666666666666666666666666" # Dummy address for an incomplete nomination

echo -e "\nVoting\n======\n\n"

# Nominate node 3
monetcli --datadir $NETDIR poa nominee new --pwd $NETDIR/keystore/node3.txt --moniker node3 --from node3 \
     -h $HOST -p $PORT $addr3

# Nominate dummy node 4
monetcli --datadir $NETDIR poa nominee new --pwd $NETDIR/keystore/node0.txt --moniker node4 --from node0 \
     -h $HOST -p $PORT $addr4   # Node 0 nominates our dummy node 4 as we have not generated a key

# Node 0 votes for node 3
monetcli --datadir $NETDIR poa nominee vote --pwd $NETDIR/keystore/node0.txt  --from node0 \
     -h $HOST -p $PORT --verdict true $addr3 -d

# Node 1 votes for node 3
monetcli --datadir $NETDIR poa nominee vote --pwd $NETDIR/keystore/node1.txt  --from node1 \
     -h $HOST -p $PORT --verdict true $addr3

# Node 2 votes for node 3
monetcli --datadir $NETDIR poa nominee vote --pwd $NETDIR/keystore/node2.txt  --from node2 \
     -h $HOST -p $PORT --verdict true $addr3

# Node 3 should be decided here

# Node 2 votes for node 4
monetcli --datadir $NETDIR poa nominee vote --pwd $NETDIR/keystore/node2.txt  --from node2 \
     -h $HOST -p $PORT --verdict true $addr4


# Pause to allow join blocks to be committed.
sleep 1

echo "Starting node 3"

giverny network push $NET node3

# We allow node3 time to join and sync with the other nodes
# We could replace this is a loop that monitors the completed rounds of all the 
# node and break when they all level out.
# But a simple sleep works too...

sleep 5 

$mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')

ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi


# Capture original files from network
wget -O $TMP_DIR/orig.genesis.json  http://$HOST:$PORT/export
ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi

wget -O $TMP_DIR/orig.peers.json  http://$HOST:$PORT/peers
ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi


# Amend config of node0 to maintenance mode 
NODE0IP=$(docker inspect --format="{{.NetworkSettings.Networks.$NET.IPAddress}}" node0)

echo "Node 0 IP is $NODE0IP"

docker cp node0:/.monet/monetd-config/monetd.toml $TMP_DIR/orig.monetd.toml
sed -e 's/maintenance-mode = "false"/maintenance-mode = "true"/g;s/bootstrap = false/bootstrap = true/g' $TMP_DIR/orig.monetd.toml > $TMP_DIR/maint.monetd.toml
docker cp $TMP_DIR/maint.monetd.toml node0:/.monet/monetd-config/monetd.toml 

# Stop Node0

docker exec node0 kill 1

sleep 5

# Start Node 0 in Maintenance mode

docker start node0

# Get files required to restart without loss of state

# Allow time to initialise the node
sleep 10

echo "Getting maintenance mode export"
wget -O $TMP_DIR/genesis.json  http://$NODE0IP:$PORT/export
ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi

wget -O $TMP_DIR/peers.json  http://$NODE0IP:$PORT/peers
ex=$?
if [ $ex -ne 0 ] ; then
    exit $ex
fi

# Generate account balances

echo "Before Balances"
for n in $(seq 1 $ACCTS)
do
  echo -n "."
  ADDR=$(sed -e 's/",.*$//g;s/^.*":"//g'  $HOME/.monettest/keystore/Test$n.json)
  wget -O $TMP_DIR/before$n.json  http://$NODE0IP:$PORT/account/$ADDR
done
echo ""

  cmp --silent $TMP_DIR/genesis.json $TMP_DIR/orig.genesis.json || ( echo "Genesis files are different" && exit 5)
  cmp --silent $TMP_DIR/peers.json $TMP_DIR/orig.peers.json || ( echo "Peers files are different")





# Overwrite docker config

find $NETDIR/docker -name genesis.json -exec cp $TMP_DIR/genesis.json {} \;
find $NETDIR/docker -name peers.json -exec cp $TMP_DIR/peers.json {} \;
find $NETDIR/docker -name peers.genesis.json -exec cp $TMP_DIR/peers.json {} \;


echo -e "\nRestart\n=======\n\n"


# Stop nodes

docker stop node3
docker stop node2
docker stop node1
docker stop node0

# Destroy the node

docker container prune -f

# Running nodes:
docker ps --all

giverny network push $NET node0
giverny network push $NET node1
giverny network push $NET node2
giverny network push $NET node3


sleep 5

 


# We compare the account totals that we generated earlier. 
# Generate account balances

echo "After Balances"
for n in $(seq 1 $ACCTS)
do
#  echo "Test$n"
  ADDR=$(sed -e 's/",.*$//g;s/^.*":"//g'  $HOME/.monettest/keystore/Test$n.json)
  wget -O $TMP_DIR/after$n.json  http://$HOST:$PORT/account/$ADDR
  cmp --silent $TMP_DIR/before$n.json $TMP_DIR/after$n.json || ( echo "Test $n  files are different" && exit 3)
  
done
echo ""


# $mydir/../../scripts/testlastblock.sh $( giverny network dump $NET | awk -F "|" '{print $2}')


#//TODO remove echo
rm -rf $TMP_DIR