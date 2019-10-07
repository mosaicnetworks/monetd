package network

import (
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

If the <node> value is provided, this command stops just that node. Otherwise it
stops all the nodes. Additionaly, if the --remove flag is set, the nodes are 
also deleted.
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
		return docker.StopNetwork(args[0], stopAndDelete)
	}

	return docker.StopNode(args[0], args[1], stopAndDelete)

}
