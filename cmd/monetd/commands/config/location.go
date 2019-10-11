package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show default configuration files",
		Long:  "Show the default locations of monetd configuration files.",
		RunE:  locationConfig,
	}

	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {

	fmt.Println("Monetd Config        : " + filepath.Join(configuration.DefaultConfigDir(), configuration.MonetTomlFile))

	fmt.Println("Babble Peers         : " + filepath.Join(configuration.DefaultConfigDir(), configuration.BabbleDir, configuration.PeersJSON))
	fmt.Println("Babble Genesis Peers : " + filepath.Join(configuration.DefaultConfigDir(), configuration.BabbleDir, configuration.PeersGenesisJSON))
	fmt.Println("Babble Private Key   : " + filepath.Join(configuration.DefaultConfigDir(), configuration.BabbleDir, configuration.BabblePrivKey))

	fmt.Println("EVM-Lite Genesis     : " + filepath.Join(configuration.DefaultConfigDir(), configuration.EthDir, configuration.GenesisJSON))

	fmt.Println("Babble Database      : " + filepath.Join(configuration.DefaultDataDir(), configuration.BabbleDB))
	fmt.Println("EVM-Lite Database    : " + filepath.Join(configuration.DefaultDataDir(), configuration.EthDB))

	fmt.Println("Keystore        : " + configuration.DefaultKeystoreDir())

	return nil
}
