#!/bin/bash


# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# cd to e2e folder
cd $mydir/..



# Parse all the command line options
while [ $# -gt 0 ]; do
 case "$1" in
    heartbeat)
        SF=$HOME/monet_heartbeat.txt
        # for i in 5 10 15 20 25 30 40 50 60 80 100 125 150 200 300 400
        for i in 10 20 50 100 200
        do
        scripts/start.sh --network=transfer_03_10 --init-ip=172.77.5.10 --heartbeat=${i}ms
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 1" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 2" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 3" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        make stop TEST=transfer_03_10
        done
    ;;
    maxpool)
        SF=$HOME/monet_maxpool.txt
        for i in 2 1 3 4
        do
        scripts/start.sh --network=transfer_03_10 --init-ip=172.77.5.10 --heartbeat=20ms --max-pool=$i
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 1" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 2" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=20 --transactions=1000 --stats-code="$i 3" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        make stop TEST=transfer_03_10
        done
    ;;
    accounts)
        SF=$HOME/monet_accounts.txt
        for i in 100 75 50 35 20 
        do
        scripts/start.sh --network=transfer_03_10 --init-ip=172.77.5.10 --heartbeat=20ms
        sleep 2
        tools/build-trans.sh --accounts=$i --transactions=1000 --stats-code="$i 1" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=$i --transactions=1000 --stats-code="$i 2" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        sleep 2
        tools/build-trans.sh --accounts=$i --transactions=1000 --stats-code="$i 3" --stats-file="$SF" \
            --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin
        make stop TEST=transfer_03_10
        done
    ;;
    dirtyaccounts)
        SF=$HOME/monet_accounts.txt
        for j in 1 2 3
        do
            scripts/start.sh --network=transfer_03_10 --init-ip=172.77.5.10 --heartbeat=20ms

            cnt=1    
            for i in $(printf "100\n300\n50\n250\n200\n150" | shuf)
            do
               sleep 2  
               tools/build-trans.sh --accounts=$i --transactions=1000 --stats-code="$cnt 1 dirty" --stats-file="$SF" \
                 --faucet-config-dir=/home/jon/.giverny/networks/transfer_03_10/keystore --round-robin

               cnt=$((cnt+1))  
            done

            make stop TEST=transfer_03_10
        done
        

    ;;

    *)
            echo "Not found"
            exit 1
    ;;



 esac 
 shift
done


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
