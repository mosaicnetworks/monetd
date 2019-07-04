# Tests

The Monet hub has extensive unit-testing. Make test uses the Go toolset to run tests. A sample (passing) session is included here. 

```bash
[...]monetd $ make test

Monet Tests


?   	.../monetd/cmd/monetcfgsrv	[no test files]
?   	.../monetd/cmd/monetcli	[no test files]
?   	.../monetd/cmd/monetcli/commands	[no test files]
?   	.../monetd/cmd/monetcli/commands/config	[no test files]
?   	.../monetd/cmd/monetcli/commands/keys	[no test files]
?   	.../monetd/cmd/monetcli/commands/network	[no test files]
?   	.../monetd/cmd/monetcli/commands/testnet	[no test files]
?   	.../monetd/cmd/monetd	[no test files]
?   	.../monetd/cmd/monetd/commands	[no test files]
?   	.../monetd/src/common	[no test files]
?   	.../monetd/src/version	[no test files]

EVM-Lite Tests


?   	.../vendor/.../evm-lite/src/common	[no test files]
?   	.../vendor/.../evm-lite/src/config	[no test files]
?   	.../vendor/.../evm-lite/src/consensus	[no test files]
ok  	.../vendor/.../evm-lite/src/consensus/babble	0.081s
?   	.../vendor/.../evm-lite/src/consensus/raft	[no test files]
?   	.../vendor/.../evm-lite/src/consensus/solo	[no test files]
?   	.../vendor/.../evm-lite/src/engine	[no test files]
?   	.../vendor/.../evm-lite/src/service	[no test files]
?   	.../vendor/.../evm-lite/src/service/templates	[no test files]
ok  	.../vendor/.../evm-lite/src/state	5.128s
?   	.../vendor/.../evm-lite/src/version	[no test files]

Babble Tests


ok  	.../vendor/.../babble/src/babble	0.169s
ok  	.../vendor/.../babble/src/common	0.013s
?   	.../vendor/.../babble/src/crypto	[no test files]
ok  	.../vendor/.../babble/src/crypto/keys	0.125s
ok  	.../vendor/.../babble/src/hashgraph	11.113s
?   	.../vendor/.../babble/src/mobile	[no test files]
ok  	.../vendor/.../babble/src/net	0.076s
ok  	.../vendor/.../babble/src/node	41.183s
ok  	.../vendor/.../babble/src/peers	0.067s
?   	.../vendor/.../babble/src/proxy	[no test files]
ok  	.../vendor/.../babble/src/proxy/dummy	0.080s
ok  	.../vendor/.../babble/src/proxy/inmem	0.108s
ok  	.../vendor/.../babble/src/proxy/socket	0.047s
?   	.../vendor/.../babble/src/proxy/socket/app	[no test files]
?   	.../vendor/.../babble/src/proxy/socket/babble	[no test files]
?   	.../vendor/.../babble/src/service	[no test files]
?   	.../vendor/.../babble/src/version	[no test files]
```

