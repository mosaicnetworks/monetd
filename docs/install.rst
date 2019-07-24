.. _install_rst:

Monet Hub Installation
======================


Dependencies
------------

The key components of the Monet Hub are written in
`Golang <https://golang.org/>`__. Hence, the first step is to install **Go 
version 1.9 or above** which is both the programming language and a CLI tool for 
managing Go code. Go is very opinionated and will require you to `define a
workspace <https://golang.org/doc/code.html#Workspaces>`__ where all your go 
code will reside. The simplest test of a go installation is:

.. code:: bash

    $ go version

Solidity Compiler
~~~~~~~~~~~~~~~~~

The Monet Hub uses proof of authority for its validator nodes. This is 
implemented using a smart contract written in
`Solidity <https://solidity.readthedocs.io/en/develop/introduction-to-smart-contracts.html>`__,
with the initial peers set embedded in it, being placed in the genesis block. To 
build the genesis block, at least one of the initial peers will need to have the 
Solidity Compiler solc available to be able to compile the contract into the 
genesis block.

Please refer to the `solidity compiler installation
instructions <https://solidity.readthedocs.io/en/develop/installing-solidity.html>`__.

Other requirements
~~~~~~~~~~~~~~~~~~

Bash scripts used in this project assume the use of GNU versions of coreutils. 
Please ensure you have GNU versions of these programs installed:-

example for macOS:

.. code:: bash

    # --with-default-names makes the `sed` and `awk` commands default to gnu sed and gnu awk respectively.
    brew install gnu-sed gawk --with-default-names

Installation
------------

Clone the `repository <https://github.com/mosaicnetworks/monetd>`__ in the 
appropriate GOPATH subdirectory:

.. code:: bash

    $ mkdir -p $GOPATH/src/github.com/mosaicnetworks/
    $ cd $GOPATH/src/github.com/mosaicnetworks
    [...]/mosaicnetworks$ git clone https://github.com/mosaicnetworks/monetd.git

Monetd uses `Glide <http://github.com/Masterminds/glide>`__ to manage
dependencies.

.. code:: bash

    [...]/babble$ curl https://glide.sh/get | sh
    [...]/babble$ glide install

This will download all dependencies and put them in the **vendor** folder.

Then build and install:

.. code:: bash

    [...]/monetd$ make install

Tests
-----

Use the Go tool to run tests:

.. code:: bash

    [...]/monetd $ make test

Further documentation can be found in the :ref:`test_rst`.

Other Make Commands
-------------------

To build binaries for use in docker:

.. code:: bash

    [...]/monetd$ make docker

