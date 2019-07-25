package configuration

import (
	"fmt"
	"time"

	"github.com/mosaicnetworks/monetd/src/common"
)

var (
	defaultNodeAddr       = fmt.Sprintf("%s:%d", common.GetMyIP(), 1337)
	defaultHeartbeat      = 500 * time.Millisecond
	defaultTCPTimeout     = 1000 * time.Millisecond
	defaultCacheSize      = 50000
	defaultSyncLimit      = 1000
	defaultEnableFastSync = false
	defaultMaxPool        = 2
)

// BabbleConfig contains the configuration for the Babble node used by monetd.
// It only presents a subset of the options Babble can accept, because monetd
// forces some configurations values. In particular, the --fast-sync and
// --store flags are disabled because monetd does not support the FastSync
// protocol, and it requires a persistant database.
type BabbleConfig struct {

	// Address of Babble node (where it talks to other Babble nodes)
	BindAddr string `mapstructure:"listen"`

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
		BindAddr:   defaultNodeAddr,
		Heartbeat:  defaultHeartbeat,
		TCPTimeout: defaultTCPTimeout,
		CacheSize:  defaultCacheSize,
		SyncLimit:  defaultSyncLimit,
		MaxPool:    defaultMaxPool,
	}
}
