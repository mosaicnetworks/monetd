# monetd
Monet Hub Daemon

**This is a work in progress and not yet complete**


monetd run



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

