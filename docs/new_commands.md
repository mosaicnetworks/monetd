# New commands

**NB this is a placeholding file. The content of which 
will be converted to rst and integrated with the existing docs. 
It's just more convenient to write them here.**

The following examples are written for Linux / Unix systems. There are minor variations in filepaths for Windows / Mac systems. These will be explicitly addressed in a later revision of these documents. 


# Creating A Single Node Network

The simplest configuration is a single node. 

## Clear Previous Configurations

The following command is not necessary for a clean install, but it is good practice to run 
it anyway. This command renames any previous .monet configuration with a .\~n\~ suffix and creates a new empty configuration folder.  
```bash
monetcli config clear
```

## Create An Account

Firstly we need to create an account 
for the validator. We will generate a private and public keypair. The syntax for the `monetcli keys new` command is:

```bash
monetcli keys new [--moniker moniker] [--passphrase eth.txt]
```
You need to supply a moniker for this node. You will also need to supply a passphrase to secure your key. You can either use the `--passphrase` parameter to specify a file containing the passphrase, or you can specify interactively. 

Thus we need to run:
```bash
$ monetcli keys new --moniker node0 
Enter Passphrase:
Re-Enter Passphrase:
Address: 0x83434e68b52Ef809538224BF78472cc3F6a17bcC
```

The configuration files are written to a `.monetcli` folder. On Linux it is `$HOME/.monetcli`. You can find the location for your instance with this command:

```bash
$ monetcli network location
```

You should have a structure like this at this stage.

```bash
$ tree $HOME/.monetcli
.monetcli
├── accounts
│   └── node0
│       ├── keyfile.json
│       └── peer.toml

```

## Build the Configuration

We now build the monetd configuration files. The syntax for this command is:

```bash
$ monetcli config build [--node node0] [--address ip]
```

We need the IP address for the node you building a network upon. For a live network that would clearly be a public IP address, but for an exploratory test net, we would recommend using an internal IP address. On Linux `ifconfig` will give you IP address information. 

Thus we need to run this command, but replace `192.168.1.4` with the IP address / hostname you discover above.

```bash
$ monetcli config build --node node0  --address 192.168.1.4 
```

This command builds the configuration files for a monetd node. It adds the node given by the `--node` parameter to the initial peer set, including adding it to the initial validator whitelist in the POA smart contract. 

If the `--node` parameter is missing or the value does not correspond to a key pair previously defined, then the list of valid nodes is shown.

If the `--address` parameter is missing, a best guess IP is shown. 

The location of the configuration files depend on the OS. On Linux it is `$HOME/.monet`. You can find the location for your instance with this command:

```bash
$ monetcli config location
The Monet Configuration files are located at:
/home/user/.monet
```

You should have a structure like this at this stage.

```bash
# //TODO After writing the code, insert the relevant tree here. 
$ tree $HOME/.monet
.monet
├── accounts
│   └── node0
│       ├── keyfile.json
│       └── peer.toml

```

## Starting the Node

To start running the monetd node in a terminal window run:

```bash
$ monetd run
```

This is clearly not a production configuration, where you would use `nohup` and redirect log output to the files. 


## Testing

Start EVM-Lite-CLI in interactive mode, and run some commands:

```bash
$ evmlc i
  _____  __     __  __  __           _       _   _               ____   _       ___ 
 | ____| \ \   / / |  \/  |         | |     (_) | |_    ___     / ___| | |     |_ _|
 |  _|    \ \ / /  | |\/| |  _____  | |     | | | __|  / _ \   | |     | |      | | 
 | |___    \ V /   | |  | | |_____| | |___  | | | |_  |  __/   | |___  | |___   | | 
 |_____|    \_/    |_|  |_|         |_____| |_|  \__|  \___|    \____| |_____| |___|
                                                                                    
 Mode:        Interactive
 Data Dir:    /home/jon/.evmlc
 Config File: /home/jon/.evmlc/config.toml
 Keystore:    /home/jon/.evmlc/keystore

  Commands:

    help [command...]    Provides help for a given command.
    exit                 Exits application.
    clear                Clears interactive mode console output
    info [options]       Display information about node
    version [options]    Display current version of cli
    transfer [options]   Initiate a transfer of token(s) to an address

  Command Groups:

    accounts *           5 sub-commands.
    config *             2 sub-commands.
    poa *                6 sub-commands.

evmlc$ info -f
.-------------------------------------.
|          Key           |   Value    |
|------------------------|------------|
| consensus_events       | 0          |
| consensus_transactions | 0          |
| events_per_second      | 0.00       |
| id                     | 1022922485 |
| last_block_index       | -1         |
| last_consensus_round   | nil        |
| moniker                | node0      |
| num_peers              | 1          |
| round_events           | 0          |
| rounds_per_second      | 0.00       |
| state                  | Babbling   |
| sync_rate              | 1.00       |
| transaction_pool       | 0          |
| type                   | babble     |
| undetermined_events    | 0          |
'-------------------------------------'
evmlc$ accounts list -f
.-----------------------------------------------------------------------------.
|                  Address                   |        Balance         | Nonce |
|--------------------------------------------|------------------------|-------|
| 0x46e05762e981d040283af871DcA60A71a6786A23 | 1234000000000000000000 |     0 |
'-----------------------------------------------------------------------------'

evmlc$ accounts get -f 0x46e05762e981d040283af871DcA60A71a6786A23
.-----------------------------------------------------------------------------------------------.
|                  Address                   |            Balance            | Nonce | Bytecode |
|--------------------------------------------|-------------------------------|-------|----------|
| 0x46e05762e981d040283af871DcA60A71a6786A23 | 1,234,000,000,000,000,000,000 |     0 |          |
'-----------------------------------------------------------------------------------------------' 
evmlc$ exit   
```

