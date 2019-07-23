package commands

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/spf13/cobra"

	"github.com/mosaicnetworks/monetd/cmd/monetd/config"
)

var (
	//ConfigCmd implements the config CLI subcommand
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "manage monetd configuration",
		Long: `monetd config
		
The config subcommands manage the monet configuration file, as used by 
the monetd server process. `,
		TraverseChildren: true,
	}
)

func init() {
	//Subcommands
	configCmd.AddCommand(
		//			NewCheckCmd(),
		//			NewPublishCmd(),
		newLocationCmd(),
	//			NewShowCmd(),
	//			NewClearCmd(),
	//			NewPullCmd(),
	//			NewBuildCmd(),
	)

	//TODO remove - temporary debug out to preserve the import - we will need it shortly
	fmt.Print(config.Config.DataDir)

	// datadir is now the config for everything...

	//	viper.BindPFlags(ConfigCmd.Flags())
}

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `monetd config location
Shows the location of the configuration files for the monetd server.`,
		Args: cobra.ArbitraryArgs,
		RunE: locationConfig,
	}
	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {
	common.InfoMessage("The Monet Configuration files are located at:")
	common.InfoMessage(config.Config.DataDir)
	return nil
}
