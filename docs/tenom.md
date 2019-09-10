# Tenom

**This is a draft discussion document.**

Internally EVM-Lite balances, values and gas prices are denominated in `Attom`. All user interactions are denoted in Tenom. There are 10<sup>18</sup> attoms in one Tenom.

The Tenom symbol is: Ŧ although for user entry a capital T would usually be used. 

The Tenom symbol is U+0166 in unicode. `&#x166;` is the HTML entity for &#x166;.

In Linux to directly enter a unicode character, hold the left control key and shift, then press u. An underscore u character will appear. Press 0166 then space and the Ŧ character will appear. 

On windows, press and hold ALT and type 0166. 

In GOLANG we can just include the character literal, but `"\u0166"` will also work. 

In JS we can also use `\u0166`.


SI Prefixes:

| Prefix | Opt 1  | Opt 2      |Quickest Symbol |  Quick Symbol | Formal Symbol | Value            |
|--------|--------|------------|--------|---------------|---------------|------------------|
|        | Tenom  | Tenom      | T      | T           | TŦ            | 1                |
| milli  | Millom | Millitenom | m      | mT         | mŦ            | 10<sup>-3</sup>  |
| micro  | Microm | Microtenom | u       | uT       | μŦ            | 10<sup>-6</sup>  |
| nano   | Nanom  | Nanotenom  | n      | nT          | nŦ            | 10<sup>-9</sup>  |
| pico   | Picom  | Picotenom  | p       | pT         | pŦ            | 10<sup>-12</sup> |
| femto  | Femtom | Femtotenom | f       | fT       | fŦ            | 10<sup>-15</sup> |
| atto   | Attom  | Attotenom  | a    | aT         | aŦ            | 10<sup>-18</sup> |


## Value


A Tenom is approximately 1 CHF which to an order of magnitude is equal to 1 GBP, 1 USD and 1 EUR. 

Our initial approximate pricing would be for a transfer to cost 0.1p, i.e. approx 1mŦ, which is 10<sup>15</sup> attoms.

A transfer has approximately 21,000 gas cost. Thus the gas price would be approximately 5 x 10<sup>11</sup>  attoms. 







