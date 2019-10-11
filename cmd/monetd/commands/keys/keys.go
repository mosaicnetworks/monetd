package keys

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables for persistent Keys options
var (
	_keystore     = configuration.DefaultKeystoreDir()
	_passwordFile string
	_outputJSON   bool
)

// KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "manage keys",
	Long: `
Manage keys in the <keystore> folder.

Note that other Monet tools, like monetcli and monet-wallet, use the same 
default keystore.

+------------------------------------------------------------------------------+ 
| Please take all the necessary precautions to secure these files and remember | 
| the passwords, as it will be impossible to recover the keys without them.    |
+------------------------------------------------------------------------------+

Keys are associated with monikers and encrypted in password-protected files in
<keystore>/[moniker].json. Keyfiles contain JSON encoded objects, which Ethereum
users will recognise as the de-facto Ethereum keyfile format. Indeed, Monet and
the underlying consensus algorithm, Babble, use the same type of keys as
Ethereum. A key can be used to run a validator node, or to control an account
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
	KeysCmd.PersistentFlags().StringVar(&_keystore, "keystore", _keystore, "keystore directory")
	KeysCmd.PersistentFlags().StringVar(&_passwordFile, "passfile", "", "file containing the passphrase")
	KeysCmd.PersistentFlags().BoolVar(&_outputJSON, "json", false, "output JSON instead of human-readable format")

	viper.BindPFlags(KeysCmd.Flags())
}
