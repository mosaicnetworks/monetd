#!/bin/bash


URIFILE=$1


grep http $URIFILE | {
while read host post  
do

   echo  curl $host -d "$post" -X POST
done
}