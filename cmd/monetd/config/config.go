package config

import (
	_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/monetd/src/poa/common"
)

// Configuration object from EVM-Lite
var (
	Config = monetConfig("")
)

// default config for monetd
func monetConfig(dataDir string) *_config.Config {
	config := _config.DefaultConfig()

	config.Babble.EnableFastSync = false
	config.Babble.Store = true

	if dataDir == "" { // fall through to default settings
		file, _ := common.DefaultMonetConfigDir()
		config.SetDataDir(file)
	} else {
		config.SetDataDir(dataDir)
	}
	return config
}
