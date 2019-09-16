#!/bin/bash


URIFILE=$1

# METHOD="batch"
METHOD="notbatch"
# TMPFILE=/tmp/tmp.curl.$$



if [ "$METHOD" == "batch" ] ; then

    cmd="curl -s"

    grep http $URIFILE | {
    while read host post  
    do
    if [ "$cmd" != "curl -s" ] ; then 
        cmd="$cmd --next"
    fi 
    cmd="$cmd $host -d $post -X POST"
#    echo $cmd
    done
    $cmd
    } 
else
    grep http $URIFILE | {
    while read host post  
    do
    curl -s $host -d "$post" -X POST
    done
    }
fi

