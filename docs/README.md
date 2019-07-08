# The Monet Hub

![Monet Logo](assets/monet_logo.png) 


----

## Table of Contents

+ [Quick Start](#quick-start)
    + [Installation](#installation)
    + [Interactive Configuration](#interactive-configuration)
    + [Creating a new Test Net](#creating-a-new-test-net)
    + [Joining an existing Test Net](#joining-an-existing-test-net)
    + [Clients](#clients)
+ [Contents of the docs folder](#contents-of-the-docs-folder)

----

The monetd respository contains the tools necessary to run and maintain a validator hub in a Monet network. 

They naturally divide into 2 sections:
+ [MonetCLI](monetcli.md) -- the swiss army knife of utilities
+ [Testnet Docs](monetd.md) -- the hub server process

Full details can found at the links above, but the Quick Start section below may help you where to look. 

# Quick Start

## Installation
The installation process is covered in [here](install.md).

----
## Interactive Configuration
The general purpose guided configuration can be accessed via:
```bash
$ monetcli wizard
```

See the wizard section in [Monet CLI docs](monetcli.md) for more information.  

----
## Creating a new Test Net
To set up a new testnet with yourself as one of the initial peers use:
```bash
$ monetcli testnet
```

See the testnet section [Monet CLI docs](monetcli.md) for more information.  

----
## Joining an existing Test Net
To join an existing testnet use:
```bash
$ monetcli testjoin
```

See the testjoin section [Monet CLI docs](monetcli.md) for more information.  



----
## Monet
To join an existing testnet use:
```bash
$ monetcli testjoin
```

See the testjoin section [Monet CLI docs](monetcli.md) for more information.  

----

## Clients

Clients and wallets configured to be used with the monet hub are described [here](clients.md).


----

# Contents of the docs folder

```
├── install.md               - installation instructions
├── monetcli.md              - monetcli command documentation
├── monetd.md                - monetcfg command documentation
├── network.md               - monetcli network command docs, linked from monetcli.md
├── README.md                - this document
├── smartcontract.md         - requirements for poa smart contract for monet
├── testnet.md               - monetcli testnet command docs, linked from monetcli.md
├── wizard.md                - monetcli wizard command docs, linked from monetcli.md
└── archive                  - deprecated docs, scheduled to be removed

```


----

<sup>[Documents Index](README.md) | [GitHub repo](https://github.com/mosaicnetworks/monetd) | [Monet](https://monet.network/) | [Mosaic Networks](https://www.babble.io/)</sup>
