.. _monet_api_rst:

Monetd API
==========

``monetd`` exposes an HTTP API at the address specified by the ``--api-listen``
flag. This document contains the API specification with some basic examples
using curl. For API clients (javascript libraries, CLI, and GUI), please refer
to `Monet CLI <https://github.com/mosaicnetworks/monetcli>`__

Get Account
-----------

Retrieve information about any account.

.. code:: http

  GET /account/{address}
  returns: JSONAccount

.. code:: go

    type JSONAccount struct {
      Address string   `json:"address"`
      Balance *big.Int `json:"balance"`
      Nonce   uint64   `json:"nonce"`
      Code    string   `json:"bytecode"`
    }

Example:

.. code:: bash

    host:~$ curl http://localhost:8080/account/0xa10aae5609643848fF1Bceb76172652261dB1d6c -s | jq
    {
        "address": "0xa10aae5609643848fF1Bceb76172652261dB1d6c",
        "balance": 1234567890000000000000,
        "nonce": 0,
        "bytecode": ""
    }

Call
----

Call a smart-contract READONLY function. These calls will NOT modify the EVM
state, and the data does NOT need to be signed.

.. code:: http

  POST /call
  data: JSON SendTxArgs
  returns: JSON JsonCallRes

.. code:: go

    type SendTxArgs struct {
        From     common.Address  `json:"from"`
        To       *common.Address `json:"to"`
        Gas      uint64          `json:"gas"`
        GasPrice *big.Int        `json:"gasPrice"`
        Value    *big.Int        `json:"value"`
        Data     string          `json:"data"`
        Nonce    *uint64         `json:"nonce"`
    }

    type JSONCallRes struct {
        Data string `json:"data"`
    }

Example:

.. code:: bash

    curl http://localhost:8080/call \
    -d '{"constant":true,"to":"0xabbaabbaabbaabbaabbaabbaabbaabbaabbaabba","value":0,"data":"0x8f82b8c4","gas":1000000,"gasPrice":0,"chainId":1}' \
    -H "Content-Type: application/json" \
    -X POST -s | jq
    {
      "data": "0x0000000000000000000000000000000000000000000000000000000000000001"
    }

Submit Transaction
------------------

Send a SIGNED, NON-READONLY transaction. The client is left to compose a
transaction, sign it and RLP encode it. The resulting bytes, represented as a
Hex string, are passed to this method to be forwarded to the EVM. This is a
SYNCHRONOUS operation; it waits for the transaction to go through consensus and
returns the transaction receipt.

.. code:: http

  POST /rawtx
  data: STRING Hex representation of the raw transaction bytes
  returns: JSON JSONReceipt

.. code:: go

    type JSONReceipt struct {
        Root              common.Hash     `json:"root"`
        TransactionHash   common.Hash     `json:"transactionHash"`
        From              common.Address  `json:"from"`
        To                *common.Address `json:"to"`
        GasUsed           uint64          `json:"gasUsed"`
        CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
        ContractAddress   common.Address  `json:"contractAddress"`
        Logs              []*ethTypes.Log `json:"logs"`
        LogsBloom         ethTypes.Bloom  `json:"logsBloom"`
        Status            uint64          `json:"status"`
    }
    
Example:

.. code:: bash

    host:~$ curl -X POST http://localhost:8080/rawtx 0xf86904808398968094f7cd2ba6892341e568e9d825c4bdc2bd53b7524189031b9d1340ad2500008026a04eb7420aa52a1955d26ffb16d3a8cb8d969ae0eb6d75bb5076599c42a788e08da0178b3ddb264cdcc624121f55a95ae45de119bc44a0a85b721d8958b7ebe0553a -s | json_pp
    {
      "root": "0xda4529d2bc5e8b438edee4463637eb91d5490edb50d15e786e8d5276f2a2c8f4",
      "transactionHash": "0x3f5682786828d26946e12a08a858b6dd805d1ea8f7d39d93f1d4d5393b23f710",
      "from": "0x888980abf63d4133482e50bf8233f307e3c2b941",
      "to": "0xf7cd2ba6892341e568e9d825c4bdc2bd53b75241",
      "gasUsed": 21000,
      "cumulativeGasUsed": 21000,
      "contractAddress": "0x0000000000000000000000000000000000000000",
      "logs": [],
      "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "status": 1
    }


Get Receipt
-----------

Get a transaction receipt. When a transaction is applied to the EVM, a receipt
is saved to record if/how the transaction affected the state. This contains
such information as the address of a newly created contract, how much gas was
use, and the EVM Logs produced by the execution of the transaction.

.. code:: http

  GET /tx/{tx_hash}
  returns: JSON JSONReceipt

