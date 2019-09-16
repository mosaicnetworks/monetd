# Testing Tools

## Sample Docker Session

Start the docker network, with the logging overwritten to set verbose on. 

```bash
$ cd e2e
$ make start TEST=transfer_03_10  VERBOSE=verbose
```
Run the test suite. 

```bash
$ cd tools
$ ./build-trans.sh
```

## build-trans.sh Parameters

This is now parameterised. 

```bash
$ ./build-trans.sh --help
./build-trans.sh [-v] [--accounts=10] [--transactions=200] [--faucet="Faucet"] [--faucet-config-dir=] [--prefix=Test] [--node-name=Node] [--node-host=172.77.5.11] [--node-port=8080] [--config-dir=/home/jon/.monettest] [--temp-dir=/tmp] [-h|--help]
```

The parameters above map to the variables below:

```bash
# CLI Params default section
VERBOSE="-v"                    # Either "" or "-v"
ACCTCNT=10                      # Number of Accounts to transfer between       
TRANSCNT=200                    # Total number of transactions 
FAUCET="Faucet"                 # Faucet Account Moniker
FAUCETCONFIG=""                 # Keystore to copy faucet key from
PREFIX="Test"                   # Prefix of the Moniker for transfer monikers   
NODENAME="Node"                 # Node Name
NODEHOST="172.77.5.11"          # Node IP
NODEPORT="8080"                 # Node Port
CONFIGDIR="$HOME/.monettest"    # Monet Config Dir used for this test
OUTDIRSTEM="/tmp"               # Output Directory
```

Besides ``--accounts`` to set the number of accounts to transfer tokens between
and ``--transactions`` to set the the number of transactions, you are likely to
want to set ``--faucet`` to the moniker of an account that you have the a key
pair got. It is likely that this was created in another system. Set
``--faucet-config-dir`` to either ``$HOME/.monet/keystore`` or 
``$HOME/.giverny/networks/<network name>/keystore`` to avoid having to copy
the keys manually. 


The Transactions are generated in a ``trans.$$`` subdirectory of $OUTDIRSTEM. 
It defaults to ``/tmp``.

At the moment, the ``trans.$$`` folder is not deleted.