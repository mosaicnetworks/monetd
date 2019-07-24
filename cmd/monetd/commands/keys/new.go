package keys

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/spf13/cobra"
)

// newNewCmd returns the command that creates a new keypair
func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [moniker]",
		Short: "create a new keypair",
		Long: `
monetd keys new [moniker]

This command produces two files:

- [datadir]/keystore/[moniker].json : The encrypted keyfile
- [datadir]/keystore/[moniker].toml : Key metadata
		`,

		Args: cobra.ExactArgs(1),
		RunE: newkeys,
	}

	return cmd
}

func newkeys(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	// key is returned, but we don't want to do anything with it.
	_, err := crypto.NewKeyPair(configuration.Configuration.DataDir, moniker, PasswordFile)

	return err
}
