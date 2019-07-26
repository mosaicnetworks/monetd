package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop a network",
		Long: `
giverny network stop

Stop a node and all the nodes within it.
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStop,
	}

	addStopFlags(cmd)

	return cmd
}

func addStopFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStop(cmd *cobra.Command, args []string) error {

	return nil
}
