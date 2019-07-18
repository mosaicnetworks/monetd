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
for the validator. We will generate a private and public keypair. The syntax for the `monetcli key new` command is:

```bash
monetcli key new [--moniker moniker] [--passphrase eth.txt]
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

//TODO flesh this out.

+ connect evmlc
+ create a new key pair
+ transfer some money
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
$ monetcli key new [--moniker moniker] [--passphrase eth.txt]
```
You need to supply a moniker for this node. You will also need to supply a passphrase to secure your key. You can either use the `--passphrase` parameter to specify a file containing the passphrase, or you can specify interactively. 

You need to supply a moniker for this node. You will also need to supply a passphrase to secure your key. You can either use the `--passphrase` parameter to specify a file containing the passphrase, or you can specify interactively. 

Thus we need to run:
```bash
$ monetcli key new --moniker node0 
Enter Passphrase:
Re-Enter Passphrase:
Address: 0x83434e68b52Ef809538224BF78472cc3F6a17bcC
```


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
$ monetcli config pull --peer 192.168.1.5 --node node0  --address 192.168.1.4 
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
$ evmlc poa nominate -h 192.168.1.5 --from <your address> --moniker <your moniker>  <your address>
```

//TODO **Do we instead have a command like this:**
```bash
$ monetcli network apply
```

**Where you would interactively need to enter a passphrase, but all other required info is in the monetd configuration files.**



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
$ monetcli key new [--moniker moniker] [--passphrase eth.txt]
```

## Importing Existing Key Pairs

Imports an existing keyfile allow configuration of an existing account.

```bash
$ monetcli key import [--keyfile keyfile] [--moniker moniker] [--passphrase eth.txt]
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

+ ~~monetcli key new~~

+ ~~monetcli config clear~~

+ monetcli config build

    + If the `--nodes` parameter is missing or any of the values do not correspond to a key pair previously defined, then the list of valid nodes is shown.

    + If the `--ips` parameter is missing, a best guess IP is shown. 

    + Modify default log level to be higher level than debug - probably info. We want the monetd output to be less intimidating. 

    + set evmlc config by default 
+ monetcli config pull


+ monetcli key import

