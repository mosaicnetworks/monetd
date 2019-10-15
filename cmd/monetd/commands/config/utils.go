package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
)

// CreateMonetConfigFolders creates the standard directory layout for a monet
// configuration folder
func CreateMonetConfigFolders(configDir string) error {
	return files.CreateDirsIfNotExists([]string{
		configDir,
		filepath.Join(configDir, configuration.BabbleDir),
		filepath.Join(configDir, configuration.EthDir),
		filepath.Join(configDir, configuration.EthDir, configuration.POADir),
	})
}

// ShowIPWarnings outputs warnings if IP addresses are local and propably not
// reachable from the outside.
func ShowIPWarnings() {
	api := configuration.Global.APIAddr
	listen := configuration.Global.Babble.BindAddr
	advertise := configuration.Global.Babble.AdvertiseAddr

	if common.CheckIP(api, true) {
		common.MessageWithType(common.MsgWarning, fmt.Sprintf("Monetd service API address in monetd.toml may be internal: %s", api))
	}

	if advertise != "" && common.CheckIP(advertise, false) {
		common.MessageWithType(common.MsgWarning, fmt.Sprintf("babble.advertise address in monetd.toml may be internal: %s \n", listen))
	} else if common.CheckIP(listen, false) {
		common.MessageWithType(
			common.MsgWarning,
			fmt.Sprintf("babble.listen address in monetd.toml may be internal: %s. Consider setting an advertise address.", listen),
		)
	}
}
