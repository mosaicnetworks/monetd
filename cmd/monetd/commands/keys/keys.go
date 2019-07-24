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
This command manages keys in the [datadir]/keystore folder.

Each key is associated with a moniker and encrypted in a password protected 
file. The moniker is a friendly name preventing users from having to type or 
copy/paste long character strings in the terminal. The password-protected file 
contains a JSON formatted string, which Ethereum users will recognise as the 
de-facto Ethereum keyfile format. Indeed, Monet and the underlying consensus 
algorithm, Babble, use the same type of keys as Ethereum. The same key can be 
used to run a validator node, or to control an account in Monet with a Tenom 
balance.

To use a key as part of a validator node running monetd, it will have to be 
decrypted with the password and copied over to [datadir]/babble/priv_key. The 
command  'monetd config build' does this automatically, but it can also be done 
manually with the help of the 'monetd keys inspect --private' command. 

Note that other Monet tools, like monetcli and monet-wallet, use the same 
default [datadir]/keystore.

+------------------------------------------------------------------------------+ 
| Please take all the necessary precautions to secure these files and remember | 
| the password, as it will be impossible to recover the key without them.      |
+------------------------------------------------------------------------------+
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
