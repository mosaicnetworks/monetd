# monetd


----

## Table of Contents

+ [Ethereum with Babble consensus](#ethereum-with-babble-consensus)
+ [USAGE](#usage)
+ [Configuration](#configuration)
+ [API](#api)
   + [Get controlled accounts](#get-controlled-accounts)
   + [Get any account](#get-any-account)
   + [Send transactions from controlled accounts](#send-transactions-from-controlled-accounts)
   + [Get Transaction receipt](#get-transaction-receipt)
   + [Send raw signed transactions](#send-raw-signed-transactions)
+ [Get consensus info](#get-consensus-info)
+ [CLIENT](#client)

----


## Ethereum with Babble consensus

We took the [Go-Ethereum](https://github.com/ethereum/go-ethereum)
implementation (Geth) and extracted the EVM and Trie components to create a
lean and modular version with interchangeable consensus.

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
including [Babble](https://github.com/mosaicnetworks/babble).

## USAGE

All the configuration required to run a node is stored under a directory with a 
very specific structure. By default, `monetd` will look for this directory in 
`$HOME/.monet` (on UNIX systems), but it is possible to override this with the 
`--datadir` flag.

`datadir` must contain a set of files defining the network that this node is 
attempting to join or create. Please refer to `monetcli` for a tool to manage 
this configuration. 

In particular:

* **eth/genesis.json**: defines the accounts prefunded in the state, and the POA 
                        smart-contract.
* **babble/genesis.peers.json**: defines Babble's initial peer-set.
* **babble/peers.json**: defines Babble's current peer-set
* **babble/priv_key**: contains the validator's private key for Babble.

Further options pertaining to the operation of the node are read from the 
[datadir]/monetd.toml file, or overwritten by the following flags.

```
Flags:
      --babble.cache-size int          Number of items in LRU caches (default 50000)
      --babble.heartbeat duration      Heartbeat time milliseconds (time between gossips) (default 500ms)
      --babble.listen string           IP:PORT of Babble node (default ":1337")
      --babble.max-pool int            Max number of pool connections (default 2)
      --babble.service-listen string   IP:PORT of Babble HTTP API service (default ":8000")
      --babble.sync-limit int          Max number of Events per sync (default 1000)
      --babble.timeout duration        TCP timeout milliseconds (default 1s)
  -d, --datadir string                 Top-level directory for configuration and data (default "/home/martin/.monet")
      --eth.cache int                  Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.listen string              IP:PORT of Monet HTTP API service (default ":8080")
  -h, --help                           help for run
```

Example of a monet.toml file:

```
[babble]
heartbeat = "50ms"
timeout = "200ms"
listen = ":1337"
sync-limit = 500

[eth]
listen = ":8080"
```

## Configuration

The application writes data and reads configuration from the directory specified
by the --datadir flag. The directory structure must respect the following
stucture:

```
host:~/.monetd$ tree
├── babble
│   ├── peers.json
│   └── priv_key
├── eth
│   ├── genesis.json
└── monetd.toml
```

The Ethereum genesis file defines Ethereum accounts and is stripped of all the 
Ethereum POW stuff. This file is useful to predefine a set of accounts that own 
all the initial Ether at the inception of the network.

Example Ethereum genesis.json defining two account:
```json
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
```  

## API
The Service exposes an API at the address specified by the --eth.listen flag for
clients to interact with the node and the network.  

### Get controlled accounts

This endpoint returns all the accounts that are controlled by the evm-lite
instance. These are the accounts whose private keys are present in the keystore.

example:
```bash
host:~$ curl http://[api_addr]/accounts -s | json_pp
{
   "accounts" : [
      {
         "address" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
         "balance" : 1337000000000000000000,
         "nonce": 0
      },
      {
         "address" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
         "balance" : 1337000000000000000000,
         "nonce": 0
      }
   ]
}
```
### Get any account

This method retrieves the information about any account, not just the ones whose 
keys are included in the keystore.  

```bash
host:~$ curl http://[api_addr]/account/0x629007eb99ff5c3539ada8a5800847eacfc25727 -s | json_pp
{
    "address":"0x629007eb99ff5c3539ada8a5800847eacfc25727",
    "balance":1337000000000000000000,
    "nonce":0
}
```

### Send transactions from controlled accounts

Send a transaction from an account controlled by the evm-lite instance. The
transaction will be signed by the service since the corresponding private key is
present in the keystore.

example: Send Ether between accounts  
```bash
host:~$ curl -X POST http://[api_addr]/tx -d '{"from":"0x629007eb99ff5c3539ada8a5800847eacfc25727","to":"0xe32e14de8b81d8d3aedacb1868619c74a68feab0","value":6666}' -s | json_pp
{
   "txHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf"
}
```

### Get Transaction receipt

example:
```bash
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

```

### Send raw signed transactions

Most of the time, one will require to send transactions from accounts that are
not  controlled by the evm-lite instance. The transaction will be assembled,
signed  and encoded on the client side. The resulting raw signed transaction
bytes can be submitted to evm-lite through the `/rawtx` endpoint.  

example:
```bash
host:~$ curl -X POST http://[api_addr]/rawtx -d '0xf8628080830f424094564686380e267d1572ee409368e1d42081562a8e8201f48026a022b4f68bfbd4f4c309524ebdbf4bac858e0ad65fd06108c934b45a6da88b92f7a046433c388997fd7b02eb7128f4d2401ef2d10d574c42edf15875a43ee51a1993' -s | json_pp
{
    "txHash":"0x5496489c606d74ad7435568393fa2c4619e64497267f80864109277631aa849d"
}
```

## Get consensus info

The `/info` endpoint exposes a map of information provided by the consensus
system.

example:
```bash
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

```

## CLIENT

Please refer to [EVM-Lite CLI](https://github.com/mosaicnetworks/evm-lite-cli)
for Javascript utilities and a CLI to interact with the API.

----

<sup>[Documents Index](README.md) | [GitHub repo](https://github.com/mosaicnetworks/monetd) | [Monet](https://monet.network/) | [Mosaic Networks](https://www.babble.io/)</sup>