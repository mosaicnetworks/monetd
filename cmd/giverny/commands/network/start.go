package network

import (
	"errors"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forceNetwork = false

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
giverny network start

Starts a network. Does not start individual nodes
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}

	addStartFlags(cmd)

	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&forceNetwork, "force-network", forceNetwork, "force network down if already exists")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	network := args[0]

	return startDockerNetwork(network)
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

	// Load Toml file to a tree
	tree, err := files.LoadToml(networkTomlFile)
	if err != nil {
		common.ErrorMessage("Cannot load network.toml file: ", networkTomlFile)
		return err
	}

	var dockerNetworkName, dockerbaseip, dockersubnet, dockeriprange, dockergateway string

	if tree.HasPath([]string{"docker", "name"}) {
		dockerNetworkName = tree.GetPath([]string{"docker", "name"}).(string)
	}
	if tree.HasPath([]string{"docker", "baseip"}) {
		dockerbaseip = tree.GetPath([]string{"docker", "baseip"}).(string)
	}
	if tree.HasPath([]string{"docker", "subnet"}) {
		dockersubnet = tree.GetPath([]string{"docker", "subnet"}).(string)
	}
	if tree.HasPath([]string{"docker", "iprange"}) {
		dockeriprange = tree.GetPath([]string{"docker", "iprange"}).(string)
	}
	if tree.HasPath([]string{"docker", "gateway"}) {
		dockergateway = tree.GetPath([]string{"docker", "gateway"}).(string)
	}

	common.DebugMessage("Configuring Network ", dockerNetworkName)
	common.DebugMessage("Base IP:    ", dockerbaseip)
	common.DebugMessage("IP Range:   ", dockeriprange)
	common.DebugMessage("Subnet:     ", dockersubnet)
	common.DebugMessage("Gateway:    ", dockergateway)

	// Create a Docker Client

	common.DebugMessage("Connecting to Docker Client")

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	// Create a Docker Network

	networkID, err := docker.SafeCreateNetwork(cli, dockerNetworkName,
		dockersubnet, dockeriprange, dockergateway, forceNetwork)

	if err != nil {
		return err
	}

	common.DebugMessage("Created Network " + networkID)

	return nil
}

func startDocker() error {

	/*
	   #!/bin/bash

	   set -eux

	   N=${1:-4}
	   FASTSYNC=${2:-false}
	   MPWD=$(pwd)


	   docker network create \
	     --driver=bridge \
	     --subnet=172.77.0.0/16 \
	     --ip-range=172.77.0.0/16 \""
	     --gateway=172.77.5.254 \
	     babblenet

	   for i in $(seq 1 $N)
	   do
	       docker run -d --name=client$i --net=babblenet --ip=172.77.10.$i -it mosaicnetworks/dummy:0.5.0 \
	       --name="client $i" \
	       --client-listen="172.77.10.$i:1339" \
	       --proxy-connect="172.77.5.$i:1338" \
	       --discard \
	       --log="debug"
	   done

	   for i in $(seq 1 $N)
	   do
	       docker create --name=node$i --net=babblenet --ip=172.77.5.$i mosaicnetworks/babble:0.5.0 run \
	       --heartbeat=100ms \
	       --moniker="node$i" \
	       --cache-size=50000 \
	       --listen="172.77.5.$i:1337" \
	       --proxy-listen="172.77.5.$i:1338" \
	       --client-connect="172.77.10.$i:1339" \
	       --service-listen="172.77.5.$i:80" \
	       --sync-limit=500 \
	       --fast-sync=$FASTSYNC \
	       --store \
	       --log="debug"

	       docker cp $MPWD/conf/node$i node$i:/.babble
	       docker start node$i
	   done


	*/
	return nil
}
