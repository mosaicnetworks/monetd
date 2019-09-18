package configuration

import (
	"fmt"
	"time"

	"github.com/mosaicnetworks/monetd/src/common"
)

var (
	defaultNodeAddr       = fmt.Sprintf("%s:%d", common.GetMyIP(), 1337)
	defaultHeartbeat      = 200 * time.Millisecond
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

	// BindAddr is the local address:port where this node gossips with other
	// nodes. This is an IP address that should be reachable by all other nodes
	// in the cluster. By default, this is "0.0.0.0", meaning Babble will bind
	// to all addresses on the local machine and will advertise the private IPv4
	// address to the rest of the cluster. However, in some cases, there may be
	// a routable address that cannot be bound. Use AdvertiseAddr to enable
	// gossiping a different address to support this. If this address is not
	// routable, the node will be in a constant flapping state as other nodes
	// will treat the non-routability as a failure
	BindAddr string `mapstructure:"listen"`

	// AdvertiseAddr is used to change the address that we advertise to other
	// nodes in the cluster
	AdvertiseAddr string `mapstructure:"advertise"`

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

	// Moniker is a friendly name to indentify this peer
	Moniker string `mapstructure:"moniker"`
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
