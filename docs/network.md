# Monetcli network

----

## Table of Contents

+ [Parameters](#parameters)
+ [Subcommands](#subcommands)
+ [MonetCli Configuration folder structure](#monetcli-configuration-folder-structure)

----


The network subcommand deals with **network.toml**, a new file that defines a network. It can be used to generate the datadir files - although it contains no private keys. All network commands can take a flag overriding the default directory - but we anticipate it being little used. 

The top level help for the subcommand looks like this:

```bash
$ monetcli help network
Network
		
The network subcommand is used to configure a network of hubs within the monetcli configuration. The compile option builds the genesis file and pushes it to a monetd configuration. The commands available from the network command are sequenced in the wizard, testnet and testjoin commands.

Usage:
  monetcli network [command]

Available Commands:
  add         add key pair
  check       check configuration
  compile     compile configuration
  contract    set solidity contract
  generate    generate and add key pair
  location    show the location of the configuration files
  new         generate new configuration
  params      update parameters interactively
  peers       review peers list
  show        show configuration

Flags:
  -c, --config-dir string   the directory containing the network.toml file holding the monetcli configuration (default "/home/jon/.monetcli")
  -h, --help                help for network

Global Flags:
  -q, --hide-banners   hide banners
  -v, --verbose        verbose messages

Use "monetcli network [command] --help" for more information about a command.
```

### Parameters

The *monetcli network* command supports some globals parameter flags that can be appled to any of the sub commands:
+ **\-c**, **\-\-config-dir [file]** specifies the directory containing the network.toml file holding the monetcli configuration (default "$HOME/.monetcli" on Linux). N.B. the name of the actual **network.toml** file cannot be changed. 
+ **\-q**, **\-\-hide-banners** hides the banners that appear in some of the interactive commands.
+ **\-v**, **\-\-verbose** turns on extra logging output with more details of the operations being performed. 


### Subcommands

**monetcli network add [moniker] [pubkey] [ip] [isValidator]** takes a given key and adds them to the validator set in the monetcli configuration. N.B. the key is not written to a peers.json file until you invoke *monetcli network compile*.


**monetcli network check** checks whether the network.toml file defines a valid configuration and all of the required files are present and readable. **In development:** If the network.toml includes bytecode and solcs version information, it attempt to compile the smart contract and verify the result matches the supplied version.

[comment]: # (//TODO remove In development flag when it is no longer in development)

**monetcli network compile** takes a monetcli network configuration and generates an actual monet hub configuration. It implicitly runs a network check command. It populates a datadir directory including copying any keys stored within the network configuration folder. If the nodename is specified the configuration for that node is written. It is intended that the node name would allow multiple configurations be generated from the same machine - likely useful for node. The POA contract is compiled to build the genesis block. If there is no bytecode in the network.toml it is added with solcs version. Otherwise the bytecode is validated. **N.B.** this command requires an internet connection to run, unless you have run *monetcli network contract*. The default contract is downloaded directly from github. 

[comment]: # (//TODO verify the checking code is put live)


**monetcli network contract [contract]** changes the contract template to use for POA in Monet Hub. The smart contract must comply with the requirements outline in [this document](smartcontract.md). If this command is not used, then a default contract sourced directly from github is used. 

[comment]: # (//TODO - actually write that smartcontract.md file)

**monetcli network generate key [ip] [nodename]** generates a new key and adds them as peers / validators. The private keys are placed in a keystore subfolder. This command is equivalent to running *monetcli key generate* followed by *monetcli network add*

**monetcli network location** outputs the location of the monetcli configuration file. 

**monetcli network new** creates a new template network.toml. If the target already exists this command will exit without changes unless you specify the \-\-force option. If you do specify the force option, then it will rename the existing configuration with a .~n~ suffix, where n is the lowest integer where the folder does not already exist. 
 

**monetcli network params** is an interactive command to set the monetcli parameters that are pushed to monet hub configuration files. These options are:

+ **Logging level** controls which messages are written to the logs. Select from the list, they are sorted from outputting the most messages to the fewest. 
+ **eth.listen** controls where EVM-Lite listens. The default :8080 will normally be fine here.
+ **eth.cache** is the size of the EVM-Lite cache 
+ **babble.listen** IP:PORT of Babble node, which must exactly match this node's entry in peers.json
+ **babble.service-listen** IP:PORT of Babble HTTP API service
+ **babble.heartbeat** Heartbeat time milliseconds (time between gossips)
+ **babble.timeout** TCP timeout milliseconds
+ **babble.cache-size** Number of items in LRU caches
+ **babble.sync-limit** Max number of Events per sync
+ **babble.fast-sync** Enable FastSync
+ **babble.max-pool** Max number of pool connections
+ **babble.bootstrap** Bootstrap from Babble database

You have the opportunity to save or discard your changes at the end of the parameter list. Parameters which have mandatory values in the Monet Hub are not available from this sub command.


**monetcli network peers** provides an interactive interface for managing peers. You are initially show a list of all peers. You can select one from that list to view its complete details. From there you may edit or delete the peer. Delete does not touch any assoicated keys - it just removes the credentials from the list that is compiled into the peers.json file. Edit Peers allows you to edit / amend each of the stored fields for that node. **N.B.** you will need to use *network add* or *network generate* to add a peer.  


**monetcli network show** outputs the current *network.toml* file to screen. 




### MonetCli Configuration folder structure

The network configuration folder will look like:
```
.
├── genesis.sol
├── keystore
│   └── key.json
└── network.toml
```

