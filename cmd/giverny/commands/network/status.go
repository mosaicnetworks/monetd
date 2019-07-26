package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "shows network status",
		Long: `
giverny network status
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStatus,
	}

	addStatusFlags(cmd)

	return cmd
}

func addStatusFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStatus(cmd *cobra.Command, args []string) error {
	return nil
}
