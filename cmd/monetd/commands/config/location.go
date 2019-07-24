package config

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/spf13/cobra"
)

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		RunE:  locationConfig,
	}
	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {
	common.InfoMessage("The Monet Configuration files are located at:")
	// XXX this should use the [datadir] inherited from the root command
	common.InfoMessage(configuration.Configuration.DataDir)
	return nil
}
