.. _giverny_rst:

#################
Giverny Reference
#################

``giverny`` is the advanced configuration tool for the Monet Toolchain.

The current subcommands are:

- **help** --- help
- **version** --- outputs version information
- **keys** --- key management tools
- **network** --- configure and build networks
- **transactions** --- generate test transactions sets
- **parse** --- parse a genesis file

***********
Global Flag
***********

The ``--verbose`` flag, or ``-v`` for short, turns on extended messages for
each ``giverny`` command.

****
Help
****

``giverny`` has context sensitive help accessed either by running 
``giverny help`` or by adding a ``-h`` parameter to the relevant command.

*******
Version
*******

The ``version`` subcommand outputs the version number for ``monetd``, 
``EVM-Lite``, ``Babble`` and ``Geth``.

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash.

.. include:: _static/includes/giverny_version.txt
    :code: bash

****
Keys
****

The ``keys`` subcommand offers tools to manage keys.

Keys Flags
==========

In addition to the ``--verbose`` flag, the ``keys`` subcommand defines addtional
flags as follows:

.. include:: _static/includes/giverny_keys_flags.txt
    :code: bash

Import
======

The ``import`` subcommand is used to import a pre-existing key pair into the
``monetd`` keystore, assigning the given moniker and setting a passphrase.

.. include:: _static/includes/giverny_help_keys_import.txt
    :code: bash

Generate
========

The ``generate`` subcommand is used to bulk generate key pairs for a test net.
The ``--prefix`` parameter defines a prefix for the account monikers. Then the
``--min-suffix`` and ``--max-suffix`` define the range of accounts names.

E.g. ``--prefix=Acc --min-suffix=1 --max-suffix=3`` would generate accounts:
``Acc1``, ``Acc2`` and ``Acc3``.

.. include:: _static/includes/giverny_help_keys_generate.txt
    :code: bash

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
invokes the build command automatically, unless the ``--no-build`` parameter
is specified.

Syntax
------

.. include:: _static/includes/giverny_help_network_new.txt
    :code: bash

Nodes
-----

The number of nodes in this network is specified by the ``--nodes [int]`` 
parameter. The ``--initial-peers [int]`` parameter specifies the number of
initial peers. If not set it assumes that all nodes are in the initial peer set.

IP Addresses
------------

An initial IP address is supplied using the ``--initial-ip`` parameter. It is
assumed the IP address range will be assigned by simply incrementing the last
octet of the IP address for each node. N.B. the first node will be assigned the
actual IP supplied by the ``initial-ip`` parameter.

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
- ``--no-save-pass`` --- suppresses saving pass phrases in the network
  configuration keystore directory

The typical use case scenarios for these flags would be:

- None specified --- you are prompted to enter the passphrase for each node
  which is saved
- ``--pass`` only --- the specified pass phrase is used, and saved in the
  config folder
- ``--pass`` and ``--no-save-pass`` --- the specified pass phrase is used
  **and** not saved in the config folder
- ``--generate-pass`` only --- pass phrases are generated and saved
- ``--no-save-pass`` only --- you are prompted to enter the passphrase for each
  node, which is not saved in the config folder

Build
-----

By default ``giverny network new`` will run ``giverny network build``
automatically. This can be disabled by specifying the ``-no-build`` flag.

Examples
--------

An example of the new subcommand:

.. code:: bash

    $ giverny network new test11 --names e2e/sampledata/names.txt --nodes 7 --pass e2e/sampledata/pwd.txt --initial-peers 3 --initial-ip 192.168.1.19

Build
=====

The ``giverny network build`` subcommand takes a configuration created by the
``new`` subcommand and builds ``peers.json`` and ``genesis.json`` files.

``build`` can be run repeatably safely. It is envisaged that users will edit
the ``network.toml`` file to adjust token allocations or change addresses.

``--no-generate-keys`` disables the creation of any keys not already in the
keystore.

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

List
====

The ``list`` subcommand lists the configured network names.

.. code:: bash

    $ giverny network list
    benchmark
    benchnet
    bulktransfers


Dump
====

The ``dump`` subcommand outputs the nodes in a named network in bar delimited
format as below:

