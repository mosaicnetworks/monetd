#!/bin/bash


STEM=${1:-"172.77.5."}
LO=${2:-10}
HI=${3:-13}


lastlbi="potato"


for i in  $(seq $LO $HI)
do
    lbi=$(wget -O - -q http://${STEM}${i}:8080/info | jq '.last_block_index' | sed -e 's/"//g' )

    if [ "$lastlbi" == "potato" ] ; then
        lastlbi="$lbi"
        continue
    fi

    if [ "$lbi" != "$lastlbi" ] ; then
        echo "($i) Mismatched last block index $lbi $lastlbi"
        lastlbi="$lbi"
    fi
done

echo "Last Block Index is: $lbi"

for bi in $(seq 1 $lbi)
do
  lastfh="potato"  
  for i in $(seq $LO $HI)
  do
    wq=$(wget -O - -q http://${STEM}${i}:8080/block/$bi )

    fh=$(echo $wq | jq '.Body.FrameHash')
    sh=$(echo $wq | jq '.Body.StateHash')
    if [ "$lastfh" == "potato" ] ; then
        lastfh="$fh"
        lastsh="$sh"
        continue
    fi

    if [ "$fh" != "$lastfh" ] ; then
        if [ "$fh" \> "$lastfh" ] ; then
            echo "Mismatched Frame Hash block $(printf "%04d" $bi):  $lastfh $fh"
        else        
            echo "Mismatched Frame Hash block $(printf "%04d" $bi):  $fh $lastfh"    
        fi    
        lastfh="$fh"
    fi

    if [ "$sh" != "$lastsh" ] ; then
        if [ "$sh" \> "$lastsh" ] ; then
            echo "Mismatched State Hash block $(printf "%04d" $bi):  $lastsh $sh"
        else        
            echo "Mismatched State Hash block $(printf "%04d" $bi):  $sh $lastsh"    
        fi    

        lastsh="$sh"
    fi


  done
done

 # 172.77.5.10:8080/block/13