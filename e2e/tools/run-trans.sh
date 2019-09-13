#!/bin/bash


URIFILE=$1
METHOD="batch"
# TMPFILE=/tmp/tmp.curl.$$



if [ "$METHOD" == "batch" ] ; then

    cmd="curl"

    grep http $URIFILE | {
    while read host post  
    do
    if [ "$cmd" != "curl" ] ; then 
        cmd="$cmd --next"
    fi 
    cmd="$cmd $host -d $post -X POST"
    echo $cmd
    done
    $cmd
    } 
else
    grep http $URIFILE | {
    while read host post  
    do
    curl $host -d "$post" -X POST
    done
    }
fi

