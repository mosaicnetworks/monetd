package network

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "shows network status",
		Long: `
giverny network status
		`,
		Args: cobra.ArbitraryArgs,
		RunE: networkStatus,
	}

	return cmd
}

func networkStatus(cmd *cobra.Command, args []string) error {

	common.DebugMessage("Connecting to Docker Client")

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	common.InfoMessage("\n\nNetworks\n")

	nets, err := docker.GetNetworks(cli, true)
	if err != nil {
		return err
	}

	if len(nets) == 0 {
		common.ErrorMessage("No networks found")
	}

	common.InfoMessage("\n\nContainers\n")
	containers, err := docker.GetContainers(cli, true)

	if len(containers) == 0 {
		common.ErrorMessage("No containers found")
	}

	return nil
}
