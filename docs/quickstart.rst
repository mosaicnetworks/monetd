.. _quickstart_rst:

Getting Started
===============

Creating A Single Node Network
------------------------------

The simplest configuration is a single node.

To start a single node, run the following commands:

.. code::

    $ monetd keys new node0
    $ monetd config build node0
    $ monetd run

Let's walk through it step by step.

Clear Previous Configurations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following command is not necessary for a clean install, but it is good
practice to run it anyway. This command renames any previous .monet 
configuration with a .~n~ suffix and creates a new empty configuration folder.

.. code:: bash

    monetd config clear

Create An Account
~~~~~~~~~~~~~~~~~

Firstly we need to create an account for the validator. We will generate a 
private and public keypair. The syntax for the ``monetd keys new`` command is:

.. code:: bash

    monetd keys new [moniker] [--passfile]

You need to supply a moniker for this node. You will also need to supply a
passphrase to secure your key. You can either use the ``--passfile`` parameter
to specify a file containing the passphrase, or you can specify it
interactively.

Thus we need to run:

.. code:: bash

    $ monetd keys new node0 
    Enter Passphrase:
    Re-Enter Passphrase:
    Address: 0x83434e68b52Ef809538224BF78472cc3F6a17bcC

The resulting encrypted file, and all other configuration, is written to a 
``.monet`` folder. On Linux it is ``$HOME/.monet``. You can find the location
for your instance with this command:

.. code:: bash

    $ monetd config location

At this stage, you should have a structure like this:

.. code:: bash

    $ tree $HOME/.monet
    .monet
    └── keystore
        └── node0.json

To view the existing list of keys, run ``monetd keys list``:

.. code:: bash

    $ monetd keys list
    node0

Build the Configuration
~~~~~~~~~~~~~~~~~~~~~~~

We now build the monetd configuration files. The syntax for this command is:

.. code:: bash

    $ monetd config build [moniker] [--address ip] [--passfile file]

We need the IP address for the node you are building a network upon. For a live 
network that would clearly be a public IP address, but for an exploratory 
testnet, we would recommend using an internal IP address. On Linux ``ifconfig`` 
will give you IP address information. If you omit the --address parameter, 
monetd will pick your first non-loopback address. The --passfile parameter 
specifies a file containing your passphrase. We would recommnend using the 
interactive prompt to enter the passphrase.

.. code:: bash

    $ monetd config build node0  

This command builds the configuration files for a monetd node. It adds the 
account referenced by [moniker] to the initial peer set, including adding it to 
the initial validator whitelist in the POA smart contract.

The location of the configuration files depend on the OS. On Linux it is 
``$HOME/.monet``. You can find the location for your instance with this command:

.. code:: bash

    $ monetd config location
    /home/user/.monet

At this stage, it should look something like this:

.. code:: bash

    $ tree $HOME/.monet
    .monet
    ├── babble
    │   ├── peers.genesis.json
    │   ├── peers.json
    │   └── priv_key
    ├── eth
    │   ├── genesis.json
    │   └── poa
    │       ├── compile.toml
    │       ├── contract0.abi
    │       └── contract0.sol
    ├── keystore
    │   └── node0.json
    └── monetd.toml

Starting the Node
~~~~~~~~~~~~~~~~~

To start running the monetd node in a terminal window run:

.. code:: bash

    $ monetd run

This is clearly not a production configuration, where you would use ``nohup`` 
and redirect log output to the files.

Testing
~~~~~~~

Let us use MONET-CLI to to query the newly created node. First of all, install
monetcli with ``npm install -g monetcli``. For more detailed instructions,
please refer to :ref:`clients_rst`.

While monetd is still running, open another terminal and start MONET-CLI in 
interactive mode to run a couple of commands:

