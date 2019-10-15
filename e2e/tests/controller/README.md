# Controller Test

There is a smart contract that controls the location of the POA smart contract.
The default locations are set in ``src/configuration/const``  ✔

```go
	DefaultContractAddress           = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
	DefaultControllerContractAddress = "aabbaabbaabbaabbaabbaabbaabbaabbaabbaabb"
```

The controller contract and the initial POA contract are included in the
``genesis.json`` file. ✔

This test initialises a network of 3 peers and has a 4th peer join, as per the
``jointest`` test. There are 5 peers in all, ``node0`` to ``node4``. ✔

We then take the existing contract, and create a new contract using the
existing peers and deploy it. 

We then update the controller contract to point to the new contract. 


## Dependency

**N.B this test uses the node js from the ``jointest`` test.**

It is assumed that the following commands are available:

+ jq
+ monetcli


monetd config contract $(monetcli poa whitelist 2> /dev/null | jq -r '.[] | "\(.moniker)=\(.address)"')

