#!/bin/bash

ACCTMASK=$1
NET="prebaked"
GIVDIR="$HOME/.giverny/networks/$NET"
TRANS="$GIVDIR/trans"

# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


cnt=1
# echo trans_${ACCTMASK}_
# echo $TRANS/trans_${ACCTMASK}_*.json
for i in $TRANS/trans_${ACCTMASK}_*.json
do

    TRANSSET=$(basename $i .json)

    # Every 5th run, we hard reset the network to manage logs etc. 
    if [ $(($cnt % 5)) == 4 ] ; then
        giverny network stop --remove $NET  
        giverny network start $NET --use-existing --start-nodes -v
    else
        echo "($cnt) Running $TRANSSET"
        echo ""
        echo "Stopping Nodes $(docker stop Danu Jon Martin Kevin 2>&1)"
        echo "Starting Nodes $(docker start Danu Jon Martin Kevin 2>&1)"

        echo "Removing old DBs"
        for node in Danu Jon Martin Kevin
        do
            for j in $( docker exec $node ls  /.monet/babble )
            do
                if [ "$j" == "badger_db" ] ; then
                continue
                fi

                if [ "$j" == "${j##badger_}" ] ; then
                    continue
                fi

                docker exec $node rm -rf /.monet/babble/$j
            done
        done
    fi

    sleep 2

    $mydir/runprebaked.sh $TRANSSET
 cnt=$(( $cnt + 1 ))
done

echo "Finished"
