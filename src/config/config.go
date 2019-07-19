package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/mosaicnetworks/babble/src/babble"
	evml_config "github.com/mosaicnetworks/evm-lite/src/config"
)

var (
	// Base
	defaultLogLevel = "debug"

	// DefaultDataDir is the default root directory for monet configuration and
	// data files. By default, evm-lite and babble configuration will also be
	// within this directory.
	DefaultDataDir = defaultHomeDir()
)

// Config contains de configuration for MONET node
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
		DataDir:  DefaultDataDir,
		LogLevel: defaultLogLevel,
	}
}

/*******************************************************************************
FILE HELPERS
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
