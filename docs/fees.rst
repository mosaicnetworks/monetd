.. _fees_rst:

Transaction Fees
================

Every operation that modifies the state (transfer, smart-contract creations,
smart-contract call, etc.) carries a cost. Within the EVM, this cost is 
denominated in gas. For example, a simple transfer costs 21000 gas. When users 
create and submit transactions, they can set the maximum amount of gas they want
to spend, and how much om (10^-18 Tenom) they are willing to pay per unit of gas 
consumed. Therefore, if their transaction is applied, it will cost them a 
transaction fee of gas-price * gas-consumed, which is capped by gas-price * 
gas-max.

Transaction fees serve a dual purpose: to incentivise validators, and to prevent
denial of service attacks.

Distribution Among Validators
-----------------------------

Every transaction applied to the EVM is associated with a coinbase address 
(possibly empty), which receives the transaction fee. In monetd, we have 
implemented a system to fairly and securily distribute fees among validators.

Upon committing a Babble block, we fetch the corresponding validator-set from
Babble. Then we use the block hash to obtain a pseudo-random number which we 
use to select a peer from the validator-set. This peer will receive all the 
transaction fees from that block. This system is fair and secure because the 
selection process is evenly distributed and it is impossible for malicious
validators to game it by manipulating the block hash.  

Minimum Gas Price
-----------------

Validators running a monetd node can set a minimum gas price, via the 
eth.min-gas-price configuration flag, to refuse broadcasting transactions with
lower gas-prices. Note that this filtering is done at the service layer, so it
will not have any effect on transactions broadcasted by other nodes. To send a 
transaction via a node, the transaction creator must set the gas price to a 
value greater or equal to that node's minimum gas price.