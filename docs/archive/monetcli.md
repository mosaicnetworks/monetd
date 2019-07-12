# Monet-CLI

----

## Table of Contents

+ [Monet Hub tools](#monet-hub-tools)
	+ [help](#help)
	+ [version](#version)
	+ [keys](#keys)
	+ [network](#network)
	+ [config](#config)
	+ [testnet](#testnet)
	+ [testjoin](#testjoin)
	+ [wizard](#wizard)
+ [Examples](#examples)

----


# Monet Hub tools

`monetcli` provides a suite of tools for configuring and managing a Monet Hub.

There are currently 8 subcommands (listed here in the order they are in this document). 

+ [help](#help) offers more help on individual commands. 
+ [version](#version) shows a version number
+ [keys](#keys) is an Ethereum key manager with commands for generating new keys, inspecting existing keys and changing passphrases. 
+ [network](#network) is a set of subcommands for building and configuring a new monet network of hubs. These tools would normally be used before a network is started
+ [config](#config) is a set of subcommands for managing the .monet configuration files of an extant network. 
+ [testnet](#testnet) is a menu driven wizard for building a test network, enabling peers to create and share a common configuration, via a server, before starting the network
+ [testjoin](#testjoin) is a menu driven wizard for joining an existing network
+ [wizard](#wizard) is a menu driven front end to the network subcommands to guide the creation of a new network of hubs.

## help

You can get more information about the commands available by using the help command or the `--help` flag. Thus you can get top level help: 
```
$ monetcli help
Monet-CLI
	
Monetcli is the swiss army knife of tools for the Monet Hub. The README.md in the github repository is a good starting place in the documentation. For quicker access and a handy reference of flags and options.:
	
monetcli help [subcommand]
	
The best starting points are the commands testnet, testjoin or wizard.

Usage:
  monetcli [command]

Available Commands:
  config      manage monetd configuration
  help        Help about any command
  keys        an Ethereum key manager
  network     manage monet network configuration
  testjoin    join monetd testnets
  testnet     manage monetd testnets
  version     show version info
  wizard      wizard to set up a Monet Network

Flags:
  -h, --help           help for monetcli
  -q, --hide-banners   hide banners
  -v, --verbose        verbose messages

Use "monetcli [command] --help" for more information about a command.
```

Or you can get details of a subcommand:

```bash
$ monetcli help keys
an Ethereum key manager

Usage:
  monetcli keys [command]

Available Commands:
  generate    generate a new keyfile
  inspect     inspect a keyfile
  update      change the passphrase on a keyfile

Flags:
  -h, --help              help for keys
      --json              output JSON instead of human-readable format
      --passfile string   the file that contains the passphrase for the keyfile

Global Flags:
  -v, --verbose        verbose messages

Use "monetcli keys [command] --help" for more information about a command.
```


## version

This command returns the version number for the `monetcli`, and the version of the EVM-Lite, Babble and Geth librarys used to build it. The suffix (if shown) on the Monet version is the github commit for this version.

```bash
$ monetcli version
Monet Version: 0.0.2-0a516594
     EVM-Lite Version: 0.2.0
     Babble Version: 0.4.2
     Geth Version: 1.8.27

```



## keys
The keys subcommand is used to manage ethereum keys.

```bash
$ monetcli help keys
Keys
	
An Ethereum key manager.

Usage:
  monetcli keys [command]

Available Commands:
  generate    generate a new keyfile
  inspect     inspect a keyfile
  update      change the passphrase on a keyfile

Flags:
  -h, --help              help for keys
      --json              output JSON instead of human-readable format
      --passfile string   the file that contains the passphrase for the keyfile

Global Flags:
  -q, --hide-banners   hide banners
  -v, --verbose        verbose messages

Use "monetcli keys [command] --help" for more information about a command.
```

The generate command generates a new key pair. You either need to use the `--passfile option` or enter a pass phrase when prompted by the application.
```bash
$ monetcli keys generate keys.json
Passphrase: 
Address: 0x83434e68b52Ef809538224BF78472cc3F6a17bcC
$ cat keys.json
{"address":"83434e68b52ef809538224bf78472cc3f6a17bcc","crypto":{"cipher":"aes-128-ctr","ciphertext":"878c888d14cd407af2f99061f432cef2c08232b4a2f99f80d9240b9ac5cb6b24","cipherparams":{"iv":"c2ac23f51d5d79fb45ead639fa7f9d7f"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"064cac9be036d0eae1c24ebc0073e02ad773289a16c7c19235dc567d957d08df"},"mac":"1f34dbe8d834911fe5048e7c183eb0608d75719a6a989c99d243bc09fb292bb3"},"id":"54d565f2-2fe1-4ee7-af5d-8619dc6bdcce","version":3}
```

You can inspect key files using the inspect command:
```bash
$ monetcli keys inspect keys.json
Passphrase: 
Address:        0x83434e68b52Ef809538224BF78472cc3F6a17bcC
Public key:     043b463098401fe38241174a9bf28e6b1d64b2b1f7061c2d4d4a2a8a73a8e389c53547bb99fb5f93579b31ca5aeb975e3d1f4577fbf05b0698a11deb720e2670c0
```

You can change the passphrase for the key with the update command:

```bash
$ monetcli keys update keys.json
Passphrase: 
Please provide a new passphrase
Passphrase: 
Repeat passphrase: 
```

## network

The network subcommand deals with **network.toml**, a new file that defines a network. There is a separate [document](network.md) for this command.

## config

The config subcommand deals with the actual `monetd` configuration datadir. 

```bash
$ monetcli config help
monetcli config

The config subcommands manage the monet configuration file, as used by 
the monetd server process.

Usage:
  monetcli config [command]

Available Commands:
  check       check the configuration
  location    show the location of the configuration files
  publish     publish a monet node configuration from the monetcli configuration
  show        show the configuration files

Flags:
  -c, --config-dir string         the directory containing monet nodes configurations (default "/home/jon/.monetcli")
  -h, --help                      help for config
  -m, --monet-config-dir string   the directory containing monet nodes configurations (default "/home/jon/.monet")

Global Flags:
  -q, --hide-banners   hide banners
  -v, --verbose        verbose messages

Use "monetcli config [command] --help" for more information about a command.
```


**monet config check** sanity checks the datadir configuration. 

**monet config location** displays the location of a monetd configuration. 

**monet config publish** publishes a monetcli configuration to a monetd configuration. There are 2 modes of publishing a configuration. By default the configuration files for all nodes are published, with each in a separate subdirectory. If you specify \-\-node-name only the configuration for that node is published.   

**monet config show** echoes the monet.toml configuration of a monetd configuration to screen. 



## testnet

**testnet** is a menu driven wizard for building a test network, enabling peers to create and share a common configuration, via a server, before starting the network

This command is documented in [Testnet Docs](testnet.md) 

```
$ monetcli help testnet
TestNet Config

This subcommand facilitates the process of building a testnet 
on a set of separate servers (i.e. not Docker). It is a menu driven 
guided process. Take a look at docs/testnet.md for fuller instructions.

Usage:
  monetcli testnet [flags]

Flags:
  -h, --help      help for testnet
  -p, --publish   jump straight to polling for a configuration

Global Flags:
  -v, --verbose   verbose messages
 ``` 


## testjoin

Test join is a command to allow the menu driven configuration for joining an existing network. There are no options as the command is interactive.

Invoke it thus:
```
$ monetcli testjoin
✔ Existing peer:  : the.existing.peer
```

Once you have specified a peer, it is queried for a genesis file and the two peers files. 
```
Downloading files from  the.existing.peer
Downloaded  /home/jon/.monetcli/testnet/peers.genesis.json
Downloaded  /home/jon/.monetcli/testnet/peers.json
Downloaded  /home/jon/.monetcli/testnet/genesis.json
```

You then enter a passphrase for the key pair that you are about to generate. 
```
Enter Keystore Password:   : #|
Confirm Keystore Password:   : #|
Address: 0x9B39Af7F8C599e67379Ec429d41A0B71Dc21F24e
Building Data to push to Configuration Server
Pub Key  :  046a0dc579184801c1ab4144f93005af0f73778d2bad5f755bd98ad499934e6c6869c34cd8252ff79cadf1b829ecb328bb03717593c558be7b0c6040543944393d
Address  :  0x9B39Af7F8C599e67379Ec429d41A0B71Dc21F24e
```

You then enter your IP address. This is used by Babble as part of its join request. 
```
Enter your ip without the port:   : |192.168.1.18
```

There is a final confirmation as the overwritten of .monet is a destructive operation. 
```
Use the arrow keys to navigate: ↓ ↑ → ← 
? Confirm Overwriting Existing Configuration  : 
  ▸ No
    Yes
```

Files are copied from .monetcli/testnet to the appropriate folders under .monetd. NB you .evmlc config is also amended with connection details to the new network.  
```
Renaming /home/jon/.monet to /home/jon/.monet.~11~
Copying to  0 /home/jon/.monet/monetd.toml
Copying to  1 /home/jon/.monet/eth/genesis.json
Copying to  2 /home/jon/.monet/babble/peers.json
Copying to  3 /home/jon/.monet/babble/priv_key
Copying to  4 /home/jon/.monet/babble/peers.genesis.json
Copying to  5 /home/jon/.monet/eth/pwd.txt
Copying to  6 /home/jon/.monet/eth/keystore/keyfile.json
Copying to  7 /home/jon/.monet/keyfile.json
Updating evmlc config
Try running:  monetd run
```
N.B. at this point you are not authorised. You will need to pass the join.json details to an existing validator. They will nominate your node, and the existing validators need to vote on your nomination. The person who nominated you will inform you when (and if) you are approved and can thus start your node successfully. 


## wizard
**wizard** is a menu driven front end to the network subcommands to guide the creation of a new network of hubs.


```bash
$ monetcli help wizard
Wizard to set up a Monet Network

This command provides a wizard interface to the "monetcli network" commands. This provides a guided interface through the process of configuring a network.

Usage:
  monetcli wizard [flags]

Flags:
  -h, --help   help for wizard

Global Flags:
  -v, --verbose   verbose messages
```

The wizard command is described in a [separate document](wizard.md).

----


# Examples


First we create a new network.
```bash
$ monetcli network new
``` 

If you get a message that the configuration directory exists, then you need to add the **force** parameter. No data is lost, the existing directory is renamed with a .~n~ where n is the lowest non-clashing positive integer.


```bash
$ monetcli network new --force 
``` 

All of the network commands have a **verbose** option to display more information.  

Next we generate some keys. Here we specify the moniker, ip address and whether the node is a validator on the commands line. 

```bash
$ monetcli network generate node0 172.77.5.1 true
Passphrase: 
Address: 0x7e42c360141DA6d5B80109eF3101f3A210BbaA32
```

At any point we can view the configuration so far. 

```bash
$ monetcli network show

[config]
  datadir = "/home/jon/.monetcli"

[poa]
  compilerversion = ""
  contractaddress = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
  contractname = "genesis_array.sol"

[validators]
  addresses = "0x7e42c360141DA6d5B80109eF3101f3A210BbaA32"
  ips = "172.77.5.1"
  isvalidator = "true"
  monikers = "node0"
``` 

By default the POA smart contract is downloaded from github directly if not previously specified. You may overwrite this default by using the **contract** subcommand. 

```bash
$ monetcli network contract ../evm-lite/e2e/smart-contracts/genesis_array.sol 
```

You can add a peer with existing keys as follows:
```bash
$ monetcli network add node1 1bbabaababbabaababbabaababbabaababbabaab 192.168.0.1 true --verbose
```
There is inbuilt validation of the configuration settings that are run before compiling the network configuration. This can also be invoked directly:

```bash
$ monetcli network check
All checks passed
```

When you are satisfied with the configuration the actual config files for the node can be built. 

```bash
$ monetcli network compile
```
The compile option, takes the specified contract if provided, otherwise it downloads a contract from github, and inserts the initial peer set into the smart contract. This contract is then compiled and inserted into a generated genesis.json file. 

----

<sup>[Documents Index](README.md) | [GitHub repo](https://github.com/mosaicnetworks/monetd) | [Monet](https://monet.network/) | [Mosaic Networks](https://www.babble.io/)</sup>