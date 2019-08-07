.. _giverny_rst:

#################
Giverny Reference
#################

``giverny`` is the advanced configuration tool for the Monet Toolchain.

The current subcommands are:

- **help** --- help
- **version** --- outputs version information
- **keys** --- key management tools
- **server** --- configuration server management
- **network** --- configure and build networks


***********
Global Flag
***********

The ``--verbose`` flag, or ``-v`` for short, turns on extended messages for
each ``giverny`` command.


****
Help
****

``giverny`` has context sensitive help accessed either by
running ``giverny help`` or by adding a ``-h`` parameter to the relevant
command.


*******
Version
*******

The ``version`` subcommand outputs the version number for ``monetd``,
``EVM-Lite``, ``Babble`` and ``Geth``.

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash.

.. code:: bash

    $ giverny version
    Monetd Version: 0.2.1-develop-ceb36cba
        EVM-Lite Version: 0.2.1
        Babble Version: 0.5.0
        Geth Version: 1.8.27


****
Keys
****

The ``keys`` subcommand offers tools to manage keys.

Keys Flags
==========

In addition to the ``--verbose`` flag, the ``keys`` subcommand defines
addtional flags as follows:

.. code:: bash

    Global Flags:
    -g, --giverny-data-dir string   Top-level giverny directory for configuration and data (default "/home/user/.giverny")
        --json                      output JSON instead of human-readable format
    -m, --monet-data-dir string     Top-level monetd directory for configuration and data (default "/home/user/.monet")
        --passfile string           the file that contains the passphrase for the keyfile


Import
======

The ``import`` subcommand is used to import a pre-existing private key into the
``monetd`` keystore, creating the associated ``toml`` file, assigning a moniker
and setting a passphrase.

.. code:: bash

    $ giverny help keys import

    Import keys to [moniker] from private key file [keyfile].

    Usage:
    giverny keys import [moniker] [keyfile] [flags]

    Flags:
    -h, --help   help for import

    Global Flags:
    -g, --giverny-data-dir string   Top-level giverny directory for configuration and data (default "/home/user/.giverny")
        --json                      output JSON instead of human-readable format
    -m, --monet-data-dir string     Top-level monetd directory for configuration and data (default "/home/user/.monet")
        --passfile string           the file that contains the passphrase for the keyfile

******
Server
******

The ``server`` subcommand is used for adminstering a REST server used to
co-ordinate configurations between multiple nodes prior to the initial node of a
network.

The server listens on port 8088. It writes logs to
``~/.giverny/server/server.pid``. [1]_

For usage examples, see the recipes for setting up networks.

Start
=====


To start the server in the foreground:

.. code:: bash

    $ giverny server start


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
``~/.giverny/server/server.pid`` [1]_ and checks the the server is responding
on localhost:8088.

.. code:: bash

    $ giverny server status


*******
Network
*******

The ``network`` command is used to build complex monet networks. The ``new``
command generates the nodes and keys for a network, and automatically calls
the ``build`` command which generates and builds ``genesis.json`` and
``peers.json`` files. You can adjust the network by editting the
``network.toml`` file. The ``location`` command outputs the relevant paths.
The ``push`` command is used to push a giverny network node configuration to a
docker or actual node so it can be used by ``monetd``. ``start``, ``stop`` and
``status`` are used to manage the docker instance.


The *network name* and *node names* must contain only standard letters
(i.e. no accented versions), digits (0--9) or underscores (_).

Location
========

The ``giverny network location`` subcommand takes a single optional parameter
``network_name``. If the network is specified it outputs the location of key
files and folders for that network. If not, only the root giverny configuration
folder is output.

Example without a network name:

.. code:: bash

    $ giverny network location
    /home/user/.giverny

Example with a network specified:

.. code:: bash

    $ giverny network location node7
    Network                 : node7
    Giverny Config Dir      : /home/user/.giverny
    Giverny Networks Dir    : /home/user/.giverny/networks/node7
    Giverny KeyStore Dir    : /home/user/.giverny/networks/node7/keystore
    Peers JSON              : /home/user/.giverny/networks/node7/peers.json
    Genesis JSON            : /home/user/.giverny/networks/node7/genesis.json
    Monetd TOML             : /home/user/.giverny/networks/node7/monetd.toml
    Network TOML            : /home/user/.giverny/networks/node7/network.toml



New
===

The ``new`` subcommand creates a new test network configuration. It also

Syntax
------

.. code:: bash

    $ giverny network new -h
    Created Directory:  /home/user/.giverny/server

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
    -g, --giverny-data-dir string   Top-level giverny directory for configuration and data (default "/home/user/.giverny")
    -m, --monet-data-dir string     Top-level monetd directory for configuration and data (default "/home/user/.monet")
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

- None specified --- you are prompted to enter the passphrase for each node
  which is not saved
- ``--pass`` only --- the specified pass phrase is used, but not saved in the
  config folder
- ``--pass`` and ``--save-pass`` --- the specified pass phrase is used **and**
  saved in the config folder
- ``--generate-pass`` only --- pass phrases are generated and saved
- ``--save-pass`` only --- you are prompted to enter the passphrase for each
  node, which is saved in the config folder


Build
-----

By default ``giverny network new`` will run ``giverny network build``
automatically. This can be disabled by specifying the ``-no-build`` flag.


Examples
--------

An example of the new subcommand:

.. code:: bash

    $ giverny network new test11 --names sampledata/names.txt --nodes 7 --pass sampledata/pwd.txt --initial-peers 3 --initial-ip 192.168.1.19



Build
=====

The ``giverny network build`` subcommand takes a configuration created by the
``new`` subcommand and builds ``peers.json`` and ``genesis.json`` files.

``build`` can be run repeatably safely. It is envisaged that users will edit
the ``network.toml`` file to adjust token allocations or change addresses.

A "built" network will have a file structure like this:

.. code:: bash

    test7
    ├── compile.toml
    ├── contract0.abi
    ├── contract0.sol
    ├── genesis.json
    ├── keystore
    │   ├── Amelia.json
    │   ├── Amelia.txt
    │   ├── Becky.json
    │   ├── Becky.txt
    │   ├── Chloe.json
    │   ├── Chloe.txt
    │   ├── Danu.json
    │   ├── Danu.txt
    ├── monetd.toml
    ├── network.toml
    └── peers.json

Export
======

The ``export`` subcommand takes a configuration that has been generated and
exports it to the exports subfolder of the giverny configuration folders as a
zip file. The ``network export`` command has a mandatory network name
parameter, and optionally one or more node names. If the node names are
omitted, all of the nodes for that network are exported.


Import
======

The ``import`` subcommand takes a configuration previously exported by the
``export`` and configures ``monetd`` to use the new configuration. You will
always need to specify a network name and a node name for the import. The
source for the import can be configured thus:

- ``--from-exports`` --- from the exports subfolder in the giverny
  configuration folders. This is the default output location for the ``export``
  command.
- ``--server`` --- from a giverny server. The giverny server will look in the
  exports subfolder in the giverny configuration folders on the instance it is
  running on. N.B. do not run the giverny server on any instance with live
  key pairs or sensitive configuration, as it may be exposed.
- ``--dir`` --- specify the folder the export zip is in. Do not rename the zip
  file. This is used when a secondary channel is used to communicate the keys.



.. [1] This location is for Linux instances. Mac and Windows uses a different
       path. The path for your instance can be ascertain with this command:
       ``giverny network location``
