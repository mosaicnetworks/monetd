.. _overview_rst:

Overview
========

This document describes the tools for operating a Monet Toolchain node, and a
couple of important concepts regarding the account model. In other documents,
we provide guidance on using these tools to perform common tasks, as well as a
complete reference of commands and API functions.

Tools
-----

Monetd
~~~~~~

**monetd** is the server process that connects to other nodes, participates in
the consensus algorithm, and maintains its own copy of the application state.
Additionaly, the **giverny** program facilitates the creation of local Monet
Toolchain networks for testing purposes. We don't expect most people to use
**giverny** as it is mostly a development tool.

**monetd** and **giverny** are written in Go, and reside in the same git
repository because they share significant source code. Please follow the
:ref:`installation instructions<install_rst>` to get started.

Monetcli
~~~~~~~~

**monetcli** is the client-side program that interacts with a running Monet
Toolchain node, and enables users to make transfers, query accounts, deploy and
call smart-contracts, or participate in the PoA governance mechanism.
``monetcli`` is a `Node.js <https://nodejs.org/>`__ project. It can be installed easily with ``npm
install -g monetcli``.

Accounts
--------

What is an account?
~~~~~~~~~~~~~~~~~~~

The Monet Toolchain, and thus MONET, uses the same account model as Ethereum.
Accounts represent identities of external agents and are associated with a
balance (and storage for Contract accounts). They rely on public key
cryptography to sign transactions so that the EVM can securely validate the
identity of a transaction sender.

Using the same account model as Ethereum doesn’t mean that existing Ethereum
accounts automatically have the same balance in MONET (or vice versa). In
Ethereum, balances are denoted in Ether, the cryptocurrency maintained by the
public Ethereum network. On the other hand, every MONET network (even a single
node network) maintains a completely separate ledger and may use any name for
the corresponding coin. The official MONET token is Tenom.

What follows is mostly inspired from the `Ethereum
Docs <http://ethdocs.org/en/latest/account-management.html>`__:

Accounts are objects in the EVM State. They come in two types: Externally owned
accounts, and Contract accounts. Externally owned accounts have a balance, and
Contract accounts have a balance and storage. The EVM State is the state of all
accounts which is updated with every transaction. The underlying consensus
engine ensures that every participant in a Monet Toolchain network processes the same
transactions in the same order, thereby arriving at the same State. The use of
Contract accounts with the EVM makes it possible to deploy and use
*SmartContracts* which we will explore in another document.

What is an account file?
~~~~~~~~~~~~~~~~~~~~~~~~

This is best explained in the `Ethereum
Docs <http://ethdocs.org/en/latest/account-management.html>`__:

   Every account is defined by a pair of keys, a private key, and public key.
   Accounts are indexed by their address which is derived from the public key
   by taking the last 20 bytes. Every private key/address pair is encoded in a
   keyfile. Keyfiles are JSON text files which you can open and view in any
   text editor. The critical component of the keyfile, your account’s private
   key, is always encrypted, and it is encrypted with the password you enter
   when you create the account.

Transactions
------------

A transaction is a signed data package that contains instructions for the EVM.
It can contain instructions to move coins from one account to another, create a
new Contract account, or call an existing Contract account. Transactions are
encoded using the custom Ethereum scheme, RLP, and contain the following
fields:

-  The recipient of the message.
-  A signature identifying the sender and proving their intention to send the
   transaction.
-  The number of coins to transfer from the sender to the recipient.
-  An optional data field, which can contain the message sent to a contract.
-  A STARTGAS value, representing the maximum number of computational steps the
   transaction execution is allowed to take.
-  a GASPRICE value, representing the fee the sender is willing to pay for gas.
   One unit of gas corresponds to the execution of one atomic instruction,
   i.e., a computational step.
