package config

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/configuration"
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
	fmt.Println(configuration.Configuration.DataDir)
	return nil
}
