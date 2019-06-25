# Monet-CLI
## Monet Hub tools

## USAGE


```
Monet-CLI

Usage:
  monetcli [command]

Available Commands:
  help        Help about any command
  keys        An Ethereum key manager
  version     Show version info

Flags:
  -h, --help   help for monetcli

Use "monetcli [command] --help" for more information about a command.
```

The keys subcommand is used to manage ethereum keys.

```bash
An Ethereum key manager

Usage:
  monetcli keys [command]

Available Commands:
  generate    Generate a new keyfile
  inspect     Inspect a keyfile
  update      change the passphrase on a keyfile

Flags:
  -h, --help              help for keys
      --json              output JSON instead of human-readable format
      --passfile string   the file that contains the passphrase for the keyfile

Use "monetcli keys [command] --help" for more information about a command.
```

## Configuration

### Network

The network subcommand deals with **network.toml**, a new file that defines a network. It can be used to generate the datadir files - although it contains no private keys. All network commands can take a flag overriding the default directory - but we anticipate it being little used. 

**monetcli network new** is an interactive tool that allows the user to create a template network.toml. They will choose the following items:

- output directory for network.toml. We don't allow the actual file name to be changed. 
- which prebaked template to use (choose from a list)
- choose whether to generate a new key pair and add them as a peer/validator. 
- which smart contract to use - default to the one in the repo

**monetcli network check** checks whether the network.toml file defines a valid configuration. If the network.toml includes bytecode and solcs version information, it attempt to compile the smart contract and verify the result matches the supplied version. 

**monetcli network show** outputs the current *network.toml* file. 

**monetcli network generate key [ip] [nodename]** generates a new key and adds them as peers / validators. The private keys are placed in a keystore subfolder. 

**monetcli network add key [pubkey] [ip] [nodename]** adds a given key and adds them to the validator set.

The network configuration filter will look like:
```
.
├── genesis.sol
├── keystore
│   └── key.json
└── network.toml
```

**monetcli network compile [output-dir] [nodename]** takes a network file and generates an actual monet hub configuration. It implicitly runs a network check command. It populates a datadir directory including copying any keys stored within the network configuration folder. If the nodename is specified the configuration for that node is written. It is intended that the node name would allow multiple configurations be generated from the same machine - likely useful for node. The POA contract is compiled to build the genesis block. If there is no bytecode in the network.toml it is added with solcs version. Otherwise the bytecode is validated. 

**This functionality is currently implemented in bash scripts calling solcs. This may end up not being a go command. **

### Config
The config subcommand deals with the actual monetd configuration datadir. 

**monet config check** sanity checks the datadir configuration. 




