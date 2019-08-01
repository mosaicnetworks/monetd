#!/bin/bash

set -eu

NET=${1:-"net9"}

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

rm -rf $HOME/.giverny/networks/$NET

giverny network new $NET --initial-ip 172.77.5.10 --names $mydir/../../sampledata/names.txt --nodes 4 --save-pass

giverny network build $NET -v

giverny network start $NET --use-existing --start-nodes -v