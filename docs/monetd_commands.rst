.. _monetd_commands_rst:

Monetd Commands
===============

Monetd provides the core commands needed to configure and run a Monet
node. Monetd has context sensitive help accessed either by
running ``monetd help`` or by adding a ``-h`` parameter to the relevant
command. 

.. code:: bash
    $ monetd help
    MONET-Daemon
        
    Monetd provides the core commands needed to configure and run a Monet
    node. The minimal quickstart configuration is:

        $ monetd config clear
        $ monetd keys new node0
        $ monetd config build node0
        $ monetd run

    See the documentation at https://monetd.readthedocs.io/ for further information.

    Usage:
    monetd [command]

    Available Commands:
    config      manage monetd configuration
    help        Help about any command
    keys        monet key manager
    run         run a MONET node
    version     show version info

    Flags:
    -d, --datadir string   Top-level directory for configuration and data (default "/home/jon/.monet")
    -h, --help             help for monetd
        --log string       trace, debug, info, warn, error, fatal, panic (default "debug")
    -v, --verbose          verbose messages

    Use "monetd [command] --help" for more information about a command.


There are 5 subcommands. ``help`` is described above. The other 4 commands are
described in separate sections below:

- **help** --- show help for the command and subcommands
- **version** --- shows the current version of monet and subsystems
- **keys** --- creates and manages keys
- **config** --- creates and manages configurations
- **run** --- runs the monet daemon, i.e. starts a node


Global Parameters
-----------------

Global Parameters are available for all subcommands.

- **-d, --datadir string** --- overrides the default location of the configuration file
- **-h, --help** --- help command as discussed above
- **--log string** --- Set to one of trace, debug, info, warn, error, fatal, panic. Selects
  the logging level for daemon subcommands such as ``run``, the higher the level, the less events will be logged. 
- **-v, --verbose** --- turns on verbose messages for the non-daemon commands. 
  Defaults to false.


Version
-------

The version command outputs the version number for ``monetd``, ``EVM-Lite``, 
``Babble`` and ``Geth``. 

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash. 

.. code:: bash

    $ monetd version
    Monet Version: 0.2.1-develop-ceb36cba
        EVM-Lite Version: 0.2.1
        Babble Version: 0.5.0
        Geth Version: 1.8.27


Keys
----

The keys subcommand is used to manage monet keys. There are 4 subcommands, each
described in a seperate section below:

- **inspect** --- inspect a keyfile
- **list** --- list keyfiles
- **new** --- create a new keyfile
- **update** --- change the passphrase on a keyfile

The keys subcommand writes and reads keys from the ``keystore`` sub-folder in the
monet configuration folder. You can see the location for your instance with this
command:

.. code:: bash

    $ monetd config location -x

The help for the keys command is:

.. code:: bash

    monetd help keys

    This command manages keys in the [datadir]/keystore folder.

    Each key is associated with a moniker and encrypted in a password protected 
    file. The moniker is a friendly name preventing users from having to type or 
    copy/paste long character strings in the terminal. The password-protected file 
    contains a JSON formatted string, which Ethereum users will recognise as the 
    de-facto Ethereum keyfile format. Indeed, Monet and the underlying consensus 
    algorithm, Babble, use the same type of keys as Ethereum. The same key can be 
    used to run a validator node, or to control an account in Monet with a Tenom 
    balance.

    To use a key as part of a validator node running monetd, it will have to be 
    decrypted with the password and copied over to [datadir]/babble/priv_key. The 
    command  'monetd config build' does this automatically, but it can also be done 
    manually with the help of the 'monetd keys inspect --private' command. 

    Note that other Monet tools, like monetcli and monet-wallet, use the same 
    default [datadir]/keystore.

    +------------------------------------------------------------------------------+ 
    | Please take all the necessary precautions to secure these files and remember | 
    | the password, as it will be impossible to recover the key without them.      |
    +------------------------------------------------------------------------------+

    Usage:
    monetd keys [command]

    Available Commands:
    inspect     inspect a keyfile
    list        list keyfiles
    new         create a new keyfile
    update      change the passphrase on a keyfile

    Flags:
    -h, --help              help for keys
        --json              output JSON instead of human-readable format
        --passfile string   file containing the passphrase

    Global Flags:
    -d, --datadir string   Top-level directory for configuration and data (default "/home/jon/.monet")
        --log string       trace, debug, info, warn, error, fatal, panic (default "debug")
    -v, --verbose          verbose messages

    Use "monetd keys [command] --help" for more information about a command.

