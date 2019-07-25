// Package configuration holds shared configuration structs for Monet, EVM-Lite and Babble.

package configuration

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/mosaicnetworks/babble/src/babble"
	evml_config "github.com/mosaicnetworks/evm-lite/src/config"
)

var (
	// Base
	defaultLogLevel   = "debug"
	defaultDataDir, _ = common.DefaultMonetConfigDir()

	// Global is a global Config object that is used by commands in cmd/ to
	// manipulate configuration options.
	Global = monetConfig(defaultDataDir)
)

// default config for monetd
func monetConfig(dataDir string) *Config {
	_config := DefaultConfig()

	if dataDir == "" { // fall through to default settings
		file, _ := common.DefaultMonetConfigDir()
		_config.SetDataDir(file)
	} else {
		_config.SetDataDir(dataDir)
	}
	return _config
}

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

// SetDataDir updates the root data directory as well as the config for EVM-Lite
// and Babble
func (c *Config) SetDataDir(datadir string) {
	c.BaseConfig.DataDir = datadir
	if c.Eth != nil {
		c.Eth.SetDataDir(fmt.Sprintf("%s/eth", datadir))
	}
	if c.Babble != nil {
		c.Babble.SetDataDir(fmt.Sprintf("%s/babble", datadir))
	}
}

// ToEVMLConfig extracts evm-lite configuration and returns a config object as
// used by the evm-lite library.
func (c *Config) ToEVMLConfig() *evml_config.Config {
	evmlConfig := evml_config.DefaultConfig()

	evmlConfig.DataDir = c.DataDir
	evmlConfig.LogLevel = c.LogLevel
	evmlConfig.Genesis = c.Eth.Genesis
	evmlConfig.DbFile = c.Eth.DbFile
	evmlConfig.EthAPIAddr = c.Eth.EthAPIAddr
	evmlConfig.Cache = c.Eth.Cache

	return evmlConfig
}

// ToBabbleConfig extracts the babble configuration and returns a config object
// as used by the Babble library. It enforces the values of Store and
// EnableFastSync to true and false respectively.
func (c *Config) ToBabbleConfig() *babble.BabbleConfig {
	babbleConfig := babble.NewDefaultConfig()

	babbleConfig.DataDir = c.Babble.DataDir
	babbleConfig.BindAddr = c.Babble.BindAddr
	babbleConfig.ServiceAddr = c.Babble.ServiceAddr
	babbleConfig.MaxPool = c.Babble.MaxPool

	babbleConfig.NodeConfig.HeartbeatTimeout = c.Babble.Heartbeat
	babbleConfig.NodeConfig.TCPTimeout = c.Babble.TCPTimeout
	babbleConfig.NodeConfig.CacheSize = c.Babble.CacheSize
	babbleConfig.NodeConfig.SyncLimit = c.Babble.SyncLimit
	babbleConfig.NodeConfig.Bootstrap = c.Babble.Bootstrap

	babbleConfig.Store = true
	babbleConfig.NodeConfig.EnableFastSync = false

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
}

// DefaultBaseConfig returns the default top-level configuration for EVM-Babble
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  defaultDataDir,
		LogLevel: defaultLogLevel,
	}
}
