.. _readme_rst:

The Monet Hub
=============

.. figure:: assets/monet_logo.png
   :alt: Monet Logo

   Monet Logo


The monetd respository contains the tools necessary to run and maintain
a validator hub in a Monet network.

They naturally divide into 2 sections: 

+ :ref:`monetcli_rst` -- the swiss army knife of utilities 
+ :ref:`monetd_rst` -- the hub server process

Full details can found at the links above, but the Quick Start section
below may help you where to look.

Quick Start
===========

Installation
------------

The installation process is covered in the installation
instructions :ref:`install_rst` .

--------------

Creating a new Test Net
-----------------------

To set up a new testnet with yourself as one of the initial peers use:

.. code:: bash

    $ monetcli testnet

See the testnet section :ref:`monetcli_testnet` for
more information.

N.B. You will need access to a running ``monetcfgsrv`` instance as
described in the testnet section and the linked document.

--------------

Joining an existing Test Net
----------------------------

To join an existing testnet use:

.. code:: bash

    $ monetcli testjoin

See the testjoin section in :ref:`monetcli_rst` for more
information.

--------------

Interactive Configuration
-------------------------

The general purpose guided configuration can be accessed via:

.. code:: bash

    $ monetcli wizard

This tool is used for more flexible configurations than the tools above,
which offer less options and are thus more straightforward. See the
wizard section in :ref:`monetcli_rst` for more information.

--------------

Clients
-------

Clients and wallets configured to be used with the monet hub are
described in :ref:`clients_rst`.

