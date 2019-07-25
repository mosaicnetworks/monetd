package config

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"
)

type urlList struct {
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

	babbleRootURL := "http://" + otherAddress + ":" + common.DefaultBabblePort
	ethRootURL := "http://" + otherAddress + ":" + common.DefaultEVMLitePort

	filesList := []*urlList{
		&urlList{URL: babbleRootURL + "/genesispeers",
			Dest: filepath.Join(configDir, common.BabbleDir, common.PeersGenesisJSON)},
		&urlList{URL: babbleRootURL + "/peers",
			Dest: filepath.Join(configDir, common.BabbleDir, common.PeersJSON)},
		&urlList{URL: ethRootURL + "/genesis",
			Dest: filepath.Join(configDir, common.EthDir, common.GenesisJSON)},
	}

	for _, filemap := range filesList {
		err := files.DownloadFile(filemap.URL, filemap.Dest)
		if err != nil {
			common.ErrorMessage("Error downloading genesis peers")
			return err
		}
		common.DebugMessage("Downloaded ", filemap.Dest)
	}

	// Write TOML file for monetd based on global config object
	if err := dumpConfigTOML(configDir, common.MonetTomlFile); err != nil {
		return err
	}

	return nil
}
