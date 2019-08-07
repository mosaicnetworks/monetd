.. _monet_api_rst:

Monetd API
==========

Monetd exposes an HTTP API at the address specified by the ``--api-listen``
flag. This document contains the API specification with some basic examples
using curl. For API clients (javascrip libraries, CLI, and GUI), please refer
to `EVM-Lite CLI <https://github.com/mosaicnetworks/evm-lite-cli>`__

Get Account
-----------

Retrieve information about any account.

.. code::

  GET /account/{address}
  returns: JsonAccount

.. code:: go

    type JsonAccount struct {
      Address string   `json:"address"`
      Balance *big.Int `json:"balance"`
      Nonce   uint64   `json:"nonce"`
      Code    string   `json:"bytecode"`
    }

Example:

.. code::

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

.. code::

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

    type JsonCallRes struct {
        Data string `json:"data"`
    }

Example:

.. code::

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
Hex string, are passed to this method to be forwarded to the EVM. This is an
ASYNCHRONOUS operation and the effect on the State should be verified by
fetching the transaction's receipt.

.. code::

  POST /rawtx
  data: STRING Hex representation of the raw transaction bytes
  returns: JSON JsonTxRes

.. code:: go

    type JsonTxRes struct {
        TxHash string `json:"txHash"`
    }

Example:

.. code:: bash

    host:~$ curl -X POST http://localhost:8080/rawtx -d '0xf8600180830f424094a10aae5609643848ff1bceb76172652261db1d6c648026a03c14b99e14420e34c15885ff3afc1043aa6e4e13e2be4691d74a772cde44819ba0652b202ab93908544ea4d7d89567fa462fa719f381e54aa6781ba96c2e9e0e90' -s | json_pp
    {
        "txHash":"0x96764078446cfbaec6265f173fb5a2411b7c272052640bca622252494a74dbb4"
    }

Get Receipt
-----------

Get a transaction receipt. When a transaction is applied to the EVM, a receipt
is saved to record if/how the transaction affected the state. This contains
such information as the address of a newly created contract, how much gas was
use, and the EVM Logs produced by the execution of the transaction.

.. code::

  GET /tx/{tx_hash}
  returns: JSON JsonReceipt

.. code:: go

    type JsonReceipt struct {
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

Get information about Babble.

.. code::

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

.. code::

  GET /poa
  returns: JsonContract

.. code:: go

    type JsonContract struct {
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

This endpoint returns the content of the genesis.json file.

.. code::

  GET /genesis
  returns: JSON Genesis

.. code::

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

.. code::

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

.. code::

  GET /block/{index}
  returns: JSON Block

.. code::

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

.. code::

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

.. code::

    Get /peers
    returns: []Peer

.. code::

    type Peer struct {
        NetAddr   string
        PubKeyHex string
        Moniker   string
    }

Example:

.. code::

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

.. code::

    GET /genesispeers
    returns: []Peer