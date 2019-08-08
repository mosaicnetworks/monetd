package configuration

var (
	defaultCache = 128
)

// EthConfig contains the configuration relative to the accounts, EVM, trie/db,
// and service API
type EthConfig struct {
	// Megabytes of memory allocated to internal caching (min 16MB / database forced)
	Cache int `mapstructure:"cache"`
}

// DefaultEthConfig return the default configuration for Eth services
func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		Cache: defaultCache,
	}
}
