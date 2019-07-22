package keys

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	passwordFile string
	outputJSON   bool
	monikerParam string
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
		newImportCmd(),
	)

	//Commonly used command line flags
	KeysCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(KeysCmd.Flags())
}
