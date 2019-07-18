package commands

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/mosaicnetworks/monetd/src/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// _config and _logger are global variables
	_config = monetConfig(defaultHomeDir())
	_logger = defaultLogger()

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

		_logger.Level = logLevel(_config.LogLevel)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(
		NewRunCmd(),
		VersionCmd,
	)

	// set global flags
	RootCmd.PersistentFlags().StringP("datadir", "d", _config.DataDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().String("log", _config.LogLevel, "debug, info, warn, error, fatal, panic")

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
	_config = config.DefaultConfig()

	// first unmarshal to read from cli flags
	if err := viper.Unmarshal(_config); err != nil {
		return err
	}

	// Trickle-down datadir config to sub-config sections (Babble and Eth). Only
	// effective if config.DataDir is currently equal to the monet default
	// (~/.monet).
	_config.SetDataDir(_config.DataDir)

	// Read from configuration file if there is one.
	// ATTENTION: CLI flags will always have precedence of these values.

	viper.SetConfigName("monetd")        // name of config file (without extension)
	viper.AddConfigPath(_config.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// fmt.Printf("No config file monetd.toml found in %s\n", config.DataDir)
	} else {
		return err
	}

	// second unmarshal to read from config file
	return viper.Unmarshal(_config)
}

// default config for monetd
func monetConfig(dataDir string) *config.Config {
	config := config.DefaultConfig()

	config.SetDataDir(dataDir)

	return config
}

// default logger (debug)
func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	return logger
}

func defaultHomeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "MONET")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "MONET")
		} else {
			return filepath.Join(home, ".monet")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func logLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
