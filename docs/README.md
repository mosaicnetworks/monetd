# The Monet Hub

The monetd respository contains the tools necessary to run and maintain a validator hub in a Monet network. 

They naturally divide into 2 sections:
+ [MonetCLI](monetcli.md) -- the swiss army knife of utilities
+ [Testnet Docs](monetd.md) -- the hub server process

Full details can found at the links above, but the Quick Start section below may help you where to look. 

## Quick Start

The installation process is covered in the [Root Readme](../README.md).

----

The general purpose guided install can be accessed via:
```bash
$ monetcli wizard
```

See the wizard section in [Monet CLI docs](monetcli.md) for more information.  

----
To set up a new testnet with yourself as one of the initial peers use:
```bash
$ monetcli testnet
```

See the testnet section [Monet CLI docs](monetcli.md) for more information.  

----
To join an existing testnet use:
```bash
$ monetcli testjoin
```

See the testjoin section [Monet CLI docs](monetcli.md) for more information.  

