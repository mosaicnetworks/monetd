// Package configuration holds shared configuration structs for Monet, EVM-Lite and Babble.
package configuration

import (
	"fmt"

	babble_config "github.com/mosaicnetworks/babble/src/config"
	evml_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	// Base
	defaultDataDir, _ = DefaultMonetConfigDir()

	// Global is a global Config object used by commands in cmd/ to manipulate
	// configuration options.
	Global = DefaultConfig()
)

// Config contains the configuration for MONET node
type Config struct {

	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for evm-lite
	Eth *EthConfig `mapstructure:"eth"`

	// Options for Babble
	Babble *BabbleConfig `mapstructure:"babble"`
}

// DefaultConfig returns the default configuration for a MONET node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Eth:        DefaultEthConfig(),
		Babble:     DefaultBabbleConfig(),
	}
}

// LogLevel returns a logrus-style log-level based on the verbose option.
func (c *Config) LogLevel() string {
	if c.Verbose {
		return "debug"
	}
	return "info"
}

// Logger returns a new prefixed logrus Entry with custom formatting
func (c *Config) Logger(prefix string) *logrus.Entry {
	if c.logger == nil {
		c.logger = logrus.New()
		c.logger.Level = common.LogLevel(c.LogLevel())
		c.logger.Formatter = new(prefixed.TextFormatter)
	}
	return c.logger.WithField("prefix", prefix)
}

// ToEVMLConfig extracts evm-lite configuration and returns a config object as
// used by the evm-lite library.
func (c *Config) ToEVMLConfig() *evml_config.Config {
	evmlConfig := evml_config.DefaultConfig()

	evmlConfig.DataDir = c.DataDir
	evmlConfig.LogLevel = c.LogLevel()
	evmlConfig.EthAPIAddr = c.APIAddr
	evmlConfig.Genesis = fmt.Sprintf("%s/%s/%s", c.DataDir, EthDir, GenesisJSON)
	evmlConfig.DbFile = fmt.Sprintf("%s/%s/%s", c.DataDir, EthDir, Chaindata)
	evmlConfig.Cache = c.Eth.Cache
	evmlConfig.MinGasPrice = c.Eth.MinGasPrice

	return evmlConfig
}

// ToBabbleConfig extracts the babble configuration and returns a config object
// as used by the Babble library. It enforces the values of Store and
// EnableFastSync to true and false respectively.
func (c *Config) ToBabbleConfig() *babble_config.Config {
	babbleConfig := babble_config.NewDefaultConfig()

	babbleConfig.DataDir = fmt.Sprintf("%s/%s", c.DataDir, BabbleDir)
	babbleConfig.LogLevel = c.LogLevel()
	babbleConfig.BindAddr = c.Babble.BindAddr
	babbleConfig.MaxPool = c.Babble.MaxPool
	babbleConfig.HeartbeatTimeout = c.Babble.Heartbeat
	babbleConfig.TCPTimeout = c.Babble.TCPTimeout
	babbleConfig.CacheSize = c.Babble.CacheSize
	babbleConfig.SyncLimit = c.Babble.SyncLimit
	babbleConfig.Bootstrap = c.Babble.Bootstrap
	babbleConfig.Moniker = c.Babble.Moniker

	// Force Babble to use persistant storage.
	babbleConfig.Store = true

	// Force FastSync = false because EVM-Lite does not support Snapshot/Restore
	// yet.
	babbleConfig.EnableFastSync = false

	// An empty ServiceAddr tells Babble not to start an API server. The API
	// handlers are still registered with the DefaultServeMux, so they will be
	// served by the EVM-Lite API server automatically.
	babbleConfig.ServiceAddr = ""

	return babbleConfig
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

// BaseConfig contains the top level configuration for an EVM-Babble node
type BaseConfig struct {

	// Top-level directory of evm-babble data
	DataDir string `mapstructure:"datadir"`

	// Verbose
	Verbose bool `mapstructure:"verbose"`

	// IP/PORT of API
	APIAddr string `mapstructure:"api-listen"`

	logger *logrus.Logger
}

// DefaultBaseConfig returns the default top-level configuration for EVM-Babble
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir: defaultDataDir,
		APIAddr: DefaultAPIAddr,
	}
}
