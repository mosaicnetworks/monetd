#!/bin/bash


NODES=$*


lastbi=""

for n in $NODES
do
   url="http://$n:8080/info"
   bi=$(curl -s $url  | sed "s/},/\n/g" | sed 's/.*"last_block_index":"\([^"]*\)".*/\1/'  )

 #  echo "$n: $bi"

   if [ "$lastbi" != "" ] && [ "$lastbi" != "$bi" ] ; then
        echo "last block index mismatch. $bi & $lastbi "
        exit 401
   fi

   lastbi="$bi"
done

# If we reach here all nodes have bi last_block_index



lastsh=""
lastfh=""
lastph=""
exitcode=0

for n in $NODES
do
   url="http://$n:8080/block/$lastbi"
   raw=$(curl -s $url)

   sh=$(echo $raw  | json_pp | grep "StateHash" | sed 's/.*"StateHash"[ \t]*:[ \t]*"\([^"]*\)".*/\1/'  )
   fh=$(echo $raw  | json_pp | grep "FrameHash" | sed 's/.*"FrameHash"[ \t]*:[ \t]*"\([^"]*\)".*/\1/'  )
   ph=$(echo $raw  | json_pp | grep "PeersHash" | sed 's/.*"PeersHash"[ \t]*:[ \t]*"\([^"]*\)".*/\1/'  )


   if [ "$lastsh" != "" ] && [ "$lastsh" != "$sh" ] ; then
        echo "statehash mismatch."
        echo "   $sh"
        echo "   $lastsh"
        exitcode=402
   fi

   lastsh="$sh"

   if [ "$lastfh" != "" ] && [ "$lastfh" != "$fh" ] ; then
        echo "framehash mismatch."
        echo "   $fh"
        echo "   $lastfh"
        exitcode=403
   fi

   lastfh="$fh"


   if [ "$lastph" != "" ] && [ "$lastph" != "$ph" ] ; then
        echo "peers hash mismatch."
        echo "   $ph"
        echo "   $lastph"
        exitcode=404
   fi

   lastph="$ph"
done

exit $exitcode