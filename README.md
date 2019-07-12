# Monet Hub

[![Documentation Status](https://readthedocs.org/projects/monetd/badge/?version=latest)](https://monetd.readthedocs.io/en/latest/?badge=latest)

----

# This repo is currently in the process of being published
# It may be unstable during this process

----

----


![Monet Logo](docs/assets/monet_logo.png)


## Overview

The [MONET Hub](https://monet.network/about.html) is an always-on blockchain that supports other Mobile ad-hoc blockchains as defined in the MONET whitepaper. It is a smart-contract platform based on the Ethereum Virtual Machine and a BFT consensus algorithm called Babble (inspired by Hashgraph).

Nodes on the Monet Hub run the `monetd` daemon, which is a specific instance of [mosaicnetwork](https://mosaicnetworks.io)'s [evm-lite](https://github.com/mosaicnetworks/evm-lite)  with 
[Babble consensus](https://github.com/mosaicnetworks/babble) running with Proof Of Authority.

For the impatient, there is a [quick start document](docs/README.md), and [installation docs](docs/install.md).

----

## Table of Contents

+ [Overview](#overview)
    + [Architecture](#architecture)
+ [Installation](#installation)
+ [Usage](#usage)

----

## Architecture

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State       |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  | -Keystore   |         | -Trie       |  |
                |  |             |         | -Database   |  |
                |  +-------------+         +-------------+  |
                |         |                       ^         |     
                |         v                       |         |
                |  +-------------------------------------+  |
                |  | Engine                              |  |
                |  |                                     |  |
                |  |       +----------------------+      |  |
                |  |       | Consensus            |      |  |
                |  |       +----------------------+      |  |
                |  |                                     |  |
                |  +-------------------------------------+  |
                |                                           |
                +-------------------------------------------+

```




## Installation

Please see the [installation documentation](docs/install.md).

## Usage

Please see the documentation for each separate program, althought the [quick start document](docs/README.md) may be the best starting point:

- [monetd](docs/monetd.md) - the actual Monet Hub node daemon.
- [monetcli](docs/monetcli.md) - useful tools including key management. 
