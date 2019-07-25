package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"
)

type downloadItem struct {
	URL  string
	Dest string
}

// PullConfig pulls a monet configuration from an existing monet node.
// configDir: monetd datadir
// moniker: moniker of key to use as validator
// seflAddress: network address for monetd
// otherAddress: address of peer to pull configuration from
func PullConfig(configDir, moniker, selfAddress, otherAddress, passwordFile string) error {

	// Helpful debugging output
	common.MessageWithType(common.MsgDebug, "Building Config for  : ", moniker)
	common.MessageWithType(common.MsgDebug, "Using Network Address: ", selfAddress)
	common.MessageWithType(common.MsgDebug, "Pulling from         : ", otherAddress)
	common.MessageWithType(common.MsgDebug, "Using Password File  : ", passwordFile)

	// Retrieve the keyfile corresponding to moniker
	privateKey, err := getKey(configDir, moniker, passwordFile)
	if err != nil {
		return err
	}

	// Create Directories if they don't exist
	files.CreateMonetConfigFolders(configDir)

	// Copy they key to babble directory with appropriate permissions
	if err := dumpPrivKey(configDir, privateKey); err != nil {
		return err
	}

	rootURL := "http://" + otherAddress

	filesList := []*downloadItem{
		&downloadItem{URL: rootURL + "/genesispeers",
			Dest: filepath.Join(configDir, configuration.BabbleDir, configuration.PeersGenesisJSON)},
		&downloadItem{URL: rootURL + "/peers",
			Dest: filepath.Join(configDir, configuration.BabbleDir, configuration.PeersJSON)},
		&downloadItem{URL: rootURL + "/genesis",
			Dest: filepath.Join(configDir, configuration.EthDir, configuration.GenesisJSON)},
	}

	for _, item := range filesList {
		err := files.DownloadFile(item.URL, item.Dest)
		if err != nil {
			common.ErrorMessage(fmt.Sprintf("Error downloading %s", item.URL))
			return err
		}
		common.DebugMessage("Downloaded ", item.Dest)
	}

	// Write TOML file for monetd based on global config object
	if err := dumpConfigTOML(configDir, configuration.MonetTomlFile); err != nil {
		return err
	}

	return nil
}
