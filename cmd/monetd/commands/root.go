//Package commands implements the CLI commands for monetd
package commands

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/monetd/cmd/monetd/commands/keys"
	"github.com/mosaicnetworks/monetd/cmd/monetd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//	config = monetConfig(defaultHomeDir())
	logger = common.DefaultLogger()

	passwordFile string
	outputJSON   bool
)

/*******************************************************************************
RootCmd
*******************************************************************************/

//RootCmd is the root command for monetd
var RootCmd = &cobra.Command{
	Use:              "monetd",
	Short:            "MONET-Daemon",
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := readConfig(cmd); err != nil {
			return err
		}

		logger.Level = common.LogLevel(config.Config.LogLevel)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(
		//		InitCmd,

		keys.KeysCmd,
		newRunCmd(),
		versionCmd,
		configCmd,
	)

	// set global flags
	RootCmd.PersistentFlags().StringP("datadir", "d", config.Config.DataDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().String("log", config.Config.LogLevel, "debug, info, warn, error, fatal, panic")

	// do not print usage when error occurs
	RootCmd.SilenceUsage = true
}

/*******************************************************************************
HELPERS
*******************************************************************************/

// read config into Viper
func readConfig(cmd *cobra.Command) error {

	// Register flags with viper
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// Reset config because evm-lite's SetDataDir only updates values if they
	// are currently equal to the defaults (~/.evm-lite/*). Before this call,
	// they should be set to monetd defaults (.monetd/*).
	config.Config = _config.DefaultConfig()

	// first unmarshal to read from cli flags
	if err := viper.Unmarshal(config.Config); err != nil {
		return err
	}

	// EnableFastSync and Store are not configurable, they MUST have these
	// values:
	config.Config.Babble.EnableFastSync = false
	config.Config.Babble.Store = true

	// Trickle-down datadir config to sub-config sections (Babble and Eth). Only
	// effective if config.Config.DataDir is currently equal to the evm-lite default
	// (~/.evm-lite).
	config.Config.SetDataDir(config.Config.DataDir)

	// Read from configuration file if there is one.
	// ATTENTION: CLI flags will always have precedence of these values.

	viper.SetConfigName("monetd")              // name of config file (without extension)
	viper.AddConfigPath(config.Config.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// fmt.Printf("No config file monetd.toml found in %s\n", config.Config.DataDir)
	} else {
		return err
	}

	// second unmarshal to read from config file
	return viper.Unmarshal(config.Config)
}
