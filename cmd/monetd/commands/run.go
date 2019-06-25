package commands

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/consensus/babble"
	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config = monetConfig(defaultHomeDir())
	logger = defaultLogger()
)

func monetConfig(dataDir string) *_config.Config {
	config := _config.DefaultConfig()

	config.Babble.EnableFastSync = false
	config.Babble.Store = true

	config.SetDataDir(dataDir)

	return config
}

func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	return logger
}

/*******************************************************************************
RunCmd
*******************************************************************************/

//RunCmd is launches a node
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a MONET node",
	Long: `Run a MONET node.

  Start a daemon which acts as a full node on a MONET network. All data and 
  configuration are stored under a directory [datadir] controlled by the 
  --datadir flag ($HOME/.monet by default on UNIX systems). 
  
  [datadir] must contain a set of files defining the network that this node is 
  attempting to join or create. Please refer to monetcli for a tool to manage 
  this configuration. 
  
  Further options pertaining to the operation of the node are read from the 
  [datadir]/monetd.toml file, or overwritten by the following flags.`,

	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := readConfig(cmd); err != nil {
			return err
		}

		logger.Level = logLevel(config.LogLevel)

		logger.WithField("Base", fmt.Sprintf("%+v", config.BaseConfig)).Debug("Config Base")
		logger.WithField("Babble", fmt.Sprintf("%+v", config.Babble)).Debug("Config Babble")
		logger.WithField("Eth", fmt.Sprintf("%+v", config.Eth)).Debug("Config Eth")

		return nil
	},
	RunE: runBabble,
}

func init() {
	// Base config
	RunCmd.Flags().StringP("datadir", "d", config.DataDir, "Top-level directory for configuration and data")

	// Babble config
	RunCmd.Flags().String("babble.listen", config.Babble.BindAddr, "IP:PORT of Babble node")
	RunCmd.Flags().String("babble.service-listen", config.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	RunCmd.Flags().Duration("babble.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	RunCmd.Flags().Duration("babble.timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	RunCmd.Flags().Int("babble.cache-size", config.Babble.CacheSize, "Number of items in LRU caches")
	RunCmd.Flags().Int("babble.sync-limit", config.Babble.SyncLimit, "Max number of Events per sync")
	RunCmd.Flags().Int("babble.max-pool", config.Babble.MaxPool, "Max number of pool connections")

	// Eth config
	RunCmd.PersistentFlags().String("eth.listen", config.Eth.EthAPIAddr, "IP:PORT of Monet HTTP API service")
	RunCmd.PersistentFlags().Int("eth.cache", config.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")

	viper.BindPFlags(RunCmd.Flags())
}

/*******************************************************************************
READ CONFIG AND RUN
*******************************************************************************/

// Read Config into Viper
func readConfig(cmd *cobra.Command) error {
	config = _config.DefaultConfig()

	// unmarshal a first time to read from cli flags
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	// EnableFastSync and Store are not configurable, they MUST have these
	// values:
	config.Babble.EnableFastSync = false
	config.Babble.Store = true

	config.SetDataDir(config.DataDir)

	viper.SetConfigName("monetd")       // name of config file (without extension)
	viper.AddConfigPath(config.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		logger.Debugf("No config file monetd.toml found in %s", config.DataDir)
	} else {
		return err
	}

	// unmarshal a second time to read from config file
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// Run the EVM-Lite / Babble engine
func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(config.Babble, logger)
	engine, err := engine.NewEngine(*config, babble, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}

/*******************************************************************************
HELPERS
*******************************************************************************/

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
