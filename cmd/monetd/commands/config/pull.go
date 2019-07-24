package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "pull the configuration files from a node",
		RunE:  pullConfig,
	}

	addBuildFlags(cmd)

	return cmd
}

func addPullFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&nodeParam, "node", "n", "", "the directory name containing monet nodes configurations")
	cmd.Flags().StringVarP(&addressParam, "address", "a", "", " ip address/host name of this node")
	cmd.Flags().StringVarP(&passwordFile, "passfile", "p", "", "the file that contains the passphrase for the keyfile")
	cmd.Flags().StringVar(&existingPeer, "peer", "", "the address of an existing peer")
	viper.BindPFlags(cmd.Flags())
}

func pullConfig(cmd *cobra.Command, args []string) error {
	return nil
}
