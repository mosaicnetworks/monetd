# Monet Hub

```
                                      /\ \__    /\ \    
  ___ ___      ___     ___       __   \ \ ,_\   \_\ \   
/' __` __`\   / __`\ /' _ `\   /'__`\  \ \ \/   /'_` \  
/\ \/\ \/\ \ /\ \L\ \/\ \/\ \ /\  __/   \ \ \_ /\ \L\ \ 
\ \_\ \_\ \_\\ \____/\ \_\ \_\\ \____\   \ \__\\ \___,_\
 \/_/\/_/\/_/ \/___/  \/_/\/_/ \/____/    \/__/ \/__,_ /
```

![Monet Logo](docs/assets/monet_logo.png)

**//TODO** Feature graphic - probably the `monet` hub one.

----

## Table of Contents

+ [Overview](#overview)
    + [Architecture](#architecture)
+ [Installation](#installation)
+ [Usage](#usage)

----

**//TODO** Sales pitch. Giac? Maybe move the Overview section. 

This repo is a thin wrapper around [evm-lite](https://github.com/mosaicnetworks/evm-lite) and [babble](https://github.com/mosaicnetworks/babble) which exposes only the functionality required to operate a node on the [MONET Hub](https://monet.network/about.html), whilst enforcing required shared parameters such as enabling POA. 

For the impatient, there is a [quick start document](docs/README.md), and [installation docs](docs/install.md).

The executables defined in this repository are: 

- [monetd](docs/monetd.md) - the actual Monet Hub node daemon.
- [monetcli](docs/monetcli.md) - a tool to manage configuration for `monetd`. 

## Overview

The Monet Hub is an always-on blockchain that supports other Mobile ad-hoc 
blockchains as defined in the MONET whitepaper. It is a smart-contract platform
based on the Ethereum Virtual Machine and a BFT consensus algorithm called
Babble (inspired by Hashgraph).

Nodes on the Monet Hub run the `monetd` daemon, which is a specific instance of 
[mosaicnetwork](https://mosaicnetworks.io)'s 
[evm-lite](https://github.com/mosaicnetworks/evm-lite)  with 
[Babble consensus](https://github.com/mosaicnetworks/babble) running with Proof Of Authority.

### Architecture

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
