package config

import (
	"fmt"
	"path/filepath"

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
		RunE:  locationConfig,
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
		fmt.Println("Babble Dir    : " + filepath.Join(configuration.Global.DataDir, configuration.BabbleDir))
		fmt.Println("EVM-Lite Dir  : " + filepath.Join(configuration.Global.DataDir, configuration.EthDir))
		fmt.Println("Keystore Dir  : " + filepath.Join(configuration.Global.DataDir, configuration.KeyStoreDir))
		fmt.Println("Config File   : " + filepath.Join(configuration.Global.DataDir, configuration.MonetTomlFile))
		fmt.Println("Wallet Config : " + filepath.Join(configuration.Global.DataDir, configuration.WalletTomlFile))
		fmt.Println("Peers         : " + filepath.Join(configuration.Global.DataDir, configuration.BabbleDir, configuration.PeersJSON))
		fmt.Println("Genesis Peers : " + filepath.Join(configuration.Global.DataDir, configuration.BabbleDir, configuration.PeersGenesisJSON))
		fmt.Println("Genesis File  : " + filepath.Join(configuration.Global.DataDir, configuration.EthDir, configuration.GenesisJSON))
	} else {
		fmt.Println(configuration.Global.DataDir)
	}

	return nil
}
