package transactions

import (
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

var (
	_keystore = monetconfig.DefaultKeystoreDir()
	_giverny  = configuration.GivernyConfigDir
)

//TODO duplicates the definition in networks package.
//Probably better to publish them and use them directly.
const (
	givernyNetworksDir     = "networks"
	givernyKeystoreDir     = "keystore"
	givernyTransactionsDir = "trans"
	networkTomlFileName    = "network.toml"
)

//TransCmd implements the transactions subcommand
var TransCmd = &cobra.Command{
	Use:   "transactions",
	Short: "giverny transactions",
	Long: `Server
	
The giverny transaction command is used to generate sets of transactions for 
testing networks.`,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	TransCmd.AddCommand(
		newGenerateCmd(),
		newSoloCmd(),
	)

	TransCmd.PersistentFlags().StringVarP(&_keystore, "keystore", "m", _keystore, "keystore directory")
	TransCmd.PersistentFlags().StringVarP(&_giverny, "dir", "d", _giverny, "giverny directory")

}
