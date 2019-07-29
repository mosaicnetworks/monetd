.. _giverny_rst:

#######
Giverny
#######

``giverny`` is the advanced configuration tool. 

The current subcommands are:

- **help** --- help
- **version** --- outputs version information
- **keys** --- key management tools
- **server** --- configuration server management
- **network** --- configure and build monet networks

**This documentation is being lovingly crafted now, and will be ready very soon.**

****
Help
****

``giverny`` has context sensitive help accessed either by
running ``giverny help`` or by adding a ``-h`` parameter to the relevant
command. 


*******
Version
*******

The version command outputs the version number for ``monetd``, ``EVM-Lite``, 
``Babble`` and ``Geth``. 

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash. 

.. code:: bash

    $ giverny version
    Monet Version: 0.2.1-develop-ceb36cba
        EVM-Lite Version: 0.2.1
        Babble Version: 0.5.0
        Geth Version: 1.8.27


****
Keys
****

The keys sub-command offers tools to manage keys. 

Import
======

The import sub-command is used to import a pre-existing private key into the
monet keystore, creating the associated ``toml`` file, assigning a moniker and 
setting a passphrase. 

.. code:: bash

    $ giverny help keys import

    Import keys to [moniker] from private key file [keyfile].

    Usage:
    giverny keys import [moniker] [keyfile] [flags]

    Flags:
    -h, --help   help for import

    Global Flags:
    -d, --datadir string    Top-level directory for configuration and data (default "/home/jon/.monet")
        --passfile string   the file that contains the passphrase for the keyfile


******
Server
******

The ``server`` sub-command is used for admninstering a REST server used to co-ordinate 
configurations between multiple nodes prior to the initial node of a network. 

The server listens on port 8088. It writes logs to ``.monet/giverny/server.log``.

For instructions on how to use, see the recipes for setting up networks. 

Start
=====


To start the server in the foreground:

.. code:: bash

    $ giverny server start --background





To start the server in the background:

.. code:: bash

    $ giverny server start --background




Stop
====

To stop a server running in the background: 

.. code:: bash

    $ giverny server stop


Status
======

Reports on the status of the server. It both checks for the PID file in 
``.monet/giverny/server.pid`` and checks the the server is responding on 
localhost:8088. 

.. code:: bash

    $ giverny server status


*******
Network
*******

The *network name* and *node names* must contain only standard letters (i.e. no accented versions), digits (0--9) or
underscores (_). 


New
===

The ``new`` sub-command creates a new test network configuration. 

Syntax
------

.. code:: bash

    $ giverny network new -h
    Created Directory:  /home/jon/.giverny/server

    giverny network build

    Usage:
    giverny network new [network_name] [flags]

    Flags:
        --generate-pass       generate pass phrases
    -h, --help                help for new
        --initial-ip string   initial IP address of range
        --initial-peers int   number of initial peers
        --names string        filename of a file containing a list of node monikers
        --pass string         filename of a file containing a passphrase
        --save-pass           save pass phrase entered on command line

    Global Flags:
    -g, --giverny-data-dir string   Top-level giverny directory for configuration and data (default "/home/jon/.giverny")
    -m, --monet-data-dir string     Top-level monetd directory for configuration and data (default "/home/jon/.monet")
    -n, --nodes int                 number of nodes in this configuration (default 4)
    -v, --verbose                   verbose messages

Nodes
-----

The number of nodes in this network is specified by the 
``--nodes [int]`` parameter. The ``--initial-peers [int]`` parameter specifies
the number of initial peers. If not set it assumes that all nodes are in the 
initial peer set. 

IP Addresses
------------

An initial IP address is supplied using the ``--initial-ip`` parameter. 
It is assumed the IP address range will be assigned by simply incrementing the
last octet of the IP address for each node. N.B. the first node will be assigned
the actual IP supplied by the ``initial-ip`` parameter.


Node Names
----------

The default node names are a standard prefix of *node* with a unique integer 
suffix. You can override the default and supply a list of node names, which are
used in the order supplied, via the ``--names`` parameter. 

Node names must contain only standard Latin alphabet characters (ie *a--z* or 
*A--Z* with no accents), underscores (_), or digits (*0--9*).

Pass Phrases
------------

There are numerous pass phrase flags for the ``new`` subcommand. 

- ``--pass [passfile]`` --- uses the given pass phrase file for all nodes
- ``--generate-pass`` --- generates a unique passphrase for each key pair and
  writes it to a file nodename.txt in the network configuration keystore 
  directory
- ``--save-pass`` --- saves pass phrases in the network configuration keystore 
  directory  

The typical use case scenarios for these flags would be:

- None specified --- you are prompted to enter the passphrase for each node which is not saved
- ``--pass`` only --- the specified pass phrase is used, but not saved in the config folder
- ``--pass`` and ``--save-pass`` --- the specified pass phrase is used **and** saved in the config folder
- ``--generate-pass`` only --- pass phrases are generated and saved
- ``--save-pass`` only --- you are prompted to enter the passphrase for each node, which is saved in the config folder


An example of the new subcommand:

.. code:: bash

    $ giverny network new test11 --names sampledata/names.txt --nodes 7 --pass sampledata/pwd.txt --initial-peers 3 --initial-ip 192.168.1.19

