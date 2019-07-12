# Proof of Authority Commands

We explain how to use `evmlc`. We will walk through nominating a new `address`, voting for that address and other PoA commands.

## 1. Run `evmlc` in interactive mode

```bash
$ evmlc i

  _____  __     __  __  __           _       _   _               ____   _       ___
 | ____| \ \   / / |  \/  |         | |     (_) | |_    ___     / ___| | |     |_ _|
 |  _|    \ \ / /  | |\/| |  _____  | |     | | | __|  / _ \   | |     | |      | |
 | |___    \ V /   | |  | | |_____| | |___  | | | |_  |  __/   | |___  | |___   | |
 |_____|    \_/    |_|  |_|         |_____| |_|  \__|  \___|    \____| |_____| |___|

 Mode:        Interactive
 Data Dir:    /Users/danu/.evmlc
 Config File: /Users/danu/.evmlc/config.toml
 Keystore:    /Users/danu/.evmlc/keystore


evmlc$
```

## 2. List accounts and create an account to nominate

We will need to create an account to nominate as a validator for the network. Firstly we can view our accounts by running `accounts list -f` (`-f` specified formatted output).

```bash
evmlc$ accounts list -f

.-----------------------------------------------------------------------------.
|                  Address                   |        Balance         | Nonce |
|--------------------------------------------|------------------------|-------|
| 0x702B0ad02a7a6056EB16A697A96d849c228F5fB4 | 1337000000000000000000 |     0 |
'-----------------------------------------------------------------------------'
```

Now we will need to create an account. We can do this by running `accounts create` and following the prompts on the screen.

```bash
evmlc$ accounts create

? Passphrase:  [hidden]
? Re-enter passphrase:  [hidden]

{"version":3,"id":"1153fee6-79e7-46d3-a3cd-93cb86dd71f5","address":"221eff07bd1bf1e1fe21a069523413218c32be42","crypto":{"ciphertext":"a672a0c40304717ac36fab3d69f3e07d7703
6f7cb0669ae299f58a47b3af9efc","cipherparams":{"iv":"d818e6da4347a5f8a8bf8a9ef7bc36b3"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"0d50a248ac6ad
f0933af1d1a4316d1339cd042f489f63ed5218589f0b4963618","n":8192,"r":8,"p":1},"mac":"84168ff91a8191f37c738e93d8bec07226eccf2e1928e544cb2b35797d6ea125"}}
```

## 3. List whitelist

As a sanity check, we will need to see the entire whitelist. We can do this by running `poa whitelist -f`.

```bash
evmlc$ poa whitelist -f

.----------------------------------------------------------------------.
|         Moniker         |                  Address                   |
|-------------------------|--------------------------------------------|
|         node0           | 0x702b0ad02a7a6056eb16a697a96d849c228f5fb4 |
'----------------------------------------------------------------------'
```

## 4. Nominate a new node

We will nominate the created account `221eff07bd1bf1e1fe21a069523413218c32be42` to go through election using the command `poa nominate`.

```bash
evmlc$ poa nominate

? From:  702b0ad02a7a6056eb16a697a96d849c228f5fb4
? Passphrase:  [hidden]
? Nominee:  0x221eFf07BD1bF1e1FE21A069523413218c32bE42
? Moniker:  node1

You (0x702b0ad02a7a6056eb16a697a96d849c228f5fb4) nominated 'node1' (0x221eff07bd1bf1e1fe21a069523413218c32be42)
```

## 5. List nominees

Now that we have nominated an address we can view the nominee list by running `poa nominee list -f`

```bash
evmlc$ poa nominee list -f

.------------------------------------------------------------------------------.
| Moniker |                  Address                   | Up Votes | Down Votes |
|---------|--------------------------------------------|----------|------------|
| Node1   | 0x221eff07bd1bf1e1fe21a069523413218c32be42 |        0 |          0 |
'------------------------------------------------------------------------------'
```

## 6. Vote for the nominee

We can now vote for the nominee by running `poa vote` and following the on-screen prompts.

```bash
evmlc$ poa vote

? From:  702b0ad02a7a6056eb16a697a96d849c228f5fb4
? Passphrase:  [hidden]
? Nominee:  0x221eff07bd1bf1e1fe21a069523413218c32be42
? Verdict:  Yes

You (0x702b0ad02a7a6056eb16a697a96d849c228f5fb4) voted 'Yes' for '0x221eff07bd1bf1e1fe21a069523413218c32be42'.
Election completed with the nominee being 'Accepted'.
```

Since we were the only whitelisted address, the only vote a nominee needs to get whitelisted is ours.

## 7. Check whitelist

We now check the updated whitelist to see if the nominee was officially accepted.

```bash
evmlc$ poa whitelist -f
.----------------------------------------------------------------------.
|         Moniker         |                  Address                   |
|-------------------------|--------------------------------------------|
| Node0                   | 0x702b0ad02a7a6056eb16a697a96d849c228f5fb4 |
| Node1                   | 0x221eff07bd1bf1e1fe21a069523413218c32be42 |
'----------------------------------------------------------------------'
```
