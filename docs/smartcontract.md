# Monet Hub POA Smart Contract

This document describes the requirements for a smart contract to implement POA in a Monet hub. 
[comment]: # (//TODO write this document)


----

## Table of Contents

+ [Solidity](#solidity)
    + [Version](#version)
    + [Constructor](#constructor)
    + [Modifier](#modifier)
    + [Functions linked to Babble](#functions-linked-to-babble)
    + [Payable calls](#payable-calls)
    + [Decision Function](#decision-function)
    + [Information Calls](#information-calls)
    + [Events](#events)
+ [Generated Genesis Whitelist Section](#generated-genesis-whitelist-section)

----




## Solidity

### Version

The first line of the contract is a pragma specifying the solidity version required. Currently this is set to greater than or equal to `0.4.22`.
```
pragma solidity >=0.4.22;
```

### Constructor

The contract is embedded in the genesis block. This means that there is no conventional constructor. It is possible to add a hook to payable function calls to set an initial state if it has not already been initialised.  


### Modifier

`checkAuthorisedModifier` is used to restrict access to payable functions. The internals of that function could be ameneded to your new scheme.

### Functions linked to Babble 

Babble calls one of the following functions to verify whether a peer making a join request as access. Any replacement smark contract aill need to implement these functions. 

```
   function checkAuthorised(address _address) public view returns (bool)
   function checkAuthorisedPublicKey(bytes32  _publicKey) public view returns (bool)
```


### Payable calls


Functions that the client tools expect to be present. 

```
function submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)
function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)
```

### Decision Function

This function decides when a vote is complete. Currently it requires all people on the whitelist to approve. It is anticipated that some form of majority voting would be implemented to prevent paralysis if a peer drops out. 

```
function checkForNomineeVoteDecision(address _nomineeAddress) private returns (bool decided, bool voteresult)
```

### Information Calls

The following information calls are available:

```
function getNomineeElection(address _address) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
function getNomineeCount() public view returns (uint count)
function getNomineeAddressFromIdx(uint idx) public view returns (address NomineeAddress)
function getNomineeElectionFromIdx(uint idx) public view returns (address nominee, address proposer, uint yesVotes, uint noVotes)
function getWhiteListCount() public view returns (uint count)
function getWhiteListAddressFromIdx(uint idx) public view returns (address WhiteListAddress)
function getYesVoteCount(address _nomineeAddress)  public view returns (uint count)
function getYesVoterFromIdx(address _nomineeAddress, uint _idx)  public view returns (address voter)
function getNoVoteCount(address _nomineeAddress)  public view returns (uint count)
function getNoVoterFromIdx(address _nomineeAddress, uint _idx) public view returns (address voter)
function getMoniker(address _address) public view returns (bytes32 moniker)
```

### Events

The following events are emitted by the smart contract. It is envisaged that the same events would be emitted by any replacement contract. 

```
/// @notice Event emitted when the vote was reached a decision
/// @param _nominee The address of the nominee
/// @param _yesVotes The total number of yes votes cast for the nominee to date
/// @param _noVotes The total number of no votes cast for the nominee to date
/// @param _accepted The decision, true for added to the whitelist, false for rejected
    event NomineeDecision(
        address indexed _nominee,
        uint  _yesVotes,
        uint _noVotes,
        bool indexed _accepted
    );
```

```
/// @notice Event emitted when a nominee vote is cast
/// @param _nominee The address of the nominee
/// @param _voter The address of the person who cast the vote
/// @param _yesVotes The total number of yes votes cast for the nominee to date
/// @param _noVotes The total number of no votes cast for the nominee to date
/// @param _accepted The vote, true for accept, false for rejected
    event NomineeVoteCast(
        address indexed _nominee,
        address indexed _voter,
        uint  _yesVotes,
        uint _noVotes,
        bool indexed _accepted
    );
```



```
/// @notice Event emitted when a nominee is proposed
/// @param _nominee The address of the nominee
/// @param _proposer The address of the person who proposed the nominee
    event NomineeProposed(
        address indexed _nominee,
        address indexed _proposer
    );
```

```
/// @notice Event emitted to announce a moniker
/// @param _address The address of the user
/// @param _moniker The moniker of the user
    event MonikerAnnounce(
        address indexed _address,
        bytes32 indexed _moniker
    );
```


## Generated Genesis Whitelist Section



