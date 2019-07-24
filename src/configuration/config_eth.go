package configuration

import "fmt"

var (
	defaultEthAPIAddr   = ":8080"
	defaultCache        = 128
	defaultEthDir       = fmt.Sprintf("%s/eth", defaultDataDir)
	defaultKeystoreFile = fmt.Sprintf("%s/keystore", defaultEthDir)
	defaultGenesisFile  = fmt.Sprintf("%s/genesis.json", defaultEthDir)
	defaultPwdFile      = fmt.Sprintf("%s/pwd.txt", defaultEthDir)
	defaultDbFile       = fmt.Sprintf("%s/chaindata", defaultEthDir)
)

// EthConfig contains the configuration relative to the accounts, EVM, trie/db,
// and service API
type EthConfig struct {

	// Genesis file
	Genesis string `mapstructure:"genesis"`

	// File containing the levelDB database
	DbFile string `mapstructure:"db"`

	// Address of HTTP API Service
	EthAPIAddr string `mapstructure:"listen"`

	// Megabytes of memory allocated to internal caching (min 16MB / database forced)
	Cache int `mapstructure:"cache"`
}

// DefaultEthConfig return the default configuration for Eth services
func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		Genesis:    defaultGenesisFile,
		DbFile:     defaultDbFile,
		EthAPIAddr: defaultEthAPIAddr,
		Cache:      defaultCache,
	}
}

// SetDataDir updates the eth configuration directories if they were set to
// default values.
func (c *EthConfig) SetDataDir(datadir string) {
	if c.Genesis == defaultGenesisFile {
		c.Genesis = fmt.Sprintf("%s/genesis.json", datadir)
	}
	if c.DbFile == defaultDbFile {
		c.DbFile = fmt.Sprintf("%s/chaindata", datadir)
	}
}
