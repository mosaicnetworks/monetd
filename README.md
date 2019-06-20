# Monet Hub
**This is a work in progress and not yet complete**

This repository is a single repository to build a Monet Hub. It builds uses a EVM_lite VM running on a Babble consensus. No other consensus is supported. It also creates tools for the management of the hub. 

The outputs from this repository are 

- [monetd](docs/monetd.md) - the actual Monet Hub Server
- [monetcli](docs/monetcli.md) - useful tools including key management. 

## Installation

First install dependencies. We use glide to manage dependencies:

```bash
[...]/monetd$ curl https://glide.sh/get | sh
[...]/monetd$ make vendor
```
This will download all dependencies and put them in the **vendor** folder; it
could take a few minutes.

Then we build and install:

```bash
[...]/monetd$ make install
```

## Usage

Please see the documentation for each separate function:


- [monetd](docs/monetd.md) - the actual Monet Hub Server
- [monetcli](docs/monetcli.md) - useful tools including key management. 

## Developer

Further notes are available [here](docs/developer.md) 