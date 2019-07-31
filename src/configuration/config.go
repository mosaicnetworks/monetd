// Package configuration holds shared configuration structs for Monet, EVM-Lite and Babble.
package configuration

import (
	"fmt"

	"github.com/mosaicnetworks/babble/src/babble"
	evml_config "github.com/mosaicnetworks/evm-lite/src/config"
)

var (
	// Base
	defaultLogLevel   = "debug"
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

// ToEVMLConfig extracts evm-lite configuration and returns a config object as
// used by the evm-lite library.
func (c *Config) ToEVMLConfig() *evml_config.Config {
	evmlConfig := evml_config.DefaultConfig()

	evmlConfig.DataDir = c.DataDir
	evmlConfig.LogLevel = c.LogLevel
	evmlConfig.EthAPIAddr = c.APIAddr
	evmlConfig.Genesis = fmt.Sprintf("%s/%s/%s", c.DataDir, EthDir, GenesisJSON)
	evmlConfig.DbFile = fmt.Sprintf("%s/%s/%s", c.DataDir, EthDir, Chaindata)
	evmlConfig.Cache = c.Eth.Cache

	return evmlConfig
}

// ToBabbleConfig extracts the babble configuration and returns a config object
// as used by the Babble library. It enforces the values of Store and
// EnableFastSync to true and false respectively.
func (c *Config) ToBabbleConfig() *babble.BabbleConfig {
	babbleConfig := babble.NewDefaultConfig()

	babbleConfig.DataDir = fmt.Sprintf("%s/%s", c.DataDir, BabbleDir)
	babbleConfig.BindAddr = c.Babble.BindAddr
	babbleConfig.MaxPool = c.Babble.MaxPool
	babbleConfig.NodeConfig.HeartbeatTimeout = c.Babble.Heartbeat
	babbleConfig.NodeConfig.TCPTimeout = c.Babble.TCPTimeout
	babbleConfig.NodeConfig.CacheSize = c.Babble.CacheSize
	babbleConfig.NodeConfig.SyncLimit = c.Babble.SyncLimit
	babbleConfig.NodeConfig.Bootstrap = c.Babble.Bootstrap

	// Force Babble to use persistant storage.
	babbleConfig.Store = true

	// Force FastSync = false because EVM-Lite does not support Snapshot/Restore
	// yet.
	babbleConfig.NodeConfig.EnableFastSync = false

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

	// Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log"`

	// IP/PORT of API
	APIAddr string `mapstructure:"api-listen"`
}

// DefaultBaseConfig returns the default top-level configuration for EVM-Babble
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  defaultDataDir,
		LogLevel: defaultLogLevel,
		APIAddr:  DefaultAPIAddr,
	}
}
