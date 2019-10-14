package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func newAWSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws [network] [output-path]",
		Short: "write aws configuration files",
		Long: `
giverny network aws

Writes AWS configuration. 
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkAWS,
	}

	return cmd
}

func networkAWS(cmd *cobra.Command, args []string) error {
	network := args[0]
	outPath := args[1]

	if !files.CheckIfExists(outPath) {
		return errors.New("cannot find the output folder, " + outPath + " for " + network)
	}

	if err := buildNetworkConfig(network, outPath); err != nil {
		return err
	}

	return nil
}

func buildNetworkConfig(networkName string, outPath string) error {

	// Set some paths
	thisNetworkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)

	// Check expected config exists
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	var conf = Config{}

	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}

	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}

	common.DebugMessage("Configuring Network ", conf.Docker.Name)

	err = exportAWSConfigs(&conf, outPath)
	if err != nil {
		return err
	}

	return nil
}

func exportAWSConfigs(conf *Config, outPath string) error {

	// Configure some paths
	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, conf.Network.Name)
	err := files.CreateDirsIfNotExists([]string{outPath})
	if err != nil {
		return err
	}

	for _, n := range conf.Nodes { // loop around nodes
		if !n.NonNode {
			if err := exportAWSNodeConfig(networkDir, outPath, &n); err != nil {
				return err
			}
		}
	}

	return nil
}

func exportAWSNodeConfig(networkDir, outPath string, n *node) error {

	netaddr := n.NetAddr
	if !strings.Contains(netaddr, ":") {
		netaddr += ":" + monetconfig.DefaultGossipPort
	}
	// Build output files

	if n.Moniker != "" { // Should not be blank here, but safety first

		monetDir := filepath.Join(outPath, n.Moniker, monetconfig.MonetdTomlDirDot)
		configDir := filepath.Join(monetDir, monetconfig.ConfigDir)
		babbleConfigDir := filepath.Join(configDir, monetconfig.BabbleDir)
		ethConfigDir := filepath.Join(configDir, monetconfig.EthDir)
		keystoreDir := filepath.Join(monetDir, monetconfig.KeyStoreDir)

		common.DebugMessage("Creating config in " + configDir)

		err := files.CreateDirsIfNotExists([]string{
			babbleConfigDir,
			ethConfigDir,
			keystoreDir,
		})
		if err != nil {
			return err
		}

		copying := []copyRecord{
			{ // monetd.toml
				from: filepath.Join(networkDir, monetconfig.MonetTomlFile),
				to:   filepath.Join(configDir, monetconfig.MonetTomlFile),
			},
			{ // eth/genesis.json
				from: filepath.Join(networkDir, monetconfig.GenesisJSON),
				to:   filepath.Join(ethConfigDir, monetconfig.GenesisJSON),
			},
			{ // babble/peers.json
				from: filepath.Join(networkDir, monetconfig.PeersJSON),
				to:   filepath.Join(babbleConfigDir, monetconfig.PeersJSON),
			},
			{ // babble/peers.genesis.json
				from: filepath.Join(networkDir, monetconfig.PeersGenesisJSON),
				to:   filepath.Join(babbleConfigDir, monetconfig.PeersGenesisJSON),
			},
			{ // keystore/<moniker>.json (private key)
				from: filepath.Join(networkDir, monetconfig.KeyStoreDir, n.Moniker+".json"),
				to:   filepath.Join(keystoreDir, n.Moniker+".json"),
			},
			{ // keystore/<moniker>.text (password)
				from: filepath.Join(networkDir, monetconfig.KeyStoreDir, n.Moniker+".txt"),
				to:   filepath.Join(keystoreDir, n.Moniker+".txt"),
			},
		}

		for _, f := range copying {
			files.CopyFileContents(f.from, f.to)
		}

		// Write a node description file containing all of the parameters needed
		// to start a container. Saves having to load and parse network.toml for
		//  every node
		nodeConfigFile := filepath.Join(outPath, n.Moniker+".toml")
		nodeConfig := dockerNodeConfig{
			Moniker: n.Moniker,
			NetAddr: strings.Split(netaddr, ":")[0],
		}

		tomlBytes, err := toml.Marshal(nodeConfig)
		if err != nil {
			return err
		}

		err = files.WriteToFile(nodeConfigFile, string(tomlBytes), 0)
		if err != nil {
			return err
		}

		// edit monetd.toml and set babble.listen appropriately
		err = setListenAddressInToml(
			filepath.Join(configDir, monetconfig.MonetTomlFile),
			netaddr)
		if err != nil {
			return err
		}

		// decrypt the validator private key, and dump it into the babble config
		// dir (priv_key)
		err = generateBabblePrivateKey(
			filepath.Join(keystoreDir, n.Moniker+".json"),
			filepath.Join(keystoreDir, n.Moniker+".txt"),
			n.Moniker,
			babbleConfigDir)
		if err != nil {
			return err
		}

	}
	return nil
}
