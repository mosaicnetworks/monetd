package docker

import (
	"errors"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/common"
)

//StopNetwork stops a network. Force removes nodes too.
func StopNetwork(networkName string, force bool) error {

	cli, err := GetDockerClient()
	if err != nil {
		return err
	}

	common.DebugMessage("List Networks")

	nets, err := GetNetworks(cli, false)
	if err != nil {
		return err
	}
	common.DebugMessage("\nNumber of Networks: " + strconv.Itoa(len(nets)))

	common.DebugMessage("Testing for network: " + networkName)
	if netID, ok := nets[networkName]; ok {
		common.DebugMessage("Network ID is: " + netID)
		common.DebugMessage("Removing Network " + networkName)

		cons, err := GetContainers(cli, false)
		if err != nil {
			return err
		}
		for moniker, containerID := range cons {
			common.DebugMessage("Stopping Container " + moniker)
			err = StopContainer(cli, containerID)
			if err != nil {
				return err
			}

			if force {
				common.DebugMessage("Removing Container " + containerID)
				err = RemoveContainer(cli, containerID, true, false, false)
				if err != nil {
					return err
				}
			}
		}

		if err := RemoveNetwork(cli, netID); err != nil {
			return err
		}
		common.DebugMessage("")
		nets, err = GetNetworks(cli, false)
		if err != nil {
			return err
		}

		common.DebugMessage("\nNumber of Networks: " + strconv.Itoa(len(nets)))
	} else {
		return errors.New("network not found")
	}

	return nil
}

//StopNode stops a node. If force is set the container is removed too.
func StopNode(networkName, nodeName string, force bool) error {

	cli, err := GetDockerClient()
	if err != nil {
		return err
	}

	cons, err := GetContainers(cli, false)
	if err != nil {
		return err
	}

	containerID, ok := cons[nodeName]
	if !ok {
		return errors.New("container " + nodeName + " is not running")
	}

	common.DebugMessage("Stopping Container " + containerID)
	err = StopContainer(cli, containerID)
	if err != nil {
		return err
	}

	if force {
		common.DebugMessage("Removing Container " + containerID)
		err = RemoveContainer(cli, containerID, true, false, false)
		if err != nil {
			return err
		}
	}

	return nil
}
