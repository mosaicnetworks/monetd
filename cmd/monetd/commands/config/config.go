package config

import "github.com/spf13/cobra"

var (
	nodeParam    string
	addressParam string
	passwordFile string
	existingPeer string
)

// ConfigCmd implements the config CLI subcommand
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage monetd configuration",
	Long: `
monetd config
		
The config subcommands manage the monet configuration file, as used by the 
monetd server process. `,
	TraverseChildren: true,
}

func init() {
	// Subcommands
	ConfigCmd.AddCommand(
		newLocationCmd(),
		newClearCmd(),
		newPullCmd(),
		newBuildCmd(),
		newContractCmd(),
	)
}
