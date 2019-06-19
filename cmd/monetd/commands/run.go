package commands

import (
	"fmt"

	_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/consensus/babble"
	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config = _config.DefaultConfig()
	logger = logrus.New()
)

//AddBabbleFlags adds flags to the Babble command
func AddBabbleFlags(cmd *cobra.Command) {
	cmd.Flags().String("babble.datadir", config.Babble.DataDir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("babble.listen", config.Babble.BindAddr, "IP:PORT of Babble node")
	cmd.Flags().String("babble.service-listen", config.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Duration("babble.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("babble.timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache-size", config.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync-limit", config.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Bool("babble.enable-fast-sync", config.Babble.EnableFastSync, "Enable FastSync")
	cmd.Flags().Int("babble.max-pool", config.Babble.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("babble.store", config.Babble.Store, "use persistent store")
	viper.BindPFlags(cmd.Flags())
}

//RunCmd is launches a node
var RunCmd = &cobra.Command{
	Use:              "run",
	Short:            "Run a Monet node",
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := bindFlagsLoadViper(cmd); err != nil {
			return err
		}

		config, err = parseConfig()
		if err != nil {
			return err
		}

		logger = logrus.New()
		logger.Level = logLevel(config.BaseConfig.LogLevel)

		config.SetDataDir(config.BaseConfig.DataDir)

		logger.WithFields(logrus.Fields{
			"Base": config.BaseConfig,
			"Eth":  config.Eth}).Debug("Config")

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		config.SetDataDir(config.BaseConfig.DataDir)

		logger.WithFields(logrus.Fields{
			"Babble": config.Babble,
		}).Debug("Config")

		return nil
	},
	RunE: runBabble,
}

func init() {
	//Subcommands
	//	RunCmd.AddCommand(
	//		NewBabbleCmd(),
	//		NewRaftCmd(),
	//		NewSoloCmd())

	//Base config
	RunCmd.Flags().StringP("datadir", "d", config.BaseConfig.DataDir, "Top-level directory for configuration and data")
	RunCmd.Flags().String("log", config.BaseConfig.LogLevel, "debug, info, warn, error, fatal, panic")

	//Eth config
	RunCmd.Flags().String("eth.genesis", config.Eth.Genesis, "Location of genesis file")
	RunCmd.Flags().String("eth.keystore", config.Eth.Keystore, "Location of Ethereum account keys")
	RunCmd.Flags().String("eth.pwd", config.Eth.PwdFile, "Password file to unlock accounts")
	RunCmd.Flags().String("eth.db", config.Eth.DbFile, "Eth database file")
	RunCmd.Flags().String("eth.listen", config.Eth.EthAPIAddr, "Address of HTTP API service")
	RunCmd.Flags().Int("eth.cache", config.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")
	//Babble config
	RunCmd.Flags().String("babble.datadir", config.Babble.DataDir, "Directory contaning priv_key.pem and peers.json files")
	RunCmd.Flags().String("babble.listen", config.Babble.BindAddr, "IP:PORT of Babble node")
	RunCmd.Flags().String("babble.service-listen", config.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	RunCmd.Flags().Duration("babble.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	RunCmd.Flags().Duration("babble.timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	RunCmd.Flags().Int("babble.cache-size", config.Babble.CacheSize, "Number of items in LRU caches")
	RunCmd.Flags().Int("babble.sync-limit", config.Babble.SyncLimit, "Max number of Events per sync")
	RunCmd.Flags().Bool("babble.enable-fast-sync", config.Babble.EnableFastSync, "Enable FastSync")
	RunCmd.Flags().Int("babble.max-pool", config.Babble.MaxPool, "Max number of pool connections")
	RunCmd.Flags().Bool("babble.store", config.Babble.Store, "use persistent store")
	viper.BindPFlags(RunCmd.Flags())
}

//------------------------------------------------------------------------------

//Retrieve the default environment configuration.
func parseConfig() (*_config.Config, error) {
	conf := _config.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

//Bind all flags and read the config into viper
func bindFlagsLoadViper(cmd *cobra.Command) error {
	// cmd.Flags() includes flags from this command and all persistent flags from the parent
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	viper.SetConfigName("evml")                    // name of config file (without extension)
	viper.AddConfigPath(config.BaseConfig.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		logger.Debugf("No config file found in %s", config.DataDir)
	} else {

		return err
	}

	return nil
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

func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(config.Babble, logger)
	engine, err := engine.NewEngine(*config, babble, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
