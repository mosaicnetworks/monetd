# Monet-CLI
## Monet Hub tools

## USAGE


```
Monet-CLI

Usage:
  monetcli [command]

Available Commands:
  help        Help about any command
  keys        An Ethereum key manager
  version     Show version info

Flags:
  -h, --help   help for monetcli

Use "monetcli [command] --help" for more information about a command.
```

The keys subcommand is used to manage ethereum keys.

```bash
An Ethereum key manager

Usage:
  monetcli keys [command]

Available Commands:
  generate    Generate a new keyfile
  inspect     Inspect a keyfile
  update      change the passphrase on a keyfile

Flags:
  -h, --help              help for keys
      --json              output JSON instead of human-readable format
      --passfile string   the file that contains the passphrase for the keyfile

Use "monetcli keys [command] --help" for more information about a command.
```
