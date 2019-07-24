package config

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/spf13/cobra"
)

//newClearCmd shows the config file path
func newClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "clears any pre-existing configuration files",
		Long: `
monetd config clear

Cleans up configuration files for the monetd server by 
renaming any pre-existing configuration.
`,
		RunE: clearConfig,
	}
	return cmd
}

func clearConfig(cmd *cobra.Command, args []string) error {
	if files.CheckIfExists(configuration.Configuration.DataDir) {
		files.SafeRenameDir(configuration.Configuration.DataDir)
	}
	return nil
}
