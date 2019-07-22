.. _monetd_keys_rst

Keys
====

The keys subcommand is used to manage ethereum keys.

.. code:: bash
    $ monetd help keys
    Using config file: /home/user/.monet/monetd.toml
    Keys

    Monet Key Manager.

    Usage:
    monetd keys [command]

    Available Commands:
    inspect     inspect a keyfile
    new         create a new keypair
    update      change the passphrase on a keyfile

    Flags:
    -h, --help              help for keys
        --json              output JSON instead of human-readable format
        --passfile string   the file that contains the passphrase for the keyfile

    Global Flags:
    -d, --datadir string   Top-level directory for configuration and data (default "/home/jon/.monet")
        --log string       debug, info, warn, error, fatal, panic (default "debug")

    Use "monetd keys [command] --help" for more information about a command.


New
---

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
-------

The **inspect** subcommand interrogates an encrypted keyfile and returns the 
public key and address. If you specify the --privatekey parameter, it also 
returns the associated private key.


Update
------

The **update** subcommand allows you to change the passphrase for an encrypted
key file.


