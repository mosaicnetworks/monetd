package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "create the configuration for a multi-node network",
		Long: `
giverny network build
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkBuild,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkBuild(cmd *cobra.Command, args []string) error {
	return nil
}