Example:

.. code:: bash

    host:~$ curl http://localhost:8080/tx/0x96764078446cfbaec6265f173fb5a2411b7c272052640bca622252494a74dbb4 -s | jq
    {
      "root": "0x348c230578e27e20a10924e925f7cddb28279561b52cab7b31750c6d4716ac21",
      "transactionHash": "0x96764078446cfbaec6265f173fb5a2411b7c272052640bca622252494a74dbb4",
      "from": "0xa10aae5609643848ff1bceb76172652261db1d6c",
      "to": "0xa10aae5609643848ff1bceb76172652261db1d6c",
      "gasUsed": 21000,
      "cumulativeGasUsed": 21000,
      "contractAddress": "0x0000000000000000000000000000000000000000",
      "logs": [],
      "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "status": 0
    }

Info
----

Get information about a Babble instance.

.. code:: http

  GET /info
  returns: JSON map

Example:

.. code:: bash

    host:-$ curl http://localhost:8080/info | jq
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

POA
---

Get details of the PoA smart-contract.

.. code:: http

  GET /poa
  returns: JSONContract

.. code:: go

    type JSONContract struct {
        Address common.Address `json:"address"`
        ABI     string         `json:"abi"`
    }

Example (trunctated output):

.. code:: bash

    host:-$ curl http://localhost:8080/poa | jq
    {
        "address": "0xabbaabbaabbaabbaabbaabbaabbaabbaabbaabba",
        "abi": "[\n\t{\n\t\t\"constant\": true,\n\t\t\"inputs\"...]"
    }


Genesis.json
------------

This endpoint returns the content of the genesis.json file in JSON format.
This allows new nodes to pull the genesis file from an existing peer.

.. code:: http

  GET /genesis
  returns: JSON Genesis

.. code:: go

    type Genesis struct {
        Alloc AccountMap
        Poa   PoaMap
    }

    type AccountMap map[string]struct {
        Code        string
        Storage     map[string]string
        Balance     string
        Authorising bool
    }

    type PoaMap struct {
        Address string
        Balance string
        Abi     string
        Code    string
    }

Example (truncated output):

.. code:: bash

  host:-$ curl://http://locahost:8080/genesis | jq
  {
    "Alloc": {
      "a10aae5609643848ff1bceb76172652261db1d6c": {
        "Code": "",
        "Storage": null,
        "Balance": "1234567890000000000000",
        "Authorising": false
      }
    },
    "Poa": {
      "Address": "0xaBBAABbaaBbAABbaABbAABbAABbaAbbaaBbaaBBa",
      "Balance": "",
      "Abi": "[\n\t{\n\t\t\"constant\": ...]",
      "Code": "6080604052600436106101095..."
    }
  }

Block
-----

Get a Babble Block by index.

.. code:: http

  GET /block/{index}
  returns: JSON Block

.. code:: go

    type Block struct {
        Body       BlockBody
        Signatures map[string]string
    }

    type BlockBody struct {
        Index                       int
        RoundReceived               int
        StateHash                   []byte
        FrameHash                   []byte
        PeersHash                   []byte
        Transactions                [][]byte
        InternalTransactions        []InternalTransaction
        InternalTransactionReceipts []InternalTransactionReceipt
    }

Example:

.. code:: bash

    host:-$ curl http://locahost:8080/block/0 | jq
    {
      "Body": {
        "Index": 0,
        "RoundReceived": 1,
        "StateHash": "VY6jFi7P5bIajdWvwZU2jU0q3KXDcp1sFx7Ye6kl1/k=",
        "FrameHash": "Nek4dF0ybGZQ1XEuJQrjmPtNrfPLAtGU4sTQSSB2iKM=",
        "PeersHash": "Gv+YqIq56l6LZWdhAsx0XEB4gjZluMaziv7hCXT5b9k=",
        "Transactions": [
          "+GSAgIMPQkCUq7qruqu6q7qruqu6q7qruqu6q7qAhOHHOSoloGCfTsLEOcMMXDX1W/78zpaZTXXK8BSR1Q8cCqicSrExoDv/0YGlpaGMJO8B6ZAJ/WAiEOKG00uzF8piaCvW3GHH"
        ],
        "InternalTransactions": [],
        "InternalTransactionReceipts": []
      },
      "Signatures": {
        "0X04F91D4429AE73229141F960B70CD2E83BF39D6EBF1B951C4E65BA9F0EE7FA2365C859CC9BF856709F78F0B9DD6BFBA450BFC7B8123197616D22E6EA8693201800": "2gtf6rkdc0q29n1isef0x2fib64qlf075uybtva6558r8onv31|2gnym6xat1ok68nqtsymcpg4x9ihj1ouwab8inode5m8eb82tb"
      }
    }

