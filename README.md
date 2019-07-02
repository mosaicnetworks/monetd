# Monet Hub

**WIP**

This repo is a thin wrapper around [evm-lite](https://github.com/mosaicnetworks/evm-lite) and [babble](https://github.com/mosaicnetworks/babble) which exposes only the 
functionality required to operate a node on the [MONET Hub](https://monet.network/about.html). 

The executables defined in this repository are: 

- [monetd](docs/monetd.md) - the actual Monet Hub node daemon.
- [monetcli](docs/monetcli.md) - a tool to manage configuration for monetd. 

There is an available [quick start document](docs/README.md).

## Overview

The Monet Hub is an always-on blockchain that supports other Mobile ad-hoc 
blockchains as defined in the MONET whitepaper. It is a smart-contract platform
based on the Ethereum Virtual Machine and a BFT consensus algorithm called
Babble (inspired by Hashgraph).

Nodes on the Monet Hub run the `monetd` daemon, which is a specific instance of 
[mosaicnetwork](https://mosaicnetworks.io)'s 
[evm-lite](https://github.com/mosaicnetworks/evm-lite)  with 
[Babble consensus](https://github.com/mosaicnetworks/babble).

## Installation

First install dependencies. We use glide to manage dependencies:

```bash
[...]/monetd$ curl https://glide.sh/get | sh
[...]/monetd$ make vendor
```

This will download all dependencies and put them in the **vendor** folder; it
could take a few minutes.

Then build and install:

```bash
[...]/monetd$ make install
```

## Usage

Please see the documentation for each separate program:

- [monetd](docs/monetd.md) - the actual Monet Hub node daemon.
- [monetcli](docs/monetcli.md) - useful tools including key management. 

## Developer

Further notes are available [here](docs/developer.md) 