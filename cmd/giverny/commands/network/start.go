package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start a multi-node network",
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
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {

	return nil
}
