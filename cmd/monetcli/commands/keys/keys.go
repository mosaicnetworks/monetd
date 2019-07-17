package keys

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	passwordFile   string
	outputJSON     bool
	monikerParam   string
	monetConfigDir string
)

//KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "an Ethereum key manager",
	Long: `Keys
	
An Ethereum key manager. `,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	KeysCmd.AddCommand(
		NewGenerateCmd(),
		NewInspectCmd(),
		NewUpdateCmd(),
		NewNewCmd(),
	)

	//Commonly used command line flags
	KeysCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(KeysCmd.Flags())
}