.. code:: bash

    giverny network dump crowdfundnet
    Amelia|172.77.5.10|0x7bBE1Df184142709d5B99C5788982D0bEE5d1167|true|false
    Becky|172.77.5.11|0xC6a29c6378C20eA9E868EdD3538Ba58d09318f81|true|false
    Chloe|172.77.5.12|0x7b225252dEe5aDa558a233c7B8B654Ef366EBe61|true|false
    Danu|172.77.5.13|0xC0d14Ed110045d7A401ecC9E57628D55e56Fd4c4|true|false

Start
=====

The ``start`` subcommand starts a docker network. Individual nodes are not
started unless the ``--start-nodes`` parameter is specified. If the
``--force-network`` parameter is set, then the network is forced down if it
is already running.


.. include:: _static/includes/giverny_help_network_start.txt
    :code: bash

Stop
====

The ``stop`` subcommand stops a docker network and all the nodes within it.

.. include:: _static/includes/giverny_help_network_stop.txt
    :code: bash

Status
======

The ``status`` subcommand shows the docker network status

.. code:: bash

    $ giverny network status

    Networks

    crowdfundnet   663db79442357cb8814b7ff40076abdd6479a2f5b24ab7087deceaf07913999a  bridge
    none   257f919e7203933bb10aadf17637552b16acb5490b5c8141815e2f19c01ff1fe  null
    bridge   37b969bb113d1707ce01328803bc57d0dc86bb349f617112d242f82ade0ada76  bridge
    host   b89aa9c3a413c14af09cbab6b3ee4450c2cf1cfdbc0449cb28d1f73e4c296d8b  host

    Containers

    /Danu   a67496705d1e5bf6dc0b92a7c4ec69d6c055dc2d08a3193d0e2f5c0fde74564b  Up 10 seconds
    /Chloe   d7da80c6e3c7b92a61975208d72e5d4864d5c5b31bb67859d3c1bbb3feb38b43  Up 11 seconds
    /Becky   dbe47cb5ff517f47f18c4307c09de95ef86fe75279657e342d3bcfee2d6f1a1e  Up 11 seconds
    /Amelia   0f08e59ef698aecc62b2c6945d8c351c7902c90f13e6c4828ef9ab7c9ee27ec3  Up 12 seconds

Push
====

The ``push`` subcommand creates a named node on a built docker network. If the
docker network has not yet been build, there is no need to push the node.

.. include:: _static/includes/giverny_help_network_push.txt
    :code: bash

************
Transactions
************

The transaction commands are used to generate transactions sets for end to end
testing of networks.

Generate
========

The ``generate`` subcommand is used to generate transaction sets from the
``network.toml`` file.

The following flags can be set:

.. code:: bash

      --count int        number of tranactions to generate (default 20)
      --faucet string    faucet account moniker (default "Faucet")
  -h, --help             help for generate
      --ips string       ips.dat file path
  -n, --network string   network name
      --surplus int      additional credit to allocate each account from the faucet above the bare minimum (default 1000000)

Solo
====

The ``solo`` subcommand is used to generate transactions sets from a single
funded account. Look in ``e2e/tools/build-trans.sh`` for an end to end example
using the ``solo`` command.

.. include:: _static/includes/giverny_help_transactions_solo.txt
    :code: bash

*****
parse
*****

The ``giverny parse`` command parses a given genesis file to report the
initial whitelist set and whether the bytecode matches the release bytecode for
your currently installed giverny version.

.. include:: _static/includes/giverny_parse.txt
    :code: bash


An example session is included below:

.. code:: bash

    $ giverny parse $HOME/.giverny/networks/joinleavetest/genesis.json

    POA Address:  0xaabbaabbaabbaabbaabbaabbaabbaabbaabbaabb

    4 peers found

    0xd813b4c2f416bf9cc038b2b3ebbcf8f0bfc6d713  node0
    0x4a47de4f72810f4f002c254d7270877e8f24e145  node1
    0xb582c7d8b6c6f496387eae9386d7e7724d96c61f  node2
    0xd593c5797fbfbcc381c74ad8c9d322d1baa3bc40  node3

    POA bytecode matches the standard contract



