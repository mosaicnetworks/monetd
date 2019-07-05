# monetcli network wizard

----

## Table of Contents

+ [Configuration](#configuration)
+ [Main Menu](#main-menu)
+ [Peers](#peers)
+ [Contract](#contract)

----


## Configuration
After completing this section successfully, there will a basic configuration file

---------------------------
|Code|Question|Answers / Outputs|
|------|-------|------------|
|**001**|Location of configuration directory [homeDir] |File Path as a string|
|-|Check if directory specified exists. If directory does not exist run monetcli network new (if fails panic) and jump to **next section**| |
|**002**|Directory exists, force creation of new config, backing up the previous version [Y/n]|If No, exit gracefully|
| -| Run monetcli network new `--force`. If fails panic|Fall through to **next section**|


## Main Menu

---------------------------
|Code|Question|Answers / Outputs|
|------|-------|------------|
|**010**|Peers, Contract, Compile, Publish|Jumps to named section|



## Peers
This section allows the setting of peers.

---------------------------
|Code|Question|Answers / Outputs|
|------|-------|------------|
|**100**|Lists Peers currently configured. Choose and action|[F]inish configuring peers - only available if there is a peer set. Jumps to **next section**.|
|-||[A]dd a peer with an existing public key. Jumps to **110**.
|-||[G]enerate a new keypair and add a peer  **120**.
|-||[E]dit Peer **130**|
|-||[I]nfo Peer **140**|
|**110**|Node Name||
|**111**|IP of Node||
|**112**|Is Validator|Run monetcli network generate, jump to **100**|


## Contract