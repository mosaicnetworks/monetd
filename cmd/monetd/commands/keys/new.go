package keys

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/spf13/cobra"
)

// newNewCmd returns the command that creates a new keypair
func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [moniker]",
		Short: "create a new keyfile",
		Long: `
This command generates a new cryptographic key-pair, and produces two files:

- [datadir]/keystore/[moniker].json : The encrypted keyfile
- [datadir]/keystore/[moniker].toml : Key metadata

[moniker] is a friendly name, which can be reused in other commands to refer to 
the key without having to type or copy a long string of characters.

If the --passfile flag is not specified, the user will be prompted to enter the
passphrase manually. Otherwise, it will be read from the file pointed to by
--passfile.
`,
		Args: cobra.ExactArgs(1),
		RunE: newKey,
	}

	return cmd
}

func newKey(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	// key is returned, but we don't want to do anything with it.
	_, err := crypto.NewKeyPair(configuration.Global.DataDir, moniker, PasswordFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
