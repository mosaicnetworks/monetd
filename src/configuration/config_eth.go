package configuration

var (
	defaultCache       = 128
	defaultMinGasPrice = "0"
)

// EthConfig contains the configuration relative to the accounts, EVM, trie/db,
// and service API
type EthConfig struct {
	// Megabytes of memory allocated to internal caching (min 16MB / database forced)
	Cache int `mapstructure:"cache"`

	// Minimum gasprice for transactions submitted via this node
	MinGasPrice string `mapstructure:"min-gas-price"`
}

// DefaultEthConfig return the default configuration for Eth services
func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		Cache:       defaultCache,
		MinGasPrice: defaultMinGasPrice,
	}
}
