package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push the configuration for a node on the named network",
		Long: `
giverny network push
		`,
		Args: cobra.ExactArgs(1),
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

	return nil
}
