package trans

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//KeysCmd is an Ethereum key manager
var TransCmd = &cobra.Command{
	Use:   "trans",
	Short: "claude trans commands",
	Long: `Trans
	
Claude Trans utilities. `,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	TransCmd.AddCommand(
		newDecodeCmd(),
	)

	//Commonly used command line flags
	//	KeysCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	//	KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(TransCmd.Flags())
}
