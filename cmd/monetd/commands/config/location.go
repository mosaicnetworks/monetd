package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	longFormat = false
)

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `
The location subcommand shows the location of the monetd configuration files. It 
respects any --datadir parameter. 

If you specify --expanded then a list of configuration folders and directories
is output.`,
		RunE: locationConfig,
	}

	addLocationFlags(cmd)
	return cmd
}

func addLocationFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&longFormat, "expanded", "x", longFormat, "show expanded information")
	viper.BindPFlags(cmd.Flags())
}

func locationConfig(cmd *cobra.Command, args []string) error {

	if longFormat {
		fmt.Println("Config root   : " + configuration.Global.DataDir)
		fmt.Println("Babble Dir    : " + filepath.Join(configuration.Global.DataDir, common.BabbleDir))
		fmt.Println("EVM-Lite Dir  : " + filepath.Join(configuration.Global.DataDir, common.EthDir))
		fmt.Println("Keystore Dir  : " + filepath.Join(configuration.Global.DataDir, common.KeyStoreDir))
		fmt.Println("Config File   : " + filepath.Join(configuration.Global.DataDir, common.MonetTomlFile))
		fmt.Println("Wallet Config : " + filepath.Join(configuration.Global.DataDir, common.WalletTomlFile))
		fmt.Println("Peers         : " + filepath.Join(configuration.Global.DataDir, common.BabbleDir, common.PeersJSON))
		fmt.Println("Genesis Peers : " + filepath.Join(configuration.Global.DataDir, common.BabbleDir, common.PeersGenesisJSON))
		fmt.Println("Genesis File  : " + filepath.Join(configuration.Global.DataDir, common.EthDir, common.GenesisJSON))
	} else {
		fmt.Println(configuration.Global.DataDir)
	}

	return nil
}
