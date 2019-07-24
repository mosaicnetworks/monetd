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
	Use:              "keys",
	Short:            "monet key manager",
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
