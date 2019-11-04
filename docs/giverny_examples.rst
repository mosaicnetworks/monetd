.. _giverny_examples_rst:

################
Giverny Examples
################

For reference, the options for ``giverny network new``:

.. include:: _static/includes/giverny_help_network_new.txt
    :code: bash


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

    giverny network new test9 --generate-pass  --names e2e/sampledata/names.txt --nodes 8 --initial-peers 4  -v


3 node network with named nodes, 2 initial peers. Passphrased prompted for on
the command line and used for all key files.

.. code:: bash

    make installgiv; rm -rf ~/.giverny/networks/test9; giverny network new test9 --save-pass  --names e2e/sampledata/withnodes.txt --nodes 3 --initial-peers 2  -v

The withnodes.txt file is interesting as it shows the expanded syntax:

.. code:: text

    Jon,192.168.1.18,1T,true
    Martin,192.168.1.3,1G,true
    Kevin,192.168.1.16,1M,false