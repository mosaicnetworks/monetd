.. _monetd_configuration_rst:

Monetd Configuration
====================

All the configuration required to run a node is stored under a directory with a
very specific structure. By default, ``monetd`` will look for this directory in
``$HOME/.monet`` (on Linux), but it is possible to override this with the
``--datadir`` flag.

The directory must respect the following stucture:

::

    host:~/.monet$ tree
    ├── babble
    │   ├── peers.genesis.json
    │   ├── peers.json
    │   └── priv_key
    ├── eth
    │   ├── genesis.json
    │   └── poa
    │       ├── compile.toml
    │       ├── contract0.abi
    │       └── contract0.sol
    ├── keystore
    │   ├── node0.json
    ├── monetd.toml


You would not normally need to access these configuration files directly. The
``monetd config`` tool provides a CLI interfaces to set up a Monet network. The
command ``monetd config location --expanded`` provides further details of the
filepaths used for your instance.

Eth
---

The eth/genesis.json file defines prefunded accounts in the state, as well as
the POA smart-contract. This file is useful to predefine a set of accounts that
own all the initial tokens at the inception of the network. In addition, the
``poa`` section contains information about the POA smart-contract.

Example genesis.json defining one prefunded account (the ABI and bytecode of
the smart-contract have been truncated):

.. code:: json

    {
        "alloc": {
                "a10aae5609643848ff1bceb76172652261db1d6c": {
                        "balance": "1234567890000000000000",
                        "moniker": "node0"
                }
        },
        "poa": {
                "address": "0xaBBAABbaaBbAABbaABbAABbAABbaAbbaaBbaaBBa",
                "abi": "[\n\t{\n\t\t\"constant\": true, ... ]",
                "code": "6080604052600436106101095760003560e01c8063..."
                }
    }

Babble
------

-  **babble/genesis.peers.json**: defines Babble's initial peer-set.

-  **babble/peers.json**: defines Babble's current peer-set

-  **babble/priv\_key**: contains the validator's private key for Babble.

Run Options
-----------

Options pertaining to the operation of the node are read from the
[datadir]/monetd.toml file, or overwritten by the following flags. It is
envisaged that you would not need to use these flags in a production
environment.

::

  Flags:
        --api-listen string           IP:PORT of Monet HTTP API service (default ":8080")
        --babble.bootstrap            Bootstrap Babble from database
        --babble.cache-size int       Number of items in LRU caches (default 50000)
        --babble.heartbeat duration   Heartbeat time milliseconds (time between gossips) (default 500ms)
        --babble.listen string        IP:PORT of Babble node (default "192.168.1.3:1337")
        --babble.max-pool int         Max number of pool connections (default 2)
        --babble.sync-limit int       Max number of Events per sync (default 1000)
        --babble.timeout duration     TCP timeout milliseconds (default 1s)
        --eth.cache int               Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
    -h, --help                        help for run

  Global Flags:
    -d, --datadir string   Top-level directory for configuration and data (default "/home/user/.monet")
        --log string       trace, debug, info, warn, error, fatal, panic (default "debug")
    -v, --verbose          verbose messages

Example of a monet.toml file:

::

  datadir = "/home/user/.monet"
  log = "debug"
  api-listen = ":8080"

  [babble]
    listen = "192.168.1.3:1337"
    heartbeat = "500ms"
    timeout = "1s"
    cache-size = 50000
    sync-limit = 1000
    max-pool = 2
    bootstrap = false

  [eth]
    cache = 128