.. code:: bash

    monetcli i
     __  __                          _        ____   _       ___ 
    |  \/  |   ___    _ __     ___  | |_     / ___| | |     |_ _|
    | |\/| |  / _ \  | '_ \   / _ \ | __|   | |     | |      | | 
    | |  | | | (_) | | | | | |  __/ | |_    | |___  | |___   | | 
    |_|  |_|  \___/  |_| |_|  \___|  \__|    \____| |_____| |___|
                                                                 
    Mode:        Interactive
    Data Dir:    /home/user/.monet
    Config File: /home/user/.monet/monetcli.toml
    Keystore:    /home/user/.monet/keystore
   
     Commands:
   
       help [command...]                    Provides help for a given command.
       exit                                 Exits application.
       accounts create [options]            Creates an encrypted keypair locally
       accounts get [options] [address]     Fetches account details from a connected node
       accounts list [options]              List all accounts in the local keystore directory
       accounts update [options] [address]  Update passphrase for a local account
       accounts import [options]            Import an encrypted keyfile to the keystore
       config set [options]                 Set values of the configuration inside the data directory
       config view [options]                Output current configuration file
       poa check [options] [address]        Check whether an address is on the whitelist
       poa info [options]                   Display Proof of Authority information
       poa nominate [options] [address]     Nominate an address to proceed to election
       poa nomineelist [options]            List nominees for a connected node
       poa vote [options] [address]         Vote for an nominee currently in election
       poa whitelist [options]              List whitelist entries for a connected node
       transfer [options]                   Initiate a transfer of token(s) to an address
       info [options]                       Display information about node
       version [options]                    Display current version of cli
       debug                                Toggle debug mode
       clear                                Clear output on screen

    monetcli$ info
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
    
    monetcli$ accounts list
    .-----------------------------------------------------------------------------.
    |                  Address                   |        Balance         | Nonce |
    |--------------------------------------------|------------------------|-------|
    | 0xa10aae5609643848fF1Bceb76172652261dB1d6c | 1234567890000000000000 |     0 |
    '-----------------------------------------------------------------------------'

    monetcli$ accounts get 0xa10aae5609643848fF1Bceb76172652261dB1d6c
    .-----------------------------------------------------------------------------------------------.
    |                  Address                   |            Balance            | Nonce | Bytecode |
    |--------------------------------------------|-------------------------------|-------|----------|
    | 0xa10aae5609643848fF1Bceb76172652261dB1d6c | 1,234,567,890,000,000,000,000 |     0 |          |
    '-----------------------------------------------------------------------------------------------' 
       

So we have a prefunded account. The same account is used as a validator in
Babble, and as a Tenom-holding account in the ledger. This is the same account, 
node0, that we just created in the previous steps, with the encrypted private
key residing in ~/.monet/keystore.

Now, let's create a new key using monetcli, and transfer some tokens to it.