Parameters
~~~~~~~~~~

All of the keys subcommands support the ``--passfile`` flag. This allows you to 
pass the path to a plain text file containing the passphrase for your key. This
removes the interactive prompt to enter the passphrase that is the default mechanism. 


Monikers
~~~~~~~~

Keys generated by monetd have a moniker associated with them. The moniker is used
to manage the keys as it is far more user friendly that an Ethereum address or
public key. The moniker is stored in 2 version. One is unchanged from how it 
is entered. The other "safe version" is used to generate file numbers. The safe
version is created by replacing all non alphanumeric characters with underscores. 
N.B. Alphanumeric letters in this context only includes unaccented standard 
Roman letters (i.e. just 26 uppercase and 26 lowercase). You will be prevented from
creating a 2nd key which produces the same "safe" moniker. You can use either the
unchanged or safe moniker in any command (the safe version is used internally), 
but if using the unchanged version care should be taken to use suitable quotes.


New
~~~

The ``new`` subcommand generates a new key pair and associates it with the specified moniker. 
You will be prompted for a passphrase which is used to encrypt the keyfile. 
It writes the encrypted keyfile to the Monetd keystore area by default. The moniker
used must be unique within your keystore. If you attempt to create a duplicate, 
the command will abort with an error. 

.. code:: bash

    $ monetd keys new -h

    This command generates a new cryptographic key-pair, and produces two files:

    - [datadir]/keystore/[moniker].json : The encrypted keyfile
    - [datadir]/keystore/[moniker].toml : Key metadata

    [moniker] is a friendly name, which can be reused in other commands to refer to 
    the key without having to type or copy a long string of characters.

    If the --passfile flag is not specified, the user will be prompted to enter the
    passphrase manually. Otherwise, it will be read from the file pointed to by
    --passfile.

    Usage:
    monetd keys new [moniker] [flags]

The moniker supplied in the command above must be in the list of moniker 
produced by ``monetd keys list``.


.. code:: bash

    $ monetd keys new node0 
    Passphrase: 
    Repeat passphrase: 
    Address: 0x14f066E56969F10a9fc95065eA8E3Bd36cf51d13


Inspect
~~~~~~~

The ``inspect`` subcommand interrogates an encrypted keyfile and returns the 
public key and address. If you specify the ``--private parameter``, it also 
returns the associated private key.

.. code:: bash

    $ monetd keys inspect -h

    The inspect subcommand interrogates an encrypted keyfile and returns the 
    public key and address. If you specify the --private parameter, it also 
    returns the associated private key.

    Usage:
    monetd keys inspect [moniker] [flags]

    Flags:
    -h, --help      help for inspect
        --private   include the private key in the output




.. code:: bash

    $ monetd keys inspect node0 --private
    Passphrase: 
    Address:        0x02f6f3D24E447218d396C14F3B47f9Ea369DADf9
    Public key:     0481d3528eec6138f8428932e4fe99571a4f77bd79ae13219540b0a929014cb490a4e5ced2f9e651b531522c2567b6dc5de75d485193615e768b8aa1190603d2c2
    Private key:    bc553aaa7e55c5d0f58f6897ba9bffdb88233c420da622d363f2fe4bd6d78df1
    jon@hpjon:~/go/src/github.com/mosaicnetworks/monetd$ monetd keys inspect node0 
    Passphrase: 
    Address:        0x02f6f3D24E447218d396C14F3B47f9Ea369DADf9
    Public key:     0481d3528eec6138f8428932e4fe99571a4f77bd79ae13219540b0a929014cb490a4e5ced2f9e651b531522c2567b6dc5de75d485193615e768b8aa1190603d2c2

Update
~~~~~~

The ``update`` subcommand allows you to change the passphrase for an encrypted
key file. You are prompted for the old passphrase, then you need to enter, and 
confirm, and new passphrase.

.. code:: bash

    $ monetd keys update node0 
    Passphrase: 
    Please provide a new passphrase
    Passphrase: 
    Repeat passphrase: 

List
~~~~

The ``list`` subcommand outputs a list of the nodes in your keystore. It provides a list of the valid nodes
that can be specified to the other keys subcommands.

.. code:: bash

    $ monetd keys list
    node0
    node1
    node2





Config 
------

Run
---


