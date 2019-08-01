package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/api/types/strslice"
	"github.com/mosaicnetworks/monetd/src/docker"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [network] [node]",
		Short: "push the configuration for a node on the named network",
		Long: `
giverny network push
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkPush,
	}

	addPushFlags(cmd)

	return cmd
}

const imgName = "mosaicnetworks/monetd:latest"

var imgIsRemote = false

func addPushFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkPush(cmd *cobra.Command, args []string) error {
	networkName := args[0]
	nodeName := args[1]
	return pushDockerNode(networkName, nodeName, "", imgName, imgIsRemote)
}

//PushDockerNode builds a docker node, configures it and starts it
func pushDockerNode(networkName, nodeName, networkID, imgName string, isRemoteImage bool) error {
	common.DebugMessage("Pushing network " + networkName + " node " + nodeName)

	// First we validate that the requested node has been created
	dockerpath := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, givernyDockerDir)
	if !files.CheckIfExists(dockerpath) {
		return errors.New(" cannot find docker config for network " + networkName + ". Have you run giverny network start? ")
	}

	dockerconfigpath := filepath.Join(dockerpath, nodeName)
	if !files.CheckIfExists(dockerconfigpath) {
		return errors.New(" cannot find docker config folder for node " + nodeName)
	}

	dockerconfig := filepath.Join(dockerpath, nodeName+".toml")
	if !files.CheckIfExists(dockerconfig) {
		return errors.New(" cannot find docker config toml for node " + nodeName)
	}

	// Get Node Details

	// Read key from file.
	tomlfile, err := ioutil.ReadFile(dockerconfig)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", dockerconfig, err)
	}

	config := dockerNodeConfig{}
	toml.Unmarshal(tomlfile, &config)

	common.DebugMessage("Container IP is " + config.NetAddr)

	// Start Docker Client
	common.DebugMessage("Connecting to Docker Client\n ")

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	// If we don't have a networkID we retrieve one

	if nets, err := docker.GetNetworks(cli, false); err == nil {
		if net, ok := nets[networkName]; ok {
			networkID = net
		} else {
			return errors.New("network " + networkName + " is not running")
		}
	} else {
		common.ErrorMessage("Error getting network status")
		return nil
	}

	// Check current containers to see if node already exists
	containers, err := docker.GetContainers(cli, false)

	if existingNode, ok := containers[nodeName]; ok {
		return errors.New("node " + nodeName + " already exists (" + existingNode + ")")
	}

	// Create Node
	common.DebugMessage("Creating Container ")

	containerID, err := docker.CreateContainerFromImage(cli, imgName, isRemoteImage,
		nodeName, strslice.StrSlice{"run"}, false)

	common.DebugMessage("Created Container " + containerID)

	// Copy Configuration to Node

	common.DebugMessage("Copying Config to Container ")
	err = docker.CopyToContainer(cli, containerID, dockerconfigpath, "/")
	if err != nil {
		return err
	}

	// Configure Networking

	//	func ConnectContainerToNetwork(cli *client.Client, networkID string, containerID string, ip string) error {
	common.DebugMessage("Connecting Container to Network")

	err = docker.ConnectContainerToNetwork(cli, networkID, containerID, config.NetAddr)
	if err != nil {
		return err
	}

	// Start Node
	common.DebugMessage("Starting Container ")
	err = docker.StartContainer(cli, containerID)
	if err != nil {
		return err
	}

	common.DebugMessage("Container Started")

	return nil
}
