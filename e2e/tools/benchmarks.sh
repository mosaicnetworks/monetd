#!/bin/bash


# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# cd to e2e folder
cd $mydir/..



# for i in 5 10 15 20 25 30 40 50 60 80 100 125 150 200 300 400
for i in 10 20 50 100 200
do
  scripts/start.sh --network=transfer_03_10 --init-ip=172.77.5.10 --heartbeat=${i}ms
  sleep 2
  tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 1"\
    --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
  sleep 2
  tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 2"\
    --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
  sleep 2
  tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 3"\
    --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
  make stop TEST=transfer_03_10
done