Generate a new key pair value. 
```bash
$ monetcli keys generate /tmp/keyfile.json
Passphrase: 
Address: 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c
```


```bash
evmlc$ accounts get 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c -f
.-------------------------------------------------------------------------.
|                  Address                   | Balance | Nonce | Bytecode |
|--------------------------------------------|---------|-------|----------|
| 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c |       0 |     0 |          |
'-------------------------------------------------------------------------'

evmlc$ transfer
? From:  46e05762e981d040283af871dca60a71a6786a23
? Enter password:  [hidden]
? To 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c
? Value:  5000
? Gas:  100000000
? Gas Price:  0
Transaction {
  constant: false,
  parseLogs: undefined,
  unpackfn: undefined,
  from: '46e05762e981d040283af871dca60a71a6786a23',
  to: '7b86a2be73108a94d54c0fd2a52676425ace270c',
  value: 5000,
  data: '',
  gas: 100000000,
  gasPrice: 0,
  nonce: undefined,
  chainId: 1 }
? Submit transaction Yes
Transaction submitted successfully.
evmlc$ accounts get 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c -f
.-------------------------------------------------------------------------.
|                  Address                   | Balance | Nonce | Bytecode |
|--------------------------------------------|---------|-------|----------|
| 0x7B86a2BE73108a94D54C0Fd2a52676425aCE270c | 5,000   |     0 |          |
'-------------------------------------------------------------------------'
evmlc$ exit

```


//TODO flesh this section out.

+ run crowdfunding etc.





# Joining a Network

This scenario is for when you wish to join an existing network that is already running, such as the one created in the previous example. This scenario is designed to be run on a machine other than the one is the running the existing node. 

## Clear Previous Configurations
The following command is not necessary for a clean install, but it is good practice to run 
it anyway. This command renames any previous .monet configuration with a .\~n\~ suffix and creates a new empty configuration folder.  

**NB if you run this command after running the previous example, it will move the configuration files from the previous example, breaking the conguration of the previous node**
```bash
monetcli config clear
```

## Create An Account

As for creating a new network, you need to generate your key pair for your account, exactly as per when creating a new network. The syntax of the command is: 

```bash
$ monetcli keys new [--moniker moniker] [--passphrase eth.txt]
```
You need to supply a moniker for this node. You will also need to supply a passphrase to secure your key. You can either use the `--passphrase` parameter to specify a file containing the passphrase, or you can specify interactively. 

You need to supply a moniker for this node. You will also need to supply a passphrase to secure your key. You can either use the `--passphrase` parameter to specify a file containing the passphrase, or you can specify interactively. 

Thus we need to run:
```bash
$ monetcli keys new --moniker node1 
Enter Passphrase:
Re-Enter Passphrase:
Address: 0xDd9C70C8a02D1D47c4423850b1bDc7C3bbb43422
```

//TODO update this tree
```bash
.monet
├── accounts
│   └── node0
│       ├── keyfile.json
│       └── peer.toml

```


## Pull the Configuration

We now pull the monetd configuration files from an existing peer. The syntax for this command is:

```bash
$ monetcli config pull [--peer peer_address] [--node node0] [--address ip]
```

The `--peer` parameter is the address / ip of an existing node on the network. The network's configuration is requested from this peer. 

We need the IP address for the node you building a network upon. For a live network that would clearly be a public IP address, but for an exploratory test net, we would recommend using an internal IP address. On Linux `ifconfig` will give you IP address information. 

Thus we need to run this command, but replace `192.168.1.4` with the IP address / hostname you discover above, and replace `192.168.1.5` with the address of the existing peer. 

```bash
$ monetcli config pull --peer 192.168.1.4 --node node1  --address 192.168.1.5 
```

This command builds the configuration files for a monetd node. It adds the lists of nodes given by the `--nodes` parameter to the initial peer set, including adding them to the initial validator whitelist in the POA smart contract. 

