package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
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
	common.DebugMessage("Building Config for  : ", moniker)
	common.DebugMessage("Using Network Address: ", selfAddress)
	common.DebugMessage("Pulling from         : ", otherAddress)
	common.DebugMessage("Using Password File  : ", passwordFile)
	common.DebugMessage("Using Moniker        : ", moniker)

	// Set global moniker config
	configuration.Global.Babble.Moniker = moniker

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
		err := files.DownloadFile(item.URL, item.Dest, files.BackupExisting|files.PromptIfExisting)
		if err != nil {
			common.ErrorMessage(fmt.Sprintf("Error downloading %s", item.URL))
			return err
		}
		common.DebugMessage("Downloaded ", item.Dest)
	}

	// Write TOML file for monetd based on global config object
	if err := DumpConfigTOML(configDir, configuration.MonetTomlFile); err != nil {
		return err
	}

	return nil
}
