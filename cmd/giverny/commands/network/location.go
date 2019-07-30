package network

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	mconfiguration "github.com/mosaicnetworks/monetd/src/configuration"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location [network_name]",
		Short: "show the location of the configuration files",
		Long: `
giverny network location
		`,
		Args: cobra.ArbitraryArgs,
		RunE: networkLocation,
	}

	addLocationFlags(cmd)

	return cmd
}

func addLocationFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkLocation(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		fmt.Println(configuration.GivernyConfigDir)
		return nil
	}

	networkName = strings.TrimSpace(args[0])

	if !common.CheckMoniker(networkName) {
		return errors.New("the network name, " + networkName + ", is invalid")
	}

	common.InfoMessage("Network                 : " + networkName)

	common.InfoMessage("Giverny Config Dir      : " + configuration.GivernyConfigDir)
	common.InfoMessage("Giverny Networks Dir    : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName))
	common.InfoMessage("Giverny KeyStore Dir    : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, givernyKeystoreDir))
	common.InfoMessage("Peers JSON              : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, mconfiguration.PeersJSON))
	common.InfoMessage("Peers Genesis JSON      : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, mconfiguration.PeersGenesisJSON))
	common.InfoMessage("Genesis JSON            : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, mconfiguration.GenesisJSON))
	common.InfoMessage("Monetd TOML             : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, mconfiguration.MonetTomlFile))
	common.InfoMessage("Monetd TOML             : " +
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, networkTomlFileName))

	return nil
}
