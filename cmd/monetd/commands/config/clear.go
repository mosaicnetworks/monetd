package config

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/spf13/cobra"
)

// newClearCmd returns the clear command which creates a backup of the current
// configuration folder, before clearing it completely.
func newClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "backup and clear configuration folder",
		Long: `
The clear subcommand creates a backup of the current configuration folder 
([datadir]) before deleting it.
`,
		RunE: clearConfig,
	}
	return cmd
}

func clearConfig(cmd *cobra.Command, args []string) error {
	if files.CheckIfExists(configuration.Global.DataDir) {
		files.SafeRenameDir(configuration.Global.DataDir)
	}
	return nil
}
