
# Changelog

## Unreleased

SECURITY:
FEATURES:
IMPROVEMENTS:
- monetd~babble:   Added Suspended state to babble to allow recovery.
- monetd~evm-lite: Added export endpoint which outputs a genesis.json file for
                   the current state.
- monetd~e2e:      Added "rebuild" e2e test which creates a network from the 
                   state of the current endpoint.
- monetd~evm-lite: Added support for setting the storage values for a contract
                   in the genesis block.        
- monetd~poa:      Removed the solc dependency, moving to precompiled bytecode
- giverny:         Added `giverny parse` command which extracts the peer list
                   from the POA.Storage section of a genesis.json file.                                                
BUG FIXES:

## v0.3.1 (October 25, 2019)

BUG FIXES:

- monetd~babble: Fix error reading events expired from the cache
- monetd~babble: Defer network starting listening after bootstrap process ends
- monetd~babble: Add the badgerdb truncate option to truncate value log files to
                 delete corrupt data
- monetd~config: Prevent ambiguous choice of keys in pull command

## v0.3.0 (October 16, 2019)

IMPROVEMENTS:

- monetd:        Restructure the configuration directories to facilitate role
                 separation between data and config
- monetd~babble: Badger_DB updated to latest v1.6.0
- monetd~poa:    Implement a voting scheme to evict a validator
                 joinleavetest to test nominating and evicting validators
- monetd~cli:    Additional warnings and confirmation prompts about when
                 overwriting configuration files

## v0.2.5 (October 2, 2019)

IMPROVEMENTS:

- docs: Document POA process
- monetd~babble: Add timestamp to node stats

BUG FIXES:

- monetd~babble: Intercept SIGTERM together with SIGINT
- monetd~config: Ignore error when 'config pull' tries to provide default key.

## v0.2.4 (September 18, 2019)

FEATURES:

- monetd~babble: enable advertising a different address than "babble.listen"
- e2e: enhanced testing toolset

IMPROVEMENTS:

- monetd~evm-lite: more granular mutexes around state increases service API
                   throughput.

## v0.2.3 (September 13, 2019)

FEATURES:

- evm-lite~currency: new denominations for token units

IMPROVEMENTS:

- monetd~babble: better handling of "normal" SelfParent errors
- monetd~e2e: use new js libs for currency operations

BUG-FIXES: 

- monetd: Add moniker to configuration
- evm-lite~state: handling transaction promises and errors

## v0.2.2 (September 6, 2019)

FEATURES:

- evm-lite~service: minimum gas price
- evm-lite~state: make use of coinbase address
- monetd~babble-proxy: pseudo-random coinbase selection
- giverny: build test transaction sets

## v0.2.1 (August 29, 2019)

Update glide dependencies with latest versions of Babble and EVM-Lite.

IMPROVEMENTS:

* evm-lite~service: /tx endpoint returns receipt directly.
* evm-lite~service: enable CORS
* babble~service: enable CORS

## v0.2.0 (August 8, 2019)

Refactor with new version of EVM-Lite.

IMPROVEMENTS:

* monetd: Easier configuration
* giverny: Deploy local docker networks
* tests: End to end tests
* docs: More accurate, more comprehensive documentation. 

## v0.1.0 (July 16, 2019)

Tightly integrated packaging of validator node software for Monet Hub.

FEATURES:

- monetd: Server process run by validators.
- giverny: Tool to build and deploy local testnets with docker.
