package configuration

import (
	"fmt"
	"time"
)

var (
	defaultNodeAddr       = "127.0.0.1:1337"
	defaultBabbleAPIAddr  = ":8000"
	defaultHeartbeat      = 500 * time.Millisecond
	defaultTCPTimeout     = 1000 * time.Millisecond
	defaultCacheSize      = 50000
	defaultSyncLimit      = 1000
	defaultEnableFastSync = false
	defaultMaxPool        = 2
	defaultBabbleDir      = fmt.Sprintf("%s/babble", defaultDataDir)
	defaultPeersFile      = fmt.Sprintf("%s/peers.json", defaultBabbleDir)
)

// BabbleConfig contains the configuration for the Babble node used by monetd.
// It only presents a subset of the options Babble can accept, because monetd
// forces some configurations values. In particular, the --fast-sync and
// --store flags are disabled because monetd does not support the FastSync
// protocol, and it requires a persistant database.
type BabbleConfig struct {

	// Directory containing priv_key.pem and peers.json files
	DataDir string `mapstructure:"datadir"`

	// Address of Babble node (where it talks to other Babble nodes)
	BindAddr string `mapstructure:"listen"`

	// Babble HTTP API address
	ServiceAddr string `mapstructure:"service-listen"`

	// Gossip heartbeat
	Heartbeat time.Duration `mapstructure:"heartbeat"`

	// TCP timeout
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// Max number of items in caches
	CacheSize int `mapstructure:"cache-size"`

	// Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync-limit"`

	// Max number of connections in net pool
	MaxPool int `mapstructure:"max-pool"`

	// Bootstrap from database
	Bootstrap bool `mapstructure:"bootstrap"`
}

// DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultBabbleConfig() *BabbleConfig {
	return &BabbleConfig{
		DataDir:     defaultBabbleDir,
		BindAddr:    defaultNodeAddr,
		ServiceAddr: defaultBabbleAPIAddr,
		Heartbeat:   defaultHeartbeat,
		TCPTimeout:  defaultTCPTimeout,
		CacheSize:   defaultCacheSize,
		SyncLimit:   defaultSyncLimit,
		MaxPool:     defaultMaxPool,
	}
}

// SetDataDir updates the babble configuration directories if they were set to
// to default values.
func (c *BabbleConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultBabbleDir {
		c.DataDir = datadir
	}
}
