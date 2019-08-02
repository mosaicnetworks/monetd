package network

import (
	"errors"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopAndDelete = false

func newStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [network] [node] ",
		Short: "stop a network",
		Long: `
giverny network stop

Stop a node and all the nodes within it.
		`,
		Args: cobra.MinimumNArgs(1),
		RunE: networkStop,
	}

	addStopFlags(cmd)

	return cmd
}

func addStopFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&stopAndDelete, "remove", stopAndDelete, "stop and remove node")
	viper.BindPFlags(cmd.Flags())
}

func networkStop(cmd *cobra.Command, args []string) error {

	if len(args) == 1 { // Network
		return stopNetwork(args[0], stopAndDelete)
	}

	return stopNode(args[0], args[1], stopAndDelete)

}

func stopNetwork(networkName string, force bool) error {

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	common.DebugMessage("List Networks")

	nets, err := docker.GetNetworks(cli, false)
	if err != nil {
		return err
	}
	common.DebugMessage("\nNumber of Networks: " + strconv.Itoa(len(nets)))

	common.DebugMessage("Testing for network: " + networkName)
	if netID, ok := nets[networkName]; ok {
		common.DebugMessage("Network ID is: " + netID)
		common.DebugMessage("Removing Network " + networkName)

		cons, err := docker.GetContainers(cli, false)
		if err != nil {
			return err
		}
		for moniker, containerID := range cons {
			common.DebugMessage("Stopping Container " + moniker)
			err = docker.StopContainer(cli, containerID)
			if err != nil {
				return err
			}

			if force {
				common.DebugMessage("Removing Container " + containerID)
				err = docker.RemoveContainer(cli, containerID, true, false, false)
				if err != nil {
					return err
				}
			}
		}

		if err := docker.RemoveNetwork(cli, netID); err != nil {
			return err
		}
		common.DebugMessage("")
		nets, err = docker.GetNetworks(cli, false)
		if err != nil {
			return err
		}

		common.DebugMessage("\nNumber of Networks: " + strconv.Itoa(len(nets)))
	} else {
		return errors.New("network not found")
	}

	return nil
}

func stopNode(networkName, nodeName string, force bool) error {

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	cons, err := docker.GetContainers(cli, false)
	if err != nil {
		return err
	}

	containerID, ok := cons[nodeName]
	if !ok {
		return errors.New("container " + nodeName + " is not running")
	}

	common.DebugMessage("Stopping Container " + containerID)
	err = docker.StopContainer(cli, containerID)
	if err != nil {
		return err
	}

	if force {
		common.DebugMessage("Removing Container " + containerID)
		err = docker.RemoveContainer(cli, containerID, true, false, false)
		if err != nil {
			return err
		}
	}

	return nil
}