If the `--node` parameter is missing or the value does not correspond to a key pair previously defined, then the list of valid nodes is shown.

If the `--address` parameter is missing, a best guess IP is shown. 

The configuration files are written to a `.monet` folder in $HOME. You should have a structure like this at this stage.

```bash
# //TODO After writing the code, insert the relevant tree here. 
$ tree $HOME/.monet
.monet
├── accounts
│   └── node0
│       ├── keyfile.json
│       └── peer.toml

```

## Apply to Join the Network

You next need to apply to join the network. 

```bash
# //TODO Initial version would be to apply just this 
$ evmlc poa nominate -h <existing node> --from <your address> --moniker <your moniker> --pwd <passphrase file> <your address>
```

```bash 
$ evmlc poa nominate -h 192.168.1.4 --from 0x967c3fE635d2a1e3098b58342D96D74cdD4bf792 --moniker node1  --pwd ~/.monet/eth/pwd.txt 0x967c3fE635d2a1e3098b58342D96D74cdD4bf792
```

The existing node needs to start before it can approve your node joining. **On the existing instance**:

```bash
$ evmlc i
  _____  __     __  __  __           _       _   _               ____   _       ___ 
 | ____| \ \   / / |  \/  |         | |     (_) | |_    ___     / ___| | |     |_ _|
 |  _|    \ \ / /  | |\/| |  _____  | |     | | | __|  / _ \   | |     | |      | | 
 | |___    \ V /   | |  | | |_____| | |___  | | | |_  |  __/   | |___  | |___   | | 
 |_____|    \_/    |_|  |_|         |_____| |_|  \__|  \___|    \____| |_____| |___|
                                                                                    
 Mode:        Interactive
 Data Dir:    /home/user/.evmlc
 Config File: /home/user/.evmlc/config.toml
 Keystore:    /home/user/.evmlc/keystore

  Commands:

    help [command...]                    Provides help for a given command.

...

evmlc$ poa vote
? From:  0f4b70c732aa6b03db3724c9d893e85c7c5e218a
? Passphrase:  [hidden]
? Nominee:  0x967c3fE635d2a1e3098b58342D96D74cdD4bf792
? Verdict:  Yes
You (0x0f4b70c732aa6b03db3724c9d893e85c7c5e218a) voted 'Yes' for '0x967c3fe635d2a1e3098b58342d96d74cdd4bf792'. 
Election completed with the nominee being 'Accepted'.

```

## Starting the Node

To start running the monetd node in a terminal window run:

```bash
$ monetd run
```

If you are not a validator on this network, monetd asks another peer if you are on the whitelist. If you are, it starts running. If not, it checks to see if you are on the nominee list and exits with a suitable message either telling you to apply to join the network, or confirming that voting is not yet complete. 


# Creating A More Complex Network

First clear down any pre-existing network configuration.
```bash
$ monetcli config clear
```


## Generating A New Key Pair

```bash
$ monetcli keys new [--moniker moniker] [--passphrase eth.txt]
```

## Importing Existing Key Pairs

Imports an existing keyfile allow configuration of an existing account.

```bash
$ monetcli keys import [--keyfile keyfile] [--moniker moniker] [--passphrase eth.txt]
```

## Key Pair Storage

A structure with 2 nodes create a file structure like this:  

```bash
.monet
├── accounts
│   ├── node0
│   │   ├── keyfile.json
│   │   └── peer.toml
│   └── node1
│       ├── keyfile.json
│       └── peer.toml
```

The subfolders of accounts are "safe" versions of the moniker --- i.e. stripped of spaces and special characters. The `peer.toml` file initially contains only the moniker. It would be possible to add items such as genesis block balances and whether they are pre-approved validators. 


## Outline

+ Generate 5 key pairs
+ Create a network with 3 of them. 




# Dev notes

## Further `monetcli config build` options 

+ `--peers` comma separated list of initial peers **excluding** the one specified by `--node`
+ `--peer-address` comma separated list of initial peers **excluding** the one specified by `--node`


So to configure a network of nodes: node0, node1, node2, node3 the command would be:
```bash
$ monetcli config build --node node0  --address 192.168.1.4 --peers node1,node2,node3 --peer-address host1,host2,host3
```

+ `--defaultbalance` - sets the default balance for an account
+ `--alloc-peers-only` - only funds named peers in the genesis block. Default is to fund all key pairs in the configuration



## Development TODO

+ ~~monetcli keys new~~

+ ~~monetcli config clear~~

+ monetcli config build

    + If the `--nodes` parameter is missing or any of the values do not correspond to a key pair previously defined, then the list of valid nodes is shown.

    + If the `--ips` parameter is missing, a best guess IP is shown. 

    + Modify default log level to be higher level than debug - probably info. We want the monetd output to be less intimidating. 

    + set evmlc config by default 
+ monetcli config pull


+ monetcli key import

