package keys

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	_keystore = configuration.DefaultKeystoreDir()
)

//KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:              "keys",
	Short:            "monet key manager",
	TraverseChildren: true,
}

func init() {
	//Subcommands
	KeysCmd.AddCommand(
		newImportCmd(),
		newGenerateCmd(),
	)

	//Commonly used command line flags
	KeysCmd.PersistentFlags().StringVarP(&_keystore, "keystore", "k", _keystore, "keystore directory")
	viper.BindPFlags(KeysCmd.Flags())
}