Current Peers
-------------

Get Babble's current peer-set.

.. code:: http

    Get /peers
    returns: []Peer

.. code:: go

    type Peer struct {
        NetAddr   string
        PubKeyHex string
        Moniker   string
    }

Example:

.. code:: bash

    $host:-$ curl http://localhost:8080/peers | jq
    [
      {
        "NetAddr": "192.168.1.3:1337",
        "PubKeyHex": "0X04F91D4429AE73229141F960B70CD2E83BF39D6EBF1B951C4E65BA9F0EE7FA2365C859CC9BF856709F78F0B9DD6BFBA450BFC7B8123197616D22E6EA8693201800",
        "Moniker": "node0"
      }
    ]

Genesis Peers
-------------

Get Babble's initial validator-set.

.. code:: http

    GET /genesispeers
    returns: []Peer

History
-------

Get the entire history of validator-sets. It returns a map of validator-sets
indexed by round number. The round number corresponds to the Hashgraph round at
which the corresponding validator-set took effect.

.. code:: http

    Get /history
    returns: map[int][]Peer

.. code:: go

    type Peer struct {
        NetAddr   string
        PubKeyHex string
        Moniker   string
    }

Example:

.. code:: bash

    $host:-$ curl http://localhost:8080/history | jq

    "0": [
      {
        "NetAddr": "node0.monet.network:1337",
        "PubKeyHex": "0X04717E964AB361D46268389691B2B9638A07079263097E7F685397935A0038569EA9BCF10AC5FAF3025312014192D4303FE5981CC3CB8AE3C7C0645E1510D9D4BB",
        "Moniker": "node0"
      }
    ],
    "7": [
      {
        "NetAddr": "node0.monet.network:1337",
        "PubKeyHex": "0X04717E964AB361D46268389691B2B9638A07079263097E7F685397935A0038569EA9BCF10AC5FAF3025312014192D4303FE5981CC3CB8AE3C7C0645E1510D9D4BB",
        "Moniker": "node0"
      },
      {
        "NetAddr": "node1.monet.network:1337",
        "PubKeyHex": "0X04824C5EA5E0169ECA63141F388563D14F05E0C71BC9B5AA2A10D44EC06931CC56E4C5AE18473941AC92739AE9E7B9E792DE29E9368920F56FD559ABF7139430A7",
        "Moniker": "node1"
      }
    ],
    "21": [
      {
        "NetAddr": "node0.monet.network:1337",
        "PubKeyHex": "0X04717E964AB361D46268389691B2B9638A07079263097E7F685397935A0038569EA9BCF10AC5FAF3025312014192D4303FE5981CC3CB8AE3C7C0645E1510D9D4BB",
        "Moniker": "node0"
      }
    ],
    "31": [
      {
        "NetAddr": "node0.monet.network:1337",
        "PubKeyHex": "0X04717E964AB361D46268389691B2B9638A07079263097E7F685397935A0038569EA9BCF10AC5FAF3025312014192D4303FE5981CC3CB8AE3C7C0645E1510D9D4BB",
        "Moniker": "node0"
      },
      {
        "NetAddr": "node1.monet.network:1337",
        "PubKeyHex": "0X04824C5EA5E0169ECA63141F388563D14F05E0C71BC9B5AA2A10D44EC06931CC56E4C5AE18473941AC92739AE9E7B9E792DE29E9368920F56FD559ABF7139430A7",
        "Moniker": "node1"
      }
    ],
    "44": [
      {
        "NetAddr": "node0.monet.network:1337",
        "PubKeyHex": "0X04717E964AB361D46268389691B2B9638A07079263097E7F685397935A0038569EA9BCF10AC5FAF3025312014192D4303FE5981CC3CB8AE3C7C0645E1510D9D4BB",
        "Moniker": "node0"
      },
      {
        "NetAddr": "node1.monet.network:1337",
        "PubKeyHex": "0X04824C5EA5E0169ECA63141F388563D14F05E0C71BC9B5AA2A10D44EC06931CC56E4C5AE18473941AC92739AE9E7B9E792DE29E9368920F56FD559ABF7139430A7",
        "Moniker": "node1"
      },
      {
        "NetAddr": "node2.monet.network:1337",
        "PubKeyHex": "0X0475ADBA5AD67D0A12EBD68FA696806D411012E8C1A8D82A6EA08403E9C03F16D9F81A4CB175015ABBC47E45CA6E30999454A4D2CD040040A6687BA729D96913EC",
        "Moniker": "node2"
      }
    ],
