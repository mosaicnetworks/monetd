.. _monetd_rst:

Monetd
======


Ethereum with Babble consensus
------------------------------

To build the Monet Hub solution, we took the
`Go-Ethereum <https://github.com/ethereum/go-ethereum>`__ implementation (Geth) 
and extracted the EVM and Trie components to create a lean and modular version 
with interchangeable consensus.

The EVM is a virtual machine specifically designed to run untrusted code on a 
network of computers. Every transaction applied to the EVM modifies the State 
which is persisted in a Merkle Patricia tree. This data structure allows to 
simply check if a given transaction was actually applied to the VM and can 
reduce the entire State to a single hash (merkle root) rather analogous to a 
fingerprint.

The EVM is meant to be used in conjunction with a system that broadcasts 
transactions across network participants and ensures that everyone executes the
same transactions in the same order. Ethereum uses a Blockchain and a Proof of
Work consensus algorithm. EVM-Lite makes it easy to use any consensus system,
including `Babble <https://github.com/mosaicnetworks/babble>`__.

We have taken EVM-Lite and Babble and surrounded them with a thin wrapper to
encapsulate the Monet Hub solution in a single repository.

Usage
-----

All the configuration required to run a node is stored under a directory with a
very specific structure. By default, ``monetd`` will look for this directory in
``$HOME/.monet`` (on UNIX systems), but it is possible to override this with the
``--datadir`` flag. You would not normally need to access these configuration
files directly. The ``monetd config`` tool provides CLI interfaces to set up a
Monet network.

``datadir`` must contain a set of files defining the network that this node is
attempting to join or create. Please refer to ``monetd config`` for a tool to
manage this configuration.

In particular:

-  **eth/genesis.json**: defines the accounts prefunded in the state, and the 
   POA smart-contract.
-  **babble/genesis.peers.json**: defines Babble's initial peer-set.
-  **babble/peers.json**: defines Babble's current peer-set
-  **babble/priv\_key**: contains the validator's private key for Babble.

The command ``monetd config location --expanded`` provides further details of
the filepaths used for your instance. Further options pertaining to the 
operation of the node are read from the [datadir]/monetd.toml file, or 
overwritten by the following flags. It is envisaged that you would not need to 
use these flags in a production environment.

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

Configuration
-------------

The application writes data and reads configuration from the directory specified
by the --datadir flag. The directory structure must respect the following
stucture:

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


The Ethereum genesis file defines Ethereum accounts and is stripped of all the
Ethereum POW stuff. This file is useful to predefine a set of accounts that own
all the initial Ether at the inception of the network. N.B. in a Monet 
environment you would need to edit the ``genesis.json`` file after compiling the
POA smart contract into a ``genesis.json`` file. Monet hubs have a ``poa``
section in addition to the ``alloc`` section in the ``genesis.json`` file.

Example Ethereum genesis.json defining two account:

.. code:: json

    {
       "alloc": {
            "629007eb99ff5c3539ada8a5800847eacfc25727": {
                "balance": "1337000000000000000000"
            },
            "e32e14de8b81d8d3aedacb1868619c74a68feab0": {
                "balance": "1337000000000000000000"
            }
       }
    }

API
---

The Service exposes an API at the address specified by the --api-listen flag for
clients to interact with the node and the network.

Get any account
~~~~~~~~~~~~~~~

This method retrieves the information about any account.

.. code:: bash

    host:~$ curl http://[api_addr]/account/0x629007eb99ff5c3539ada8a5800847eacfc25727 -s | json_pp
    {
        "address":"0x629007eb99ff5c3539ada8a5800847eacfc25727",
        "balance":1337000000000000000000,
        "nonce":0
    }

Send raw signed transactions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This endpoint allows sending NON-READONLY transactions ALREADY SIGNED. The
client is left to compose a transaction, sign it and RLP encode it. The
resulting bytes, represented as a Hex string, are passed to this method to be
forwarded to the EVM.

This is an ASYNCHRONOUS operation and the effect on the State should be verified
by fetching the transaction' receipt.

example:

.. code:: bash

    host:~$ curl -X POST http://[api_addr]/rawtx -d '0xf8628080830f424094564686380e267d1572ee409368e1d42081562a8e8201f48026a022b4f68bfbd4f4c309524ebdbf4bac858e0ad65fd06108c934b45a6da88b92f7a046433c388997fd7b02eb7128f4d2401ef2d10d574c42edf15875a43ee51a1993' -s | json_pp
    {
        "txHash":"0x5496489c606d74ad7435568393fa2c4619e64497267f80864109277631aa849d"
    }

Get Transaction receipt
~~~~~~~~~~~~~~~~~~~~~~~

example:

.. code:: bash

    host:~$ curl http://[api_addr]/tx/0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf -s | json_pp
    {
       "to" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
       "root" : "0xc8f90911c9280651a0cd84116826d31773e902e48cb9a15b7bb1e7a6abc850c5",
       "gasUsed" : "0x5208",
       "from" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
       "transactionHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf",
       "logs" : [],
       "cumulativeGasUsed" : "0x5208",
       "contractAddress" : null,
       "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
    }

Get consensus info
------------------

The ``/info`` endpoint exposes a map of information provided by the consensus
system.

example:

.. code:: bash

    host:-$ curl http://[api_addr]/info | json_pp
    {
       "rounds_per_second" : "0.00",
       "type" : "babble",
       "consensus_transactions" : "10",
       "num_peers" : "4",
       "consensus_events" : "10",
       "sync_rate" : "1.00",
       "transaction_pool" : "0",
       "state" : "Babbling",
       "events_per_second" : "0.00",
       "undetermined_events" : "22",
       "id" : "1785923847",
       "last_consensus_round" : "1",
       "last_block_index" : "0",
       "round_events" : "0"
    }

Client
------

Please refer to `EVM-Lite
CLI <https://github.com/mosaicnetworks/evm-lite-cli>`__ for Javascript
utilities and a CLI to interact with the API.
