#!/bin/bash

set -eu

NET=${1:-"crowdfundnet"}
INITIP=${2:-""}

if [ "$INITIP" != "" ] ; then
    INITIP="--initial-ip $INITIP"
fi

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

rm -rf $HOME/.giverny/networks/$NET

giverny network new $NET \
    $INITIP \
    --names $mydir/../networks/$NET.txt \
    --pass $mydir/../networks/pwd.txt

giverny network build $NET -v

giverny network start $NET --use-existing -v


for node in $(giverny network dump $NET | grep "|true$" | cut -f1 -d'|')
do 
    giverny network push $NET $node
done 
