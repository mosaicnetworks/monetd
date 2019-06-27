package network

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show configuration",
		Long: `
Show configuration.`,
		Args: cobra.ExactArgs(0),
		RunE: showConfig,
	}

	return cmd
}

func showConfig(cmd *cobra.Command, args []string) error {
	filename := filepath.Join(configDir, tomlName+".toml")
	message("Displaying file: ", filename)

	return common.ShowConfigFile(filename)
}
