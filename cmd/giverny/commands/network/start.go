package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	eth_keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/keystore"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//CLI flags
var forceNetwork = false
var useExisting = false
var startNodes = false

type copyRecord struct {
	from string
	to   string
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
giverny network start

Starts a network. Does not start individual nodes. 
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}

	addStartFlags(cmd)

	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&forceNetwork, "force-network", forceNetwork, "force network down if already exists")
	cmd.Flags().BoolVar(&useExisting, "use-existing", useExisting, "use existing network if already exists")
	cmd.Flags().BoolVar(&startNodes, "start-nodes", startNodes, "start nodes")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	network := args[0]

	if err := startDockerNetwork(network); err != nil {
		return err
	}

	return nil
}

func startDockerNetwork(networkName string) error {

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

	if conf.Docker.Name == "" {
		return errors.New("network " + networkName + " is not configured as a docker network")
	}

	// Create a Docker Client

	common.DebugMessage("Connecting to Docker Client")

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	// Create a Docker Network
	networkID, err := docker.SafeCreateNetwork(cli,
		conf.Docker.Name,
		conf.Docker.Subnet,
		conf.Docker.IPRange,
		conf.Docker.Gateway,
		forceNetwork, useExisting)
	if err != nil {
		return err
	}
	common.DebugMessage(fmt.Sprintf("Created Network %s (%s)", conf.Docker.Name, networkID))

	// Next we build the docker configurations to get all of the configs ready to
	// push

	err = exportDockerConfigs(&conf)
	if err != nil {
		return err
	}

	if startNodes {
		for _, n := range conf.Nodes {
			if !n.NonNode {
				common.DebugMessage("Starting node " + n.Moniker)
				if err := pushDockerNode(networkName, n.Moniker, imgName, imgIsRemote); err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func exportDockerConfigs(conf *Config) error {

	// Configure some paths
	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, conf.Network.Name)
	dockerDir := filepath.Join(networkDir, givernyDockerDir)
	err := files.CreateDirsIfNotExists([]string{dockerDir})
	if err != nil {
		return err
	}

	for _, n := range conf.Nodes { // loop around nodes
		if !n.NonNode {
			if err := exportDockerNodeConfig(networkDir, dockerDir, &n); err != nil {
				return err
			}
		}
	}

	return nil
}

func exportDockerNodeConfig(networkDir, dockerDir string, n *node) error {

	netaddr := n.NetAddr
	if !strings.Contains(netaddr, ":") {
		netaddr += ":" + monetconfig.DefaultGossipPort
	}
	// Build output files

	if n.Moniker != "" { // Should not be blank here, but safety first

		monetDir := filepath.Join(dockerDir, n.Moniker, monetconfig.MonetdTomlDirDot)
		configDir := filepath.Join(monetDir, monetconfig.ConfigDir)
		babbleConfigDir := filepath.Join(configDir, monetconfig.BabbleDir)
		ethConfigDir := filepath.Join(configDir, monetconfig.EthDir)
		keystoreDir := filepath.Join(monetDir, monetconfig.KeyStoreDir)

		common.DebugMessage("Creating config in " + monetDir)

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
		nodeConfigFile := filepath.Join(dockerDir, n.Moniker+".toml")
		nodeConfig := dockerNodeConfig{
			Moniker: n.Moniker,
			NetAddr: strings.Split(netaddr, ":")[0],
		}

		tomlBytes, err := toml.Marshal(nodeConfig)
		if err != nil {
			return err
		}

		err = files.WriteToFile(nodeConfigFile, string(tomlBytes), files.OverwriteSilently)
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

func setListenAddressInToml(toml string, listen string) error {
	// For a simple change, tree is quicker and easier than unmarshalling the
	// whole tree
	tree, err := files.LoadToml(toml)
	if err != nil {
		return err
	}

	tree.SetPath([]string{"babble", "listen"}, listen)
	files.SaveToml(tree, toml)
	if err != nil {
		return err
	}

	return nil
}

func generateBabblePrivateKey(keyfile, pwdfile, moniker, outDir string) error {

	if moniker == "" {
		return nil
	} // If account not set, do nothing

	if !files.CheckIfExists(keyfile) {
		return errors.New("cannot read keyfile: " + keyfile)
	}

	if !files.CheckIfExists(pwdfile) {
		common.DebugMessage("No passphrase file available")
		pwdfile = ""
	}

	keyjson, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfile, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := crypto.GetPassphrase(pwdfile, false)
	if err != nil {
		return err
	}

	key, err := eth_keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	addr := key.Address.Hex()

	err = keystore.DumpPrivKey(outDir, key.PrivateKey)
	if err != nil {
		return fmt.Errorf("Error writing raw key: %v", err)
	}

	common.DebugMessage("Written Private Key for " + addr)

	return nil
}
