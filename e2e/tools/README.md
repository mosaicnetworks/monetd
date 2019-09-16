# Testing Tools

## Sample Docker Session

Start the docker network, with the logging overwritten to set verbose on. 

```bash
$ cd e2e
$ make start TEST=transfer_3_10  VERBOSE=verbose
```
Run the test suite. 

```bash
$ cd tools
$ ./build-trans.sh
```

## build-trans.sh Parameters

Currently there are no CLI parameters, instead the section below is edited. 

```bash
# CLI Params section. These will become parameters
VERBOSE="-v"                 # EIther "" or "-v"
ACCTCNT=3                    # Number of Accounts to transfer between       
TRANSCNT=12                  # Total number of transactions 
FAUCET="Faucet"              # Faucet Account Moniker
PREFIX="Test"                # Prefix of the Moniker for transfer monikers   
NODENAME="Node0"             # Node Name
NODEHOST="172.77.5.10"       # Node IP
NODEPORT="8080"              # Node Port
CONFIGDIR="$HOME/.monettest" # TemporaryMonet Config Dir
OUTDIRSTEM="/tmp"            # Output Directory
```

XXX - I had to run this as well to copy the Faucet key in the right place

```
cp ~/.giverny/networks/bulktransfers/keystore/Faucet* ~/.monet/keystore
```

The Transactions are generated in a ``Trans.$$`` subdirectory of $OUTDIRSTEM. 
It defaults to ``/tmp``.