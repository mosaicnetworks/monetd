# The Monet Hub

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
└── archive                  
    └── developer.md         - out dated developer doc, scheduled to be removed

```