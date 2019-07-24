package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mosaicnetworks/monetd/src/configuration"
	pconfig "github.com/mosaicnetworks/monetd/src/poa/config"
)

// newBuildCmd initialises a bare-bones configuration for monetd
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build the configuration files",
		Long: `
monetd config build

Builds a bare-bones configuration for monetd`,
		RunE: buildConfig,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&nodeParam, "node", "n", "", "the directory name containing monet nodes configurations")
	cmd.Flags().StringVarP(&addressParam, "address", "a", "", " ip address/host name of this node")
	cmd.Flags().StringVarP(&passwordFile, "passfile", "p", "", "the file that contains the passphrase for the keyfile")
	viper.BindPFlags(cmd.Flags())
}

func buildConfig(cmd *cobra.Command, args []string) error {
	return pconfig.BuildConfig(configuration.Configuration.DataDir, nodeParam, addressParam, passwordFile)
}
