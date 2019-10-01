.. _fees_rst:

Transaction Fees
================

Every operation that modifies the state (transfer, smart-contract creation,
smart-contract call, etc.) carries a cost. Within the EVM, this cost is 
denominated in gas. For example, a simple transfer costs 21000 gas. When users 
create and submit transactions, they can set the maximum amount of gas they want
to spend, and how many ``Attoms`` (10^-18 Tenom) they are willing to pay per 
unit of gas consumed. Therefore, if their transaction is applied, it will cost
them a transaction fee of ``gas-price * gas-consumed``, which is capped by 
``gas-price * gas-max``.

Transaction fees serve a dual purpose: to incentivise validators, and to prevent
denial of service attacks.

Distribution Among Validators
-----------------------------

Every transaction applied to the EVM is associated with a coinbase address 
(possibly empty), which receives the transaction fee. In monetd, we have 
implemented a system that fairly and securily distributes fees among validators.

Upon committing a Babble block, we fetch the corresponding validator-set from
Babble. Then we use the block hash to obtain a pseudo-random number which we 
use to select a peer from the validator-set. This peer will receive all the 
transaction fees from that block. This system is fair and secure because the 
selection process is evenly distributed and it is impossible for malicious
validators to game it by manipulating the block hash.  

Minimum Gas Price
-----------------

Validators running a monetd node can set a minimum gas price, via the 
``eth.min-gas-price`` configuration flag, to refuse broadcasting transactions 
with lower gas-prices. To send a transaction via a node, the transaction creator 
must set the gas price to a value greater or equal to that node's minimum gas
price. Note that this filtering is done at the service layer, so it will not
prevent other nodes from including cheaper transactions. 

Tenom Denominations
-------------------

Internally EVM-Lite balances, values and gas prices are denominated in 
``Attom``. All user interactions are denoted in Tenom. There are ``10^18`` 
attoms in one Tenom.

The Tenom symbol is: Ŧ although for user entry a capital T would usually be used. 

SI Prefixes:

.. code:: bash

    | Prefix | Opt 1  | Opt 2      |Quickest Symbol |  Quick Symbol | Formal Symbol | Value            |
    |--------|--------|------------|----------------|---------------|---------------|------------------|
    |        | Tenom  | Tenom      | T              | T             | Ŧ             | 1                |
    | milli  | Millom | Millitenom | m              | mT            | mŦ            | 10^-3            |
    | micro  | Microm | Microtenom | u              | uT            | μŦ            | 10^-6            |
    | nano   | Nanom  | Nanotenom  | n              | nT            | nŦ            | 10^-9            |
    | pico   | Picom  | Picotenom  | p              | pT            | pŦ            | 10^-12           |
    | femto  | Femtom | Femtotenom | f              | fT            | fŦ            | 10^-15           |
    | atto   | Attom  | Attotenom  | a              | aT            | aŦ            | 10^-18           |


The Tenom symbol is U+0166 in unicode. `&#x166;` is the HTML entity for &#x166;.

In Linux to directly enter a unicode character, hold the left control key and 
shift, then press u. An underscore u character will appear. Press 0166 then
space and the Ŧ character will appear. 

On windows, press and hold ALT and type 0166. 

In GOLANG we can just include the character literal, but `"\u0166"` will also 
work. 

In JS we can also use `\u0166`.