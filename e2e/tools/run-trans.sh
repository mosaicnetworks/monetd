#!/bin/bash


URIFILE=$1
# METHOD="batch"
METHOD="notbatch"
# TMPFILE=/tmp/tmp.curl.$$

OUTFILE=$URIFILE.out
ERRFILE=$URIFILE.err
TMPFILE=$URIFILE.tmp

> $OUTFILE
> $ERRFILE

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
        
    # //TODO Remove
    #     host="35.176.237.170:8080/rawtx"

         json=$( curl -v $host -d "$post" -X POST  2> $TMPFILE); 

         rc="$(grep -v "POST is already" $TMPFILE)"       
        

   

         if [ "$(echo $rc | grep -c "HTTP/1.1 200 OK")" -gt 0 ] ; then
            echo -n "."
         else   
          { echo "$host $post" 
          echo "$rc" 
          echo "$json\n\n"            
          echo ""; } >> $ERRFILE   

          echo ""
          echo "ERROR: $URIFILE"
          echo "$host $post" 
          echo ""
          echo "$rc"           
          echo ""
          echo "$json"
          echo ""               
         fi   

         { echo "$URIFILE $host $post" 
         echo "$rc" 
         echo "$json" 
         echo "" ; } >> $OUTFILE
    done
    }
fi


if [ ! -s $ERRFILE ] ; then
    rm  $ERRFILE
fi

if [ -f "$TMPFILE" ] ; then
    rm $TMPFILE 
fi    