.. code:: bash

    monetcli$ accounts create                                                                                                                                   
    ? Passphrase:  [hidden]                                                                                                                                  
    ? Re-enter passphrase:  [hidden]                                                                                                                         
    {"version":3,"id":"89970faf-8754-468e-903c-c9d3248a08cc","address":"960c13654c477ac1d2d7f8fc7ae84d93a2225257","crypto":{"ciphertext":"7aac819c1bed442d778
    97b690e5c2f14416589c7bdd6bdd2b5df5d03584ce0ec","cipherparams":{"iv":"3d15a67d76293c3b7123f2bde76ba120"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams"
    :{"dklen":32,"salt":"730dd67f175a77c9833a230e334719292cbb735607795b1b84484e3d04783510","n":8192,"r":8,"p":1},"mac":"7535c31c277a698207d278cd1f1df90747463
    e390b822cfef7d2faf8f1fa1809"}} 

Like ``monetd keys new`` this command created a new key and wrote the encrypted
keyfile in ~/.monet/keystore. Let's double check that the key was created and 
transfer 100 tokens to it.

.. code:: bash

    monetcli$ accounts list
   .-----------------------------------------------------------------------------.
   |                  Address                   |        Balance         | Nonce |
   |--------------------------------------------|------------------------|-------|
   | 0x960c13654C477Ac1D2d7f8FC7Ae84D93A2225257 | 0                      |     0 |
   | 0xa10aae5609643848fF1Bceb76172652261dB1d6c | 1234567890000000000000 |     0 |
   '-----------------------------------------------------------------------------'

    monetcli$ transfer
    ? From:  a10aae5609643848ff1bceb76172652261db1d6c
    ? Enter password:  [hidden]
    ? To 0x960c13654C477Ac1D2d7f8FC7Ae84D93A2225257
    ? Value:  100
    ? Gas:  1000000
    ? Gas Price:  0
    { from: 'a10aae5609643848ff1bceb76172652261db1d6c',
      to: '960c13654c477ac1d2d7f8fc7ae84d93a2225257',
      value: 100,
      gas: 1000000,
      gasPrice: 0 }
    ? Submit transaction Yes
    Transaction submitted successfully.

    monetcli$ accounts list
    .-----------------------------------------------------------------------------.
    |                  Address                   |        Balance         | Nonce |
    |--------------------------------------------|------------------------|-------|
    | 0x960c13654C477Ac1D2d7f8FC7Ae84D93A2225257 | 100                    |     0 |
    | 0xa10aae5609643848fF1Bceb76172652261dB1d6c | 1234567889999999999900 |     1 |
    '-----------------------------------------------------------------------------'

Joining a Network
-----------------

This scenario is for when you wish to join an existing network that is already 
running, such as the one created in the previous example.

To join an existing monetd network, run the following commands:

.. code::

    $ monetd keys new node1
    $ monetd config pull [address]:[port] --key node1
    $ monetcli poa nominate -h [address] -p [port] --from [node1 address] --pwd [password file for node1 key] --moniker node1 [node1 address]

    # wait to be accepted in the whitelist, which can be checked with
    $ monetd poa whitelist
    # or
    $ monetd poa nomineelist

    $ monetd run

Where [address] and [port] correspond to the endpoint of an existing peer in the
network. 

**This scenario is designed to be run on a machine other than the one that is 
running the existing node.**

Clear Previous Configurations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following command is not necessary for a clean install, but it is good 
practice to run it anyway. It renames any previous .monet configuration with a 
.~n~ suffix and creates a new empty configuration folder.

**NB if you run this command after running the previous example, it will move 
the configuration files from the previous example, breaking the configuration of 
the previous node**

.. code:: bash

    monetd config clear

Create An Account
~~~~~~~~~~~~~~~~~

You need to generate your key pair for your account, exactly as per when 
creating a new network. This time, we will override the default configuration 
directory. The syntax of the command is:

.. code:: bash

    $ monetd -d [datadir] keys new [moniker] [--passfile]

Thus we need to run:

.. code:: bash

    $ monetd -d ~/.monet2 keys new node1
    Passphrase:
    Repeat passphrase:
    Address: 0x5a735fC1235ce1E60eb5f9B9BCacb643a9Da27F4

Pull the Configuration From the Existing Node
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

We now pull the monetd configuration files from an existing peer. The syntax for
this command is:

.. code:: bash

    $ monetd -d [datadir] config pull [peer] [--key] [--address]

The peer parameter is the address / ip of an existing node on the network. The 
network's configuration is requested from this peer. If the address does not 
specify a port, the default API port (8080) is assumed. 

We need the IP address for the node you are building a network upon. For a live 
network that would clearly be a public IP address, but for an exploratory 
testnet, we would recommend using an internal IP address. On Linux ``ifconfig`` 
will give you IP address information. This can be set by using the --address 
flag. If not specified monetd will pick the first non-loopback address. 

The ``--key`` parameter specifies the keyfile to use by moniker. 

Thus we need to run the following command, but replace ``192.168.1.5:8080`` with 
the endpoint of the existing peer.

.. code:: bash

    $ monetd -d ~/.monet2 config pull 192.168.1.5:8080 --key node1  

Apply to Join the Network
~~~~~~~~~~~~~~~~~~~~~~~~~

If we tried to run monetd at this stage, it would not be allowed to join the 
other node because it isn't whitelisted yet. So we need to apply to the 
whitelist first.

We do so with the MONET-CLI ``poa nominate`` command. The syntax is:

.. code:: bash

    $ monetcli poa nominate -h <existing node> --from <your address> --moniker <your moniker> --pwd <passphrase file> <your address>

But we can also do it interactively. **On the existing instance (node0), run the
following interactive monetcli session**:

.. code:: bash

    monetcli i
    __  __                          _        ____   _       ___ 
   |  \/  |   ___    _ __     ___  | |_     / ___| | |     |_ _|
   | |\/| |  / _ \  | '_ \   / _ \ | __|   | |     | |      | | 
   | |  | | | (_) | | | | | |  __/ | |_    | |___  | |___   | | 
   |_|  |_|  \___/  |_| |_|  \___|  \__|    \____| |_____| |___|
                                                                
   Mode:        Interactive
   Data Dir:    /home/user/.monet
   Config File: /home/user/.monet/monetcli.toml
   Keystore:    /home/user/.monet/keystore
  
    Commands:
     [...]
    

    monetcli$ poa nominate
    ? From:  a10aae5609643848ff1bceb76172652261db1d6c
    ? Passphrase:  [hidden]
    ? Nominee:  0x5a735fC1235ce1E60eb5f9B9BCacb643a9Da27F4
    ? Moniker:  node1
    You (0xa10aae5609643848ff1bceb76172652261db1d6c) nominated 'node1' (0x5a735fc1235ce1e60eb5f9b9bcacb643a9da27f4)

    monetcli$ poa nomineelist
    .------------------------------------------------------------------------------.
    | Moniker |                  Address                   | Up Votes | Down Votes |
    |---------|--------------------------------------------|----------|------------|
    | Node1   | 0x5a735fc1235ce1e60eb5f9b9bcacb643a9da27f4 |        0 |          0 |
    '------------------------------------------------------------------------------'

Now that, we have submitted node1 to the whitelist (via node0), we need all the
entities in the current whitelist to vote for it. At the moment, only node0 is
in the whitelist, so let's cast a vote. 

.. code:: bash

    monetcli$ poa whitelist
    .------------------------------------------------------.
    | Moniker |                  Address                   |
    |---------|--------------------------------------------|
    | Node0   | 0xa10aae5609643848ff1bceb76172652261db1d6c |
    '------------------------------------------------------'

    monetcli$ poa vote
    ? From:  a10aae5609643848ff1bceb76172652261db1d6c
    ? Passphrase:  [hidden]
    ? Nominee:  0x5a735fc1235ce1e60eb5f9b9bcacb643a9da27f4
    ? Verdict:  Yes
    You (0xa10aae5609643848ff1bceb76172652261db1d6c) voted 'Yes' for '0x5a735fc1235ce1e60eb5f9b9bcacb643a9da27f4'. 
    Election completed with the nominee being 'Accepted'.

    monet$ poa whitelist
    .------------------------------------------------------.
    | Moniker |                  Address                   |
    |---------|--------------------------------------------|
    | Node0   | 0xa10aae5609643848ff1bceb76172652261db1d6c |
    | Node1   | 0x5a735fc1235ce1e60eb5f9b9bcacb643a9da27f4 |
    '------------------------------------------------------'

Finaly node1 made it into the whitelist.

Starting the Node
~~~~~~~~~~~~~~~~~

To start node2, run the simple ``monetd run`` command. You should be able see
the JoinRequest going through consensus, and being accepted by the PoA contract.

.. code:: bash

    $ monetd -d ~/.monet2 run