package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/config"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/mosaicnetworks/monetd/src/files"
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

Starts a network. Does not start individual nodes. The --force-network parameter
stops and restarts the network. 
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

	// Check expect config exists
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
	networkID, err := docker.SafeCreateNetwork(cli, conf.Docker.Name,
		conf.Docker.Subnet, conf.Docker.IPRange, conf.Docker.Gateway, forceNetwork, useExisting)
	if err != nil {
		return err
	}
	common.DebugMessage("Created Network " + networkID)

	// Next we build the docker configurations to get all of the configs ready to
	// push

	err = exportDockerConfigs(&conf)
	if err != nil {
		return err
	}

	if startNodes {
		for _, n := range conf.Nodes {
			common.DebugMessage("Starting node " + n.Moniker)
			if err := pushDockerNode(networkName, n.Moniker, networkID, imgName, imgIsRemote); err != nil {
				return err
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
		if err := exportDockerNodeConfig(networkDir, dockerDir, &n); err != nil {
			return err
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
		nodeDir := filepath.Join(dockerDir, n.Moniker)
		// Docker container will always use .monet
		monetDir := filepath.Join(nodeDir, monetconfig.MonetdTomlDirDot)

		common.DebugMessage("Creating config in " + nodeDir)
		err := files.CreateDirsIfNotExists([]string{
			nodeDir,
			monetDir,
			filepath.Join(monetDir, monetconfig.BabbleDir),
			filepath.Join(monetDir, monetconfig.EthDir),
			filepath.Join(monetDir, monetconfig.KeyStoreDir),
		})
		if err != nil {
			return err
		}

		copying := []copyRecord{
			{from: filepath.Join(networkDir, monetconfig.GenesisJSON),
				to: filepath.Join(monetDir, monetconfig.EthDir, monetconfig.GenesisJSON)},
			{from: filepath.Join(networkDir, monetconfig.PeersJSON),
				to: filepath.Join(monetDir, monetconfig.BabbleDir, monetconfig.PeersJSON)},
			{from: filepath.Join(networkDir, monetconfig.PeersGenesisJSON),
				to: filepath.Join(monetDir, monetconfig.BabbleDir, monetconfig.PeersGenesisJSON)},
			{from: filepath.Join(networkDir, monetconfig.MonetTomlFile),
				to: filepath.Join(monetDir, monetconfig.MonetTomlFile)},
			{from: filepath.Join(networkDir, monetconfig.KeyStoreDir, n.Moniker+".json"),
				to: filepath.Join(monetDir, monetconfig.KeyStoreDir, n.Moniker+".json")},
			{from: filepath.Join(networkDir, monetconfig.KeyStoreDir, n.Moniker+".txt"),
				to: filepath.Join(monetDir, monetconfig.KeyStoreDir, n.Moniker+".txt")},
		}

		for _, f := range copying {
			files.CopyFileContents(f.from, f.to)
		}

		// Write a node description file containing all of the parameters needed to start a container
		// Saves having to load and parse network.toml for every node
		nodeConfigFile := filepath.Join(dockerDir, n.Moniker+".toml")
		nodeConfig := dockerNodeConfig{Moniker: n.Moniker, NetAddr: strings.Split(netaddr, ":")[0]}

		tomlBytes, err := toml.Marshal(nodeConfig)
		if err != nil {
			return err
		}

		err = files.WriteToFile(nodeConfigFile, string(tomlBytes))
		if err != nil {
			return err
		}

		// Need to edit monetd.toml and set datadir and listen appropriately

		err = config.SetLocalParamsInToml("/.monet", filepath.Join(monetDir, monetconfig.MonetTomlFile), netaddr)
		if err != nil {
			return err
		}

		// Need to generate private key
		err = config.GenerateBabblePrivateKey(monetDir, n.Moniker)
		if err != nil {
			return err
		}

	}
	return nil
}
