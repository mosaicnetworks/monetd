.. _monetd_commands_rst:

Monetd Commands
===============




Keys
----

The keys subcommand is used to manage monet keys.

.. code:: bash

    $ monetd help keys

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

    Use "monetd keys [command] --help" for more information about a command.

New
~~~

The **new** subcommand generates a new key pair and associates it with the specified moniker. 
You will be prompted for a passphrase which is used to encrypt the keyfile. 
It writes the encrypted keyfile to the Monetd keystore area by default. 

.. code:: bash

    $ monetd keys new node0 
    Using config file: /home/jon/.monet/monetd.toml
    Passphrase: 
    Repeat passphrase: 
    Address: 0x14f066E56969F10a9fc95065eA8E3Bd36cf51d13


Inspect
~~~~~~~

The **inspect** subcommand interrogates an encrypted keyfile and returns the 
public key and address. If you specify the --private parameter, it also 
returns the associated private key.


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

The **update** subcommand allows you to change the passphrase for an encrypted
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

The **list** subcommand outputs a list of the nodes in your keystore. It provides a list of the valid nodes
that can be specified to the other keys subcommands.

