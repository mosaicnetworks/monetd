.. _install_rst:

Installing Monetd
=================

Versioning
++++++++++

Monetd versions follow `semantic versioning <https://semver.org>`__. As we are
still in the 0.x range, different versions might contain undocumented and/or
breaking changes. At this stage, the prefered way of installing monetd is
building from source, or using our public Docker images.

Docker
++++++

Docker images of monetd are available from the ``mosaicnetworks`` organisation.
Use the ``latest`` tag for the latest released version. The advantage of using
Docker containers is that they come packaged with all the necessary binary
files, including solc, and contain an isolated running environment where monetd
is sure to run.

**Example**: Mount a configuration directory, and run a node from inside a
monetd container.

.. code::

    docker run --rm -v ~/.monet:/.monet mosaicnetworks/monetd run

Downloads
+++++++++

Binary packages of monetd are available at
`<https://github.com/mosaicnetworks/monetd/releases>`__.


Building From Source
++++++++++++++++++++

Dependencies
------------

The key components of the Monet Hub are written in 
`Golang <https://golang.org/>`__. Hence, the first step is to install **Go 
version 1.9 or above**, which is both the programming language and a CLI tool
for managing Go code. Go is very opinionated and requires `defining a
workspace <https://golang.org/doc/code.html#Workspaces>`__ where all Go code 
resides. The simplest test of a Go installation is:

.. code:: bash

    $ go version

Monetd uses `Glide <http://github.com/Masterminds/glide>`__ to manage
dependencies.

.. code::

    $ curl https://glide.sh/get | sh

Solidity Compiler
~~~~~~~~~~~~~~~~~

The Monet Hub uses Proof of Authority (PoA) to manage the validator set. This is 
implemented using a smart-contract written in
`Solidity <https://solidity.readthedocs.io/en/develop/introduction-to-smart-contracts.html>`__,
with the corresponding EVM bytecode set in the genesis file. For every newly 
defined network, the smart-contract needs to be recompiled because it embeds the
initial whitelist. Hence, the Solidity compiler (solc) is a requirement to
define a new network and produce the appropriate genesis file.

Please refer to the `solidity compiler installation
instructions <https://solidity.readthedocs.io/en/develop/installing-solidity.html>`__.

**Attention**: The Node.js version of the compiler is not supported. **Do not
install via** ``npm install solc``.

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

Run the following command to download all dependencies and put them in the 
**vendor** folder.

.. code:: bash

    [...]/monetd$ make vendor

Then build and install:

.. code:: bash

    [...]/monetd$ make install


Tests
-----

``Monetd`` has both unit tests and end to end tests included in the distribution. 

Unit Tests
~~~~~~~~~~

Monetd comes with extensive units tests, both for ``monetd`` itself and for the
``evm-lite`` and ``babble`` components. The tests can be launched thus:

.. code:: bash

    [...]/monetd$ make test

    Monet Tests


    ?       .../monetd/cmd/giverny  [no test files]
    ?       .../monetd/cmd/giverny/commands [no test files]
    ?       .../monetd/cmd/giverny/commands/keys    [no test files]
    ?       .../monetd/cmd/giverny/commands/network [no test files]
    ?       .../monetd/cmd/giverny/commands/server  [no test files]
    ?       .../monetd/cmd/giverny/configuration    [no test files]
    ?       .../monetd/cmd/monetd   [no test files]
    ?       .../monetd/cmd/monetd/commands  [no test files]
    ?       .../monetd/cmd/monetd/commands/config   [no test files]
    ?       .../monetd/cmd/monetd/commands/keys     [no test files]
    ok      .../monetd/src/babble   0.058s
    ?       .../monetd/src/common   [no test files]
    ?       .../monetd/src/config   [no test files]
    ?       .../monetd/src/configuration    [no test files]
    ?       .../monetd/src/contract [no test files]
    ?       .../monetd/src/crypto   [no test files]
    ?       .../monetd/src/docker   [no test files]
    ?       .../monetd/src/files    [no test files]
    ?       .../monetd/src/peers    [no test files]
    ?       .../monetd/src/types    [no test files]
    ?       .../monetd/src/version  [no test files]

    EVM-Lite Tests

    ?       .../vendor/.../evm-lite/src/common      [no test files]
    ?       .../vendor/.../evm-lite/src/config      [no test files]
    ?       .../vendor/.../evm-lite/src/consensus   [no test files]
    ?       .../vendor/.../evm-lite/src/consensus/solo      [no test files]
    ok      .../vendor/.../evm-lite/src/currency    0.005s
    ?       .../vendor/.../evm-lite/src/engine      [no test files]
    ?       .../vendor/.../evm-lite/src/service     [no test files]
    ok      .../vendor/.../evm-lite/src/state       14.536s
    ?       .../vendor/.../evm-lite/src/version     [no test files]

    Babble Tests

    ok      .../vendor/.../babble/src/babble        0.261s
    ok      .../vendor/.../babble/src/common        0.069s
    ?       .../vendor/.../babble/src/config        [no test files]
    ?       .../vendor/.../babble/src/crypto        [no test files]
    ok      .../vendor/.../babble/src/crypto/keys   0.061s
    ok      .../vendor/.../babble/src/hashgraph     11.280s
    ?       .../vendor/.../babble/src/mobile        [no test files]
    ok      .../vendor/.../babble/src/net   0.102s
    ok      .../vendor/.../babble/src/node  41.565s
    ok      .../vendor/.../babble/src/peers 0.112s
    ?       .../vendor/.../babble/src/proxy [no test files]
    ok      .../vendor/.../babble/src/proxy/dummy   0.048s
    ok      .../vendor/.../babble/src/proxy/inmem   0.072s
    ok      .../vendor/.../babble/src/proxy/socket  0.057s
    ?       .../vendor/.../babble/src/proxy/socket/app      [no test files]
    ?       .../vendor/.../babble/src/proxy/socket/babble   [no test files]
    ?       .../vendor/.../babble/src/service       [no test files]
    ?       .../vendor/.../babble/src/version       [no test files]


End to End Tests
~~~~~~~~~~~~~~~~

The end to end tests are hosted in the ``e2e`` subfolder of the monetd 
repository. 

To run all of the tests:

.. code:: bash

    [...]monetd/e2e$ make install
    [...]monetd/e2e$ make tests


The end to end tests require docker to be installed, npm and node, and we 
would recommend installing monetcli too. The tests will take several minutes to 
run. 