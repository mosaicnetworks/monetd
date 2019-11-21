#!/bin/bash

# Check for solc
command -v solc >/dev/null 2>&1 || { echo >&2 "I require solc but it's not installed.  Aborting."; exit 1; }

# Check for json_pp
command -v json_pp >/dev/null 2>&1 || { echo >&2 "I require json_pp but it's not installed.  Aborting."; exit 1; }


GIT="$(git rev-parse --abbrev-ref HEAD) $(git show --oneline -s)"
SOLCVERSION=$(solc --version | tail -1)
OSVERSION=$(lsb_release -a 2> /dev/null | sed -e 's/$/;  /g'| paste -s)

echo "Found solc $SOLCVERSION"

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


solc --bin-runtime --abi -o . --overwrite --optimize $mydir/poa.sol
ret=$?

if [ $ret -ne 0 ] ; then
    echo >&2 "solc returned an error. Aborting."
    exit 2
fi 

BYTECODE="$(sed -r 's/(.{72})/   "\1" +\n/g' POA_Genesis.bin-runtime | sed '$s/    "//;$s/" +\n//;$s/^/    "/;$s/$/"/')"

ABI="$(cat POA_Genesis.abi  | json_pp | sed -e 's/"/\\"/g;s/   /\\t/g' | sed -E ':a;N;$!ba;s/\r{0,1}\n/\\n" +\n    "/g'   )" 

cat <<! > $mydir/bytecode.go
package genesis

// This code is generated externally as part of each release.
// The exact toolset is specified to allow the compilation environment to be recreated to
// verify this build

const (
    // SolcCompilerVersion is the solc version used to compile this bytecode
    SolcCompilerVersion = "$SOLCVERSION"
    
    // SolcOSVersion is the output of lsb_release -a for the OS used to compile this bytecode
    SolcOSVersion = "$OSVERSION"

    GitVersion = "$GIT"

	// StandardPOAContractByteCode is the bytecode for the standard POA contract precompiled
	StandardPOAContractByteCode = "" +
$BYTECODE 
    
  	// StandardPOAContractABI is the ABI for the standard POA contract precompiled
	StandardPOAContractABI = "$ABI"  

)

!

