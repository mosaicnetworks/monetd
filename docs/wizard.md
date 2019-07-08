# monetcli wizard

The `monetcli wizard`command is an interactive tool for building a monet network configuration.


----

## Table of Contents


+ [Invocation](#invocation)
+ [Configuration Files](#configuration-files)
	+ [Monetcli Network Configuration Folder](#monetcli-network-configuration-folder)
+ [Edit Menu](#edit-menu)
	+ [Add Key Pair](#add-key-pair)
	+ [Check Configuration](#check-configuration)
	+ [Compile POA Contract](#compile-poa-contract)
	+ [Generate Key Pair](#generate-key-pair)
	+ [Edit Params](#edit-params)
	+ [Peers](#peers)
	+ [Show Configuration](#show-configuration)
	+ [Version Info](#version-info)
	+ [Exit](#exit)

----

# Invocation

The tool can be invoked thus:

```bash
$ monetcli wizard
```
 
# Configuration Files
 
## Monetcli Network Configuration Folder
 
You are prompted to enter the location of the configuration folder to use. Unless there is a pressing reason to change it, we would recommend using the default. If you change the location, future executions of `monetcli` would require a `--config-dir` flag to be set. 
 
![screenshot of monetcli wizard](assets/screenshots/wizard1.png)

The default entry is OS dependent and is set to:

+ **Linux:** 		`$HOME/.monetcli 
+ **Windows:** 	`$HOME/AppData/Roaming/MONET 
+ **Mac:** 		`$HOME/Library/MONET 



If the location does not exist, you get the following 3 options:

![screenshot of monetcli wizard](assets/screenshots/wizard2.png)

You can use the arrow keys to highlight an option, then press `return` to select it. 

+ **Exit Without Saving Changes** --- exits the program
+ **Change Configuration Directory** --- takes you back to the first screen to change your answer
+ **Create New Configuration** --- creates a new configuration in the directory specified
 
In most circumstances you will selection **Create New Configuration**.

If the configuration folder already exists, you will be offered 2 additional options. **Rename Current Directory and Create New Configuration** renames the existing folder with a `.~n~` suffix where `n` is the lowest number where the folder does not already exist. If n reaches 100, the program will abort - and we recommend that you do some housekeeping on your filesystem.  **Edit Existing Configuration** does, as the name suggests, edit the existing configuration in the folder you specified.

![screenshot of monetcli wizard](assets/screenshots/wizard3.png)

**Create New Configuration**, **Rename Current Directory and Create New Configuration** and  **Edit Existing Configuration** all progress to the same menu - the latter with existing configuration, the first 2 with a clean configuration. 

# Edit Menu

The Edit menu allows you to build and amend your network configuration. 

![screenshot of monetcli wizard](assets/screenshots/wizard4.png)


## Add Key Pair
The add key pair option allows the addition of peers whose private key had been generated elsewhere. For testnets you would usually use Generate Key Pair, not Add Key Pair. 

![screenshot of monetcli wizard](assets/screenshots/wizard11.png)


## Check Configuration

This option runs checks to ensure that the `network.toml` is in a consistent state. 

![screenshot of monetcli wizard](assets/screenshots/wizard10.png)

## Compile POA Contract

This option generates an actual monet hub configuration. It implicitly runs a network check command. It populates a datadir directory including copying any keys stored within the network configuration folder. If the nodename is specified the configuration for that node is written. It is intended that the node name would allow multiple configurations be generated from the same machine - likely useful for node. The POA contract is compiled to build the genesis block. If there is no bytecode in the network.toml it is added with solcs version. Otherwise the bytecode is validated. **N.B.** this command requires an internet connection to run, unless you have run `monetcli network contract`. The default contract is downloaded directly from github.

[comment]: # (//TODO Add compile section )

## Generate Key Pair

Generate Key Pair generates a new public and private key pair, and adds it as validator to the monetcli network configuration.

![screenshot of monetcli wizard](assets/screenshots/wizard7.png)

The Node Name, or moniker, is a label used to identify this node. It will be rejected if another node with that name has already been added to this configuration.

The node address is the IP and Port that Babble gossips from. 

The passphrase is used to encrypt your private key to secure it. You will need to re-enter the phrase to confirm it. If the 2 do not match, you will be prompted to enter them again. 

The Address generated is echoed to screen. 

## Edit Params

[comment]: # (//TODO Add Edit Params section )

## Peers

[comment]: # (//TODO Add peers section )

## Show Configuration

This option echoes the current `network.toml` file to screen. 

![screenshot of monetcli wizard](assets/screenshots/wizard7.png)


## Version Info

The version info command displays version information. 


![screenshot of monetcli wizard](assets/screenshots/wizard5.png)


## Exit

The exit command closes the program. 

----



----

<sup>[Documents Index](README.md) | [GitHub repo](https://github.com/mosaicnetworks/monetd) | [Monet](https://monet.network/) | [Mosaic Networks](https://www.babble.io/)</sup>
