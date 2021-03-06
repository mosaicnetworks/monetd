.. _install_rst:

Installing monetd
=================

Versioning
++++++++++

``monetd`` versions follow `semantic versioning <https://semver.org>`__. As we
are still in the 0.x range, different versions might contain undocumented
and/or breaking changes. At this stage, the prefered way of installing
``monetd`` is building from source, or downloading binaries directly.

Docker
++++++

Docker images of ``monetd`` are available from the ``mosaicnetworks``
organisation. Use the ``latest`` tag for the latest released version. The
advantage of using Docker containers is that they come packaged with all the
necessary binary files, including solc, and contain an isolated running
environment where ``monetd`` is sure to run.

**Example**: Mount a configuration directory, and run a node from inside a
``monetd`` container.

.. code::

    docker run --rm -v ~/.monet:/.monet mosaicnetworks/monetd run

Downloads
+++++++++

Download the latest version of monetd:

- `Linux <https://dashboard.monet.network/api/downloads/monetd/?os=linux>`__
- `Mac <https://dashboard.monet.network/api/downloads/monetd/?os=mac>`__
- `Windows <https://dashboard.monet.network/api/downloads/monetd/?os=windows>`__



**Example**: Download ``monetd`` and copy it to the local bin directory.

.. code ::

    $ wget -O monetd -L "https://dashboard.monet.network/api/downloads/monetd/?os=linux"

    $ chmod 751 monetd
    $ sudo mv monetd /usr/local/bin/

Please refer to :ref:`monetd systemd<monetd_systemd_rst>` for instructions to
setup a ``systemd`` service on Linux systems.

Building From Source
++++++++++++++++++++

Dependencies
------------

The key components of the Monet Toolchain, which powers the MONET Hub, are
written in `Golang <https://golang.org/>`__. Hence, the first step is to
install **Go version 1.9 or above**, which is both the programming language and
a CLI tool for managing Go code. Go is very opinionated and requires `defining
a workspace <https://golang.org/doc/code.html#Workspaces>`__ where all Go code
resides. The simplest test of a Go installation is:

.. code:: bash

    $ go version

``monetd`` uses `Glide <http://github.com/Masterminds/glide>`__ to manage
dependencies.

.. code::

    $ curl https://glide.sh/get | sh

Solidity Compiler
~~~~~~~~~~~~~~~~~

The Monet Toolchain uses Proof of Authority (PoA) to manage the validator set.
This is implemented using a smart-contract written in
`Solidity <https://solidity.readthedocs.io/en/develop/introduction-to-smart-contracts.html>`__,
with the corresponding EVM bytecode set in the genesis file.

A standard precompiled contract is included in ``monetd`` and ``giverny`` and
will be included by default in the generated ``genesis.json`` file. If you wish
to customise the POA smart contract you will need to have the Solidity compiler
(``solc``) installed. Most users will not need to. If required, please refer to
the `solidity compiler installation instructions <https://solidity.readthedocs.io/en/develop/installing-solidity.html>`__.

Previously the Node.js version of the compiler was not supported for compiling
bytecode. This limitation no longer applies, as solc is no longer embedded in
the apps, so ``npm install solc`` is now valid.

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
