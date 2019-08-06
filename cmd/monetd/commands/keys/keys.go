package keys

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables for persistent Keys options
var (
	PasswordFile string
	OutputJSON   bool
)

// KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "monet key manager",
	Long: `
Manage keys in the [datadir]/keystore folder.

Note that other Monet tools, like monetcli and monet-wallet, use the same 
default [datadir]/keystore.

+------------------------------------------------------------------------------+ 
| Please take all the necessary precautions to secure these files and remember | 
| the passwords, as it will be impossible to recover the keys without them.    |
+------------------------------------------------------------------------------+

Keys are associated with monikers and encrypted in password-protected files in
[datadir]/keystore/[moniker].json. Keyfiles contain JSON encoded objects, which
Ethereum users will recognise as the de-facto Ethereum keyfile format. Indeed,
Monet and the underlying consensus algorithm, Babble, use the same type of keys
as Ethereum. A key can be used to run a validator node, or to control an account
with a token balance.
`,
	TraverseChildren: true,
}

func init() {
	// Keys subcommands
	KeysCmd.AddCommand(
		newInspectCmd(),
		newUpdateCmd(),
		newNewCmd(),
		newListCmd(),
	)

	// Flags that are common to all Keys subcommands
	KeysCmd.PersistentFlags().StringVar(&PasswordFile, "passfile", "", "file containing the passphrase")
	KeysCmd.PersistentFlags().BoolVar(&OutputJSON, "json", false, "output JSON instead of human-readable format")

	viper.BindPFlags(KeysCmd.Flags())
}
