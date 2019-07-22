package keys

import (
	"github.com/mosaicnetworks/monetd/cmd/monetd/commands/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	passwordFile      string
	outputJSON        bool
	monikerParam      string
	monetCliConfigDir string
)

//KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "monet key manager",
	Long: `Keys
	
Monet Key Manager. `,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	KeysCmd.AddCommand(
		NewGenerateCmd(),
		keys.NewInspectCmd(),
		keys.NewUpdateCmd(),
		keys.NewNewCmd(),
	)

	//Commonly used command line flags
	KeysCmd.PersistentFlags().StringVar(&keys.PasswordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	KeysCmd.PersistentFlags().BoolVar(&keys.OutputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(KeysCmd.Flags())
}
