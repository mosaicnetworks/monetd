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

The MONET Hub is a blockchain that supports other Mobile ad-hoc blockchains, as 
defined in the [MONET whitepaper](http://bit.ly/monet-whitepaper). It is a 
smart-contract platform based on the Ethereum Virtual Machine and a BFT 
Consensus algorithm. It leverages 
[evm-lite](https://github.com/mosaicnetworks/evm-lite) and 
[Babble](https://github.com/mosaicnetworks/babble).

More information about MONET can be found on the [website](https://monet.network/about.html).

This repository contains the code for `monetd` and `monetcli`:

- `monetd` is the server process that validators are expected to run.

- `monetcli` is a tool that helps with configuration to join or create a
  network.

For the impatient, there is a 
[quick start document](https://monetd.readthedocs.io/en/latest/README.html#quick-start), 
and [installation docs](https://monetd.readthedocs.io/en/latest/install.html).

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
                |  |       |  Babble Consensus    |      |  |
                |  |       +----------------------+      |  |
                |  |                                     |  |
                |  +-------------------------------------+  |
                |                                           |
                +-------------------------------------------+

```

## Installation

Please see the [installation documentation](https://monetd.readthedocs.io/en/latest/install.html).

## Usage

Please see the documentation for each separate program, althought the [quick start document](https://monetd.readthedocs.io/en/latest/README.html) may be the best starting point:

- [monetd](https://monetd.readthedocs.io/en/latest/monetd.html) - the server process.
- [monetcli](https://monetd.readthedocs.io/en/latest/monetcli.html) - useful tools including key management. 
