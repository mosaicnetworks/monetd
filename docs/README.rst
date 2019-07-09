.. role:: raw-html-m2r(raw)
   :format: html


The Monet Hub
=============


.. image:: assets/monet_logo.png
   :target: assets/monet_logo.png
   :alt: Monet Logo
 

----

Table of Contents
-----------------


* `Quick Start <#quick-start>`_

  * `Installation <#installation>`_
  * `Interactive Configuration <#interactive-configuration>`_
  * `Creating a new Test Net <#creating-a-new-test-net>`_
  * `Joining an existing Test Net <#joining-an-existing-test-net>`_
  * `Clients <#clients>`_

* `Contents of the docs folder <#contents-of-the-docs-folder>`_

----

The monetd respository contains the tools necessary to run and maintain a validator hub in a Monet network. 

They naturally divide into 2 sections:


* `MonetCLI <monetcli.md>`_ -- the swiss army knife of utilities
* `Testnet Docs <monetd.md>`_ -- the hub server process

Full details can found at the links above, but the Quick Start section below may help you where to look. 

Quick Start
===========

Installation
------------

The installation process is covered in `here <install.md>`_.

----

Interactive Configuration
-------------------------

The general purpose guided configuration can be accessed via:

.. code-block:: bash

   $ monetcli wizard

See the wizard section in `Monet CLI docs <monetcli.md>`_ for more information.  

----

Creating a new Test Net
-----------------------

To set up a new testnet with yourself as one of the initial peers use:

.. code-block:: bash

   $ monetcli testnet

See the testnet section `Monet CLI docs <monetcli.md>`_ for more information.  

----

Joining an existing Test Net
----------------------------

To join an existing testnet use:

.. code-block:: bash

   $ monetcli testjoin

See the testjoin section `Monet CLI docs <monetcli.md>`_ for more information.  

----

Monet
-----

To join an existing testnet use:

.. code-block:: bash

   $ monetcli testjoin

See the testjoin section `Monet CLI docs <monetcli.md>`_ for more information.  

----

Clients
-------

Clients and wallets configured to be used with the monet hub are described `here <clients.md>`_.

----

Contents of the docs folder
===========================

.. code-block::

   ├── install.md               - installation instructions
   ├── monetcli.md              - monetcli command documentation
   ├── monetd.md                - monetcfg command documentation
   ├── network.md               - monetcli network command docs, linked from monetcli.md
   ├── README.md                - this document
   ├── smartcontract.md         - requirements for poa smart contract for monet
   ├── testnet.md               - monetcli testnet command docs, linked from monetcli.md
   ├── wizard.md                - monetcli wizard command docs, linked from monetcli.md
   └── archive                  - deprecated docs, scheduled to be removed

----

:raw-html-m2r:`<sup>[Documents Index](README.md) | [GitHub repo](https://github.com/mosaicnetworks/monetd) | [Monet](https://monet.network/) | [Mosaic Networks](https://www.babble.io/)</sup>`
