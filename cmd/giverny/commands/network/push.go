package network

import (
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
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

func addPushFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkPush(cmd *cobra.Command, args []string) error {
	networkName := args[0]
	nodeName := args[1]
	return buildZip(configuration.GivernyConfigDir, networkName, nodeName)

}
