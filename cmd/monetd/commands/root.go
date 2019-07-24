//Package commands implements the CLI commands for monetd
package commands

import (
	"fmt"

<<<<<<< HEAD
	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/mosaicnetworks/monetd/cmd/monetd/commands/keys"
	"github.com/mosaicnetworks/monetd/src/configuration"
=======
	"github.com/mosaicnetworks/monetd/src/config"
	"github.com/sirupsen/logrus"
>>>>>>> origin/develop
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	_logger = common.DefaultLogger()

	passwordFile string
	outputJSON   bool
	nodeParam    string
	addressParam string
	existingPeer string
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

		_logger.Level = common.LogLevel(configuration.Configuration.LogLevel)

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
	RootCmd.PersistentFlags().StringP("datadir", "d", configuration.Configuration.DataDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().String("log", configuration.Configuration.LogLevel, "debug, info, warn, error, fatal, panic")

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
	// they should be set to monetd defaults (.monet/*).
	configuration.Configuration = configuration.DefaultConfig()

	// first unmarshal to read from cli flags
	if err := viper.Unmarshal(configuration.Configuration); err != nil {
		return err
	}

	// Trickle-down datadir config to sub-config sections (Babble and Eth). Only
	// effective if _config.DataDir is currently equal to the monet default
	// (~/.monet) on Linux.
	configuration.Configuration.SetDataDir(configuration.Configuration.DataDir)

	// Read from configuration file if there is one.
	// ATTENTION: CLI flags will always have precedence of these values.

	viper.SetConfigName("monetd")                            // name of config file (without extension)
	viper.AddConfigPath(configuration.Configuration.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// fmt.Printf("No config file monetd.toml found in %s\n", _config.DataDir)
	} else {
		return err
	}

	// second unmarshal to read from config file
	return viper.Unmarshal(configuration.Configuration)
}

// default config for monetd
func monetConfig(dataDir string) *configuration.Config {
	config := configuration.DefaultConfig()

	config.SetDataDir(dataDir)

	return config
}
