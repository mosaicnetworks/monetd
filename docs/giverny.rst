.. _giverny_rst:

Giverny
=======

``giverny`` is the advanced configuration tool. 

The current subcommands are:

- **help** --- help
- **version** --- outputs version information
- **keys** --- key management tools

**This documentation is being lovingly crafted now, and will be ready very soon.**

Help
----

``giverny`` has context sensitive help accessed either by
running ``giverny help`` or by adding a ``-h`` parameter to the relevant
command. 



Version
-------

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



Keys
----

The keys sub-command offers tools to manage keys. 

Import
~~~~~~

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


Server
------

The ``server`` sub-command is used for admninstering a REST server used to co-ordinate 
configurations between multiple nodes prior to the initial node of a network. 

The server listens on port 8088. It writes logs to ``.monet/giverny/server.log``.

For instructions on how to use, see the recipes for setting up networks. 

Start
~~~~~


To start the server in the foreground:

.. code:: bash

    $ giverny server start --background





To start the server in the background:

.. code:: bash

    $ giverny server start --background




Stop
~~~~

To stop a server running in the background: 

.. code:: bash

    $ giverny server stop


Status
~~~~~~

Reports on the status of the server. It both checks for the PID file in 
``.monet/giverny/server.pid`` and checks the the server is responding on 
localhost:8088. 

.. code:: bash

    $ giverny server status

