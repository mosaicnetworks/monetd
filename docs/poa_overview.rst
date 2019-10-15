.. _poa_overview_rst:

POA Overview
============

As mentionned in the :ref:`design document<design_rst>`, the MONET Hub relies on
a Proof of Authority (POA) system to control the validator-set of the underlying
consensus system (Babble). Here we explain how this system is implemented in 
monetd.

Whitelist
---------

The whitelist is a list of members who are allowed to run a validator node in
the network. Upon processing an ``Internal Transaction`` from Babble (related to
a join request), monetd queries the whitelist and accepts the request if and
only if the corresponding address is whitelisted. For more information about
``Internal Transactions`` and Babble's membership protocol please refer to 
`Babble docs <https//docs.babble.io/en/latest/dynamic_membership.html>`__.


POA Smart Contract
------------------

Technically, the whitelist is implemented in the 
:ref:`POA smart-contract<poa_smartcontract_rst>`, which is defined in the
genesis file, and hence automatically provisioned upon starting a node, along
with the pre-funded genesis accounts. 


Adding Addresses to the Whitelist
---------------------------------

It is possible to pre-populate the whitelist (cf. :ref:`giverny<giverny_rst>`)
to include addresses at genesis. Once the network is up an running, it is also
possible to add new addresses on the fly via a voting mechanism.

To add an address to the whitelist, someone **already on the whitelist** has to
**nominate** the address, and **everyone** on the whitelist needs to cast a 
**yes** vote on that nomination. When the nomination collects yes votes from 
everyone, the address is automatically moved from the nominee list to the 
whitelist. This logic is implemented in the POA smart-contract, and can be
called using ``monetcli``.

``monetcli whitelist`` displays the current whitelist.

``monetcli nominee list`` displays the current nomineelist.

``monetcli nominee vote`` is used to cast a vote for a nomination.



