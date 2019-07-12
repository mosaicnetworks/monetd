Examples
========

This document is an end to end example, creating a test net of three
nodes and demonstrating basic functionality within that test net.

Software Installation
---------------------

On each of the three nodes, you need to :ref:`install_rst` .
You also will need to install :ref:`clients_rst_evmlc` .

Configuration
-------------

node1
~~~~~

We need to configure and set up 3 nodes, which we will call ``node1``,
``node2``, and, ``node3``. On ``node1`` we will run ``monetcfgsrv``. In
a new terminal on ``node1`` only type:

.. code:: bash

   $ monetcfgsrv
   Starting monetcfgsrv
   192.168.1.18:8088

Node the address given (in this case: ``192.168.1.18:8088``), we will
need that for configuring the other nodes.

Leaving the ``monetcfgsrv`` window open, open a new terminal session on
``node1`` and type:

.. code:: bash

   $ monetcli testnet

if you get a message like this:

::

   This is a destructive operation. Remove/rename the following folder to proceed.
   /home/jon/.monetcli/testnet

Then you need to rename the given folder like this, and run
``monetcli testnet`` again:

.. code:: bash

   $ mv /home/jon/.monetcli/testnet /home/jon/.monetcli/testnet.~1~

.. code:: bash

   /\ \__                   /\ \__                    /\ \__
   \ \ ,_\     __     ____  \ \ ,_\    ___       __   \ \ ,_\
    \ \ \/   /'__`\  /',__\  \ \ \/  /' _ `\   /'__`\  \ \ \/
     \ \ \_ /\  __/ /\__, `\  \ \ \_ /\ \/\ \ /\  __/   \ \ \_
      \ \__\\ \____\\/\____/   \ \__\\ \_\ \_\\ \____\   \ \__\
       \/__/ \/____/ \/___/     \/__/ \/_/\/_/ \/____/    \/__/



   The configuration server is a running instance of monetcfgsrv, which should be run by one of the initial peers. If you are running it, you can use the localhost default address, otherwise you need to ask the person running it for their IP address.
   ✔ Configuration Server:  : |http://localhost:8088

You are asked to specify the configuration server. Whilst the localhost
default will work for this node, you will need to change it to the value
reported by monetcfgsrv. In this case ``http://192.168.1.18:8088``.

Press ``Enter`` to enter your value.

Next you enter your moniker - a more user-friendly name than an address
or public key. We are going to use node1.

.. code:: bash

   ✗ Enter your moniker:   : |

Next up you enter your IP address. The default value should be fine in
nearly cases.

.. code:: bash

   ✔ Enter your ip without the port:   : |192.168.1.18

Then you enter you passphrase used to secure your keys. You need to
re-enter to confirm that it was entered correctly.

.. code:: bash

   Enter Keystore Password:   : ######|
   ✗ Confirm Keystore Password:   : |

It then returns your generated address. And offers you some publishing
options. At this point you just leave ``node1`` and move to ``node2``.

.. code:: bash

   Address: 0x8141948ffAE77ce18D328c930E857DA1ba4c4A65
   Choose publish to build the configuration files.
   Choose check to see if another peer has built them and if so, use them.
   Use the arrow keys to navigate: ↓ ↑ → ←
   ? Choose your action  :
     ▸ Check if published
       Publish, no more initial peers will be allowed to be added
       Exit

node2 and node3
~~~~~~~~~~~~~~~

On ``node2`` and ``node3`` in turn, peform the actions in this
subsection:

::

   $ monetcli testnet

   /\ \__                   /\ \__                    /\ \__
   \ \ ,_\     __     ____  \ \ ,_\    ___       __   \ \ ,_\
    \ \ \/   /'__`\  /',__\  \ \ \/  /' _ `\   /'__`\  \ \ \/
     \ \ \_ /\  __/ /\__, `\  \ \ \_ /\ \/\ \ /\  __/   \ \ \_
      \ \__\\ \____\\/\____/   \ \__\\ \_\ \_\\ \____\   \ \__\
       \/__/ \/____/ \/___/     \/__/ \/_/\/_/ \/____/    \/__/



   The configuration server is a running instance of monetcfgsrv, which should be run by one of the initial peers. If you are running it, you can use the localhost default address, otherwise you need to ask the person running it for their IP address.
   ✔ Configuration Server:  : |http://localhost:8088

You are asked to specify the configuration server. You will need to
change it to the value reported by ``monetcfgsrv`` on ``node1``. In this
case ``http://192.168.1.18:8088``.

Press ``Enter`` to enter your value.

Next you enter your moniker - a more user-friendly name than an address
or public key. We are going to use node1.

.. code:: bash

   ✗ Enter your moniker:   : |

Next up you enter your IP address. The default value should be fine in
nearly cases.

.. code:: bash

   ✔ Enter your ip without the port:   : |192.168.1.18

Then you enter you passphrase used to secure your keys. You need to
re-enter to confirm that it was entered correctly.

.. code:: bash

   Enter Keystore Password:   : ######|
   ✗ Confirm Keystore Password:   : |

It then returns your generated address. And offers you some publishing
options. At this point leave monetcli running and make sure we have
reache this stage for all 3 nodes.

::

   Address: 0xc930E857DA1ba4c4A658141948ffAE77ce18D328
   Choose publish to build the configuration files.
   Choose check to see if another peer has built them and if so, use them.
   Use the arrow keys to navigate: ↓ ↑ → ←
   ? Choose your action  :
     ▸ Check if published
       Publish, no more initial peers will be allowed to be added
       Exit

Publishing
----------

For ``node1`` only select:
``Publish, no more initial peers will be allowed to be added`` by
highlighting it and pressing ``Enter``

::

   ✔ Publish, no more initial peers will be allowed to be added
   Getting peers.json
   Unmarshalling peers.json
   Peers list unmarshalled:  1 [0xc0000ef4a0]
   Adding...  node1
   Publish result: true

   Configuration has been published.
   Downloaded peersjson
   Downloaded genesisjson
   Enter your ip without the port:   : |192.168.1.18

You will need to enter the IP of this device. It should default to the
correct value. You will also need to confirm overwriting the Monet
configuration bu selecting ``Yes``.

::

   All files downloaded
   ✔ Yes
   Renaming /home/jon/.monet to /home/jon/.monet.~1~
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

The program will exit.

Next on ``node2`` and ``node3`` select ``Check if published``.

On each node in turn run:

.. code:: bash

   $ monetd run

And leave the windows open.

Watch Script
------------

//TODO fire up a watch script to show the nodes are up and running.

EVM-Lite CLI
------------

The :ref:`evmlc_getting_started_rst` document for
``evm-lite-cli`` will demonstrate how to manage accounts and make
transfers.

You can also view the :ref:`evmlc_poa_rst` document to find an example
on how to nominate a new validator.

//TODO Fire up the Wallet

//TODO Fire up the Dashboard
