# Monet Hub

[![Documentation Status](https://readthedocs.org/projects/monetd/badge/?version=latest)](https://monetd.readthedocs.io/en/latest/?badge=latest)
[![Go Report](https://goreportcard.com/badge/github.com/mosaicnetworks/monetd)](https://goreportcard.com/report/github.com/mosaicnetworks/monetd)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

----

## This repo is currently in the process of being published
## It may be unstable during this process

----
----

![Monet Logo](docs/assets/monet_logo.png)

## Overview

As defined in the [MONET whitepaper](http://bit.ly/monet-whitepaper), the Hub is 
a blockchain that supports other Mobile ad-hoc blockchains with services like 
peer-discovery, Inter-Blockchain Communication, and the Tenom token. It is a 
smart-contract platform based on the Ethereum Virtual Machine and a BFT 
Consensus algorithm. It leverages 
[evm-lite](https://github.com/mosaicnetworks/evm-lite) and 
[Babble](https://github.com/mosaicnetworks/babble).

More information about MONET can be found on the [website](https://monet.network/about.html).

This repository contains the code for `monetd` and `monetcli`:

- `monetd` is the server process that validators are expected to run.

- `monetcli` is a tool that helps with configuration to join or create a 
  network.

## Documentation

* [design](https://monetd.readthedocs.io/en/latest/design.html)
* [installation](https://monetd.readthedocs.io/en/latest/install.html)
* [monetd](https://monetd.readthedocs.io/en/latest/monetd.html)
* [monetcli](https://monetd.readthedocs.io/en/latest/monetcli.html)
* [clients](https://monetd.readthedocs.io/en/latest/clients.html)