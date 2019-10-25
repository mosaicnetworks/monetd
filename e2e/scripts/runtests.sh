#!/bin/bash

TEST=${1:-""}
NOSTOP=${2:-""}

NODUMMY="true"
RANSPECIFIEDTEST="false"

STARTPARAMS=""


set -eu

IP=172.77.5.10
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
declare -A testresults

for d in "$mydir"/../tests/*
do
    if [ ! -d "$d" ] ; then
        continue
    fi

    testname=$(basename $d)
    output="$mydir"/../tests/"$testname.out"


# If set we only run the matching test
    if [ "$TEST" != "" ] ; then
        if [ "$TEST" != "$testname" ] ; then
            continue
        fi   
        RANSPECIFIEDTEST="true" 
    fi


    if [ "$NODUMMY" == "true" ] && [ "$testname" == "dummypass" ] ; then
            continue
    fi

    if [ "$NODUMMY" == "true" ] && [ "$testname" == "dummyfail" ] ; then
            continue
    fi
    

    echo -n "$testname ."
   

    if [ -f "$mydir/../tests/$testname/start.params" ] ; then
        STARTPARAMS=$(cat "$mydir/../tests/$testname/start.params")
  #      echo $STARTPARAMS
    fi

    $mydir/start.sh --network=$testname --init-ip=$IP  $STARTPARAMS  2>&1 > $output

    echo -n "."
    set +e
    $mydir/../tests/$testname/run-test.sh  2>&1 >> $output
    testresults[$testname]=$?

    echo -n ". "${testresults[$testname]}
    set -e

    echo "."
   
    if [ "$NOSTOP" = "nostop" ] ; then
        continue
    fi
       
    $mydir/stop.sh $testname 2>&1 >> $output

done 


    if [ "$TEST" != "" ] ; then
        if [ "$RANSPECIFIEDTEST" != "true" ] ; then
            echo "Test $TEST not found. Aborting."
            exit 9            
        fi   
    fi


echo ""
echo ""
echo "Results"
echo ""

exitcode=0
for k in  "${!testresults[@]}"
do
    if [ ${testresults[$k]} -eq 0 ] ; then
        echo "[PASS] : $k"
    else      
	    echo "[FAIL] : $k (${testresults[$k]})"
        exitcode=1
    fi
done

if [ $exitcode -eq 1 ]; then
    echo ""
    echo "Logs are in $mydir/../tests/"
fi
    
exit $exitcode