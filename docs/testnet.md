#monetcli testnet

This command builds the configuration for a new test net on separate devices

We invoke a server to co-ordinate between the peers and to share the genesis information for your testnet. Then all the initial (and subsequent peers) can download their configuration from that server. 

The testnet subcommand facilitates moving the configuration files into the right folders, and ensuring that the same initial configuration is shared between all of the nodes.

## Configuration Server Installation and Invocation
One **one** machine only

We install as follows:
```bash
[...monetd] $ make installcfg 
```

Make a note of the IP address of the this device (ifconfig will tell you), as the peers will need it.

We can then invoke it:
```bash
$ monetcfgsrv
```

Just leave this window open and the server running. 

## Configuring Peers

On each peer in a new terminal session run:

```bash
$ monetcli testnet
```

If you have run the testnet subcommand before you may get this message:

```bash
$ monetcli testnet
This is a destructive operation. Remove/rename the following folder to proceed.
/home/jon/.monetcli/testnet
```
If you do, you need to delete / rename that directory. The testnet subcommand will not overwrite it automatically. 


The first question asked is the address of the configuration server. If you are not running it on your device replace localhost with the IP address of the running monetcfgsrv. Leave the :8088 on the end. It is required.

```bash
$ monetcli testnet
✔ Configuration Server:  : |http://localhost:8088
```

Next you are asked to enter, and re-enter to confirm, a passphrase to secure your keys. Do not lose this phrase as you will not be able to use the account if you do.

```
Enter Keystore Password:   : ######|
✔ Confirm Keystore Password:   : ######|
```

Next you are asked to enter a moniker to identify your node by:
```
✔ Enter your moniker:   : Jon|

```

Next you are asked to enter your IP. Our best guess is pre-filled as the default and can usually be accepted.
```
✔ Enter your ip without the port:   : |192.168.1.18
```

The program now generates a keypair for you and places you in a holding menu.
```
Address: 0x7C86f94E113d9E957a42442765Cd06969ABB1bef
Building Data to push to Configuration Server
Moniker  :  Jon
IP       :  192.168.1.18
Pub Key  :  04893ea962c86923931c99f0915cae9ca74245e3a1ee949b5e7a65eb20ff1e00601f33bc29400f522744b142b36ecc54a5b37e38a712405dba44bf5673bbfb0543
Address  :  0x7C86f94E113d9E957a42442765Cd06969ABB1bef
URL      :  http://localhost:8088/addpeer
response Status: 200 OK
response Headers: map[Content-Length:[4] Content-Type:[text/plain; charset=utf-8] Date:[Mon, 01 Jul 2019 11:37:39 GMT]]
response Body: true
Use the arrow keys to navigate: ↓ ↑ → ← 
? Choose your action  : 
  ▸ Check if published
    Publish, no more initial peers will be allowed to be added
    Exit

```
Check if published polls the monetcfgsrv to check whether the configuration for this network has been published. 

You should get each of your initial set of peers to this stage before one of them selects the Publish... option. You can see the peers that have been created by viewing the web page: http://localhost:8088/peersjson where localhost can be replaced with the IP address of the device running monetcfgsrv.

When the set is complete on one device select Publish. This device should have solc installed and accessible from the command line. The following command will generate an error if that is not the case. 

```
$ solc --version
```


You get asked for your IP. It should be prefilled correctly. 
```
✔ Publish, no more initial peers will be allowed to be added
Getting peers.json
Unmarshalling peers.json
Peers list unmarshalled:  1 [0xc0002e6b10]
Adding...  Jon
response Status: 200 OK
response Headers: map[Content-Length:[4] Content-Type:[text/plain; charset=utf-8] Date:[Mon, 01 Jul 2019 11:57:38 GMT]]
response Body: true
Publish result: true

Configuration has been published.
Getting peers.json
Getting genesis.json
✔ Enter your ip without the port:   : |192.168.1.18
```

It then downloads all of the configuration files, then prompts you for a confirmation to write them into place.
```
All files downloaded
Use the arrow keys to navigate: ↓ ↑ → ← 
? Confirm Overwriting Existing Configuration  : 
  ▸ No
    Yes
```

It then copies all the files into place, finally giving you the command to start a monetd server. 
```
✔ Yes
Renaming /home/jon/.monet to /home/jon/.monet.~9~
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

Start the server using:
```
monetd run
```

On all the other nodes, select  Check if published, which will then perform the same workflow (minus the smart contract compilation - they will use the genesis.json file generated above. 


## Developer Details



The testnet wizard places its files in $HOME/.monetcli/testnet (on Linux, other OS may vary) as below. The files marked (*) are only created on the device that publishes the network details.   

```
├── contract0.abi   (*)
├── contract0.sol   (*)
├── genesis.json
├── keyfile.json
├── monetd.toml
├── network.toml    (*)
├── peers.genesis.json
├── peers.json
├── priv_key
└── pwd.txt
```


