package config

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
)

// newClearCmd returns the clear command which creates a backup of the current
// configuration folder, before clearing it completely.
func newClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "backup and clear configuration folder",
		Long: `
Backup and delete the current configuration folder ([datadir]).
`,
		RunE: clearConfig,
	}
	return cmd
}

func clearConfig(cmd *cobra.Command, args []string) error {
	return files.SafeRename(configuration.DefaultMonetDir())
}
