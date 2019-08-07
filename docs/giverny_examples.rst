.. _giverny_examples_rst:

################
Giverny Examples
################

For reference, the options for ``giverny network new``:

.. code:: bash

    $ giverny network new -h

    giverny network build

    Usage:
    giverny network new [network_name] [flags]

    Flags:
        --generate-pass       generate pass phrases
    -h, --help                help for new
        --initial-ip string   initial IP address of range
        --initial-peers int   number of initial peers
        --names string        filename of a file containing a list of node monikers
        --no-build            disables the automatic build of a new network
        --pass string         filename of a file containing a passphrase
        --save-pass           save pass phrase entered on command line

    Global Flags:
    -g, --giverny-data-dir string   Top-level giverny directory for configuration and data (default "~/.giverny")
    -m, --monet-data-dir string     Top-level monetd directory for configuration and data (default "~/.monet")
    -n, --nodes int                 number of nodes in this configuration (default 4)
    -v, --verbose                   verbose messages



*************************
Development Test Networks
*************************

To make commands repeatable, and to reflect code changes, the following
commands can be prefixed to all the commands below:

``make installgiv; rm -rf ~/.giverny/networks/test9;``

The command above rebuilds the ``giverny`` app and removes the network
``test9``allow the ``new`` commands to be run repeatedly. If you do not remove
the previous network ``test9`` before running ``giverny network new`` then the
command aborts. The ``make installgiv`` is only required if you are making code
changes.

Adding ``-v`` or ``--verbose`` to each of these commands gives addition
information and progress messages in the command output.

New
===

8 node network, 4 initial peers, named from prebaked list of names, generated
passphrases.

.. code:: bash

    giverny network new test9 --generate-pass  --names sampledata/names.txt --nodes 8 --initial-peers 4  -v


3 node network with named nodes, 2 initial peers. Passphrased prompted for on
the command line and used for all key files.

.. code:: bash

    make installgiv; rm -rf ~/.giverny/networks/test9; giverny network new test9 --save-pass  --names sampledata/withnodes.txt --nodes 3 --initial-peers 2  -v

The withnodes.txt file is interesting as it shows the expanded syntax:

.. code:: text

    Jon,192.168.1.18,1T,true
    Martin,192.168.1.3,1G,true
    Kevin,192.168.1.16,1M,false



Export Network
==============

The export command writes the configuration of one or more nodes to a zip file.

To export the configuration of all nodes in a network, type this:

.. code:: bash

    $ giverny network export test9


Take a look in ``~/.giverny/exports``. [#]_ There should be numerous files
named ``test9_[node].zip``. These can be applied to ``monetd`` directly on the
same instance by:


.. code:: bash

    $ giverny network import test9 Danu --from-exports

Alternatively you can use a secondary channel such as slack or e-mail to send
that zip file and then load it --- without changing the name of the file:

.. code:: bash

    $ giverny network import test9 Danu --dir ~/Downloads


Or you can use giverny server and pull it directly. Assuming that you have run
``giverny server start`` on the instance you ran the exports you can:

.. code:: bash

    $ giverny network import test9 Danu --server 192.168.1.4



.. [#] This location is for Linux instances. Mac and Windows uses a different
       path. The path for your instance can be ascertain with this command:
       ``giverny network location``
