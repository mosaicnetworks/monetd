# monetd
Monet Hub Daemon

**This is a work in progress and not yet complete**

This repository is a single repository to build a Monet Hub. It builds uses a EVM_lite VM running on a Babble consensus. No other consensus is supported. 

Currently we are using the poa branch of EVM-Lite and and the develop branch of Babble.

It is intended that the invocation for a validator would look something like:

```bash
monetd run
```

## Repository Roadmap

- ~~Initial Repo Creation, building a single executable~~
- Review of the configuration options available. Most command line options will be removed.
	- We need to categorise the remaining items into core Monet properties, and tuning options. Validators 		should be able to tune, but not alter a fundamental property.mk
	- Whilst the review is in progress it would be sensible to categorise the tuning parameters into those than can be changed whilst running and those that can't. 
- Keygen functionality recreated and available elsewhere - either via the Client or monetcli (or both)
- UI for validators to manage the parameters that they are allowed to set only. We could use the same approach to create a UI for Hub configuration - but it is probably not worth the effort.
- Modify version to show EVM_Lite, Babble and Monet versions. All 3 are critical. Could also show GETH version. Any others?



## Monetd Command Line

This is as per currently. More of these options will be removed. 

```bash
$ monetd help
Monet-Daemon

Usage:
  monetd [command]

Available Commands:
  help        Help about any command
  run         Run a Monet node
  version     Show version info

Flags:
  -h, --help   help for monetd

Use "monetd [command] --help" for more information about a command.

```



```bash
$ monetd help run
Run a Monet node

Usage:
  monetd run [flags]

Flags:
      --babble.cache-size int          Number of items in LRU caches (default 50000)
      --babble.datadir string          Directory contaning priv_key.pem and peers.json files (default "/home/jon/.evm-lite/babble")
      --babble.enable-fast-sync        Enable FastSync
      --babble.heartbeat duration      Heartbeat time milliseconds (time between gossips) (default 500ms)
      --babble.listen string           IP:PORT of Babble node (default ":1337")
      --babble.max-pool int            Max number of pool connections (default 2)
      --babble.service-listen string   IP:PORT of Babble HTTP API service (default ":8000")
      --babble.store                   use persistent store
      --babble.sync-limit int          Max number of Events per sync (default 1000)
      --babble.timeout duration        TCP timeout milliseconds (default 1s)
  -d, --datadir string                 Top-level directory for configuration and data (default "/home/jon/.evm-lite")
      --eth.cache int                  Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string                  Eth database file (default "/home/jon/.evm-lite/eth/chaindata")
      --eth.genesis string             Location of genesis file (default "/home/jon/.evm-lite/eth/genesis.json")
      --eth.keystore string            Location of Ethereum account keys (default "/home/jon/.evm-lite/eth/keystore")
      --eth.listen string              Address of HTTP API service (default ":8080")
      --eth.pwd string                 Password file to unlock accounts (default "/home/jon/.evm-lite/eth/pwd.txt")
  -h, --help                           help for run
      --log string                     debug, info, warn, error, fatal, panic (default "debug")

```





## Changes from EVM-Lite

The following EVM-Lite commands have been deprecated from monetd:

```bash
$ evml help
EVM-Lite

Usage:
  evml [command]

Available Commands:
  help        Help about any command
  keys        An Ethereum key manager
  run         Run a node
  version     Show version info

Flags:
  -h, --help   help for evml

```

- **keys** is deprecated in monetd.

```bash
$ evml help run
Run a node

Usage:
  evml run [command]

Available Commands:
  babble      Run the evm-lite node with Babble consensus
  raft        Run the evm-lite node with Raft consensus
  solo        Run the evm-lite node with Solo consensus (no consensus)

Flags:
  -d, --datadir string        Top-level directory for configuration and data (default "/home/jon/.evm-lite")
      --eth.cache int         Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string         Eth database file (default "/home/jon/.evm-lite/eth/chaindata")
      --eth.genesis string    Location of genesis file (default "/home/jon/.evm-lite/eth/genesis.json")
      --eth.keystore string   Location of Ethereum account keys (default "/home/jon/.evm-lite/eth/keystore")
      --eth.listen string     Address of HTTP API service (default ":8080")
      --eth.pwd string        Password file to unlock accounts (default "/home/jon/.evm-lite/eth/pwd.txt")
  -h, --help                  help for run
      --log string            debug, info, warn, error, fatal, panic (default "debug")

Use "evml run [command] --help" for more information about a command.
```

- **raft** is deprecated in monetd. 
- **solo** is deprecated in monetd. 
- **babble** is no longer an option. Monetd *only* uses babble as a consensus engine.

All the flags bar --datadir, --help and --log are deprecated in monetd. All those options are set in the configuration file. 

```bash
$ evml help run babble
Run the evm-lite node with Babble consensus

Usage:
  evml run babble [flags]

Flags:
      --babble.cache-size int          Number of items in LRU caches (default 50000)
      --babble.datadir string          Directory contaning priv_key.pem and peers.json files (default "/home/jon/.evm-lite/babble")
      --babble.enable-fast-sync        Enable FastSync
      --babble.heartbeat duration      Heartbeat time milliseconds (time between gossips) (default 500ms)
      --babble.listen string           IP:PORT of Babble node (default ":1337")
      --babble.max-pool int            Max number of pool connections (default 2)
      --babble.service-listen string   IP:PORT of Babble HTTP API service (default ":8000")
      --babble.store                   use persistent store
      --babble.sync-limit int          Max number of Events per sync (default 1000)
      --babble.timeout duration        TCP timeout milliseconds (default 1s)
  -h, --help                           help for babble

Global Flags:
... [as per evml run above]

```

All flags bar --babble.datadir are deprecated. 





