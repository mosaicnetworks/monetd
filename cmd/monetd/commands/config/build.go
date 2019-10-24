package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"

	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/genesis"
	"github.com/mosaicnetworks/monetd/src/keystore"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newBuildCmd initialises a bare-bones configuration for monetd
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [moniker]",
		Short: "create the configuration for a single-node network",
		Long: `
Create the configuration for a single-node network.

Use the keystore account identified by [moniker] to define a network with a
single node. All the accounts in <keystore> are also credited with a large 
amount of tokens in the genesis file. This command is mostly used for testing.
If the --address flag is omitted, the first non-loopback address is used.
`,
		Args: cobra.ExactArgs(1),
		RunE: buildConfig,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&_keystore, "keystore", _keystore, "keystore directory")
	cmd.Flags().StringVar(&_configDir, "config", _configDir, "output directory")
	cmd.Flags().StringVar(&_addressParam, "address", _addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&_passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func buildConfig(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	address := fmt.Sprintf("%s:%s", _addressParam, configuration.DefaultGossipPort)

	// Some debug output confirming parameters
	common.DebugMessage("Building Config for   : ", moniker)
	common.DebugMessage("Using Network Address : ", address)
	common.DebugMessage("Using Password File   : ", _passwordFile)

	// set global config moniker
	configuration.Global.Babble.Moniker = moniker

	// Retrieve the keyfile corresponding to moniker
	privateKey, err := keystore.GetKey(_keystore, moniker, _passwordFile)
	if err != nil {
		return err
	}

	// Create Directories if they don't exist
	CreateMonetConfigFolders(_configDir)

	// Copy the key to babble directory with appropriate permissions
	err = keystore.DumpPrivKey(
		filepath.Join(_configDir, configuration.BabbleDir),
		privateKey)
	if err != nil {
		return err
	}

	pubKey := hex.EncodeToString(eth_crypto.FromECDSAPub(&privateKey.PublicKey))

	// Create a peer-set whith a single node
	peers := []*peers.Peer{
		peers.NewPeer(pubKey, address, moniker),
	}

	// Write peers.json and peers.genesis.json
	if err := dumpPeers(_configDir, peers); err != nil {
		return err
	}

	// Create the eth/genesis.json file
	err = genesis.GenerateGenesisJSON(
		filepath.Join(_configDir, configuration.EthDir),
		_keystore,
		peers,
		nil,
		configuration.DefaultContractAddress,
		configuration.DefaultControllerContractAddress)
	if err != nil {
		return err
	}

	// Write TOML file for monetd based on global config object
	err = configuration.DumpGlobalTOML(
		_configDir,
		configuration.MonetTomlFile,
		true)
	if err != nil {
		return err
	}

	return nil
}

// dumpPeers takes a list of peers and dumps it into peers.json and
// peers.genesis.json in the babble directory
func dumpPeers(configDir string, peers []*peers.Peer) error {
	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}

	// peers.json
	jsonFileName := filepath.Join(configDir, configuration.BabbleDir, configuration.PeersJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)

	// peers.genesis.json
	jsonFileName = filepath.Join(configDir, configuration.BabbleDir, configuration.PeersGenesisJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)

	return nil
}
