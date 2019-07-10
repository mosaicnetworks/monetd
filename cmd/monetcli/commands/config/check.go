package config

import (
	"path/filepath"

	"github.com/mosaicnetworks/babble/src/babble"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

func checkConfig(cmd *cobra.Command, args []string) error {
	// extract code to the
	return checkConfigParams()
}

func checkConfigParams() error {

	// Here we don't throw errors for most errors as check should find errors and exit cleanly.

	// Check directories we expect exist.
	common.MessageWithType(common.MsgInformation, "Monet Configuration Directory is set to: ", monetConfigDir)
	babbledir := filepath.Join(monetConfigDir, common.BabbleDir)
	ethdir := filepath.Join(monetConfigDir, common.EthDir)
	if !directoryChecks(monetConfigDir, "Monet Configuration Directory") ||
		!directoryChecks(babbledir, "Babble Configuration Directory") ||
		!directoryChecks(ethdir, "EVM-Lite Configuration Directory") {
		common.MessageWithType(common.MsgError, "Stopping.")
		return nil
	}

	toml := filepath.Join(monetConfigDir, common.MonetdTomlName+common.TomlSuffix)

	fileToCheck := []string{
		toml,
		filepath.Join(monetConfigDir, common.DefaultKeyfile),
		filepath.Join(monetConfigDir, common.BabbleDir, common.PeersGenesisJSON),
		filepath.Join(monetConfigDir, common.BabbleDir, common.PeersJSON),
		filepath.Join(monetConfigDir, common.BabbleDir, babble.DefaultKeyfile),
		filepath.Join(monetConfigDir, common.EthDir, common.GenesisJSON),
		filepath.Join(monetConfigDir, common.EthDir, common.PwdFile),
		//		filepath.Join(monetConfigDir, "eth", "missing.txt"),
	}

	for _, f := range fileToCheck {
		common.MessageWithType(common.MsgDebug, "Checking file: ", f)
		if !common.CheckIfExists(f) {
			common.MessageWithType(common.MsgWarning, "Missing File ", f)
		}
	}

	// Next we check monetd.toml

	common.MessageWithType(common.MsgInformation, "Monet Configuration File is set to: ", toml)

	tree, err := common.LoadToml(toml)
	if err != nil {
		common.MessageWithType(common.MsgError, "Cannot parse Monet Configuration File")
	}

	pathToCheck := [][]string{
		{"datadir"},
		{"log"},
		{"babble", "listen"},
		{"babble", "service-listen"},
		{"babble", "heartbeat"},
		{"babble", "timeout"},
		{"babble", "cache-size"},
		{"babble", "sync-limit"},
		{"babble", "max-pool"},
		{"eth", "listen"},
		{"eth", "cache"},
		//		[]string{"eth", "missing"},
	}

	common.MessageWithType(common.MsgInformation, "Checking for missing parameters in Monet Configuration File")

	for _, t := range pathToCheck {
		if !tree.HasPath(t) {
			common.MessageWithType(common.MsgWarning, "Monet Configuration File missing parameter ", t)
		}
	}

	common.MessageWithType(common.MsgInformation, "Checked for missing parameters in Monet Configuration File")
	babbleListen := tree.GetPath([]string{"babble", "listen"}).(string)
	common.MessageWithType(common.MsgInformation, "Babble listening on: ", babbleListen)

	foundPeer, err := common.CheckPeersAddress(filepath.Join(monetConfigDir, common.BabbleDir), babbleListen)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error parsing peers", err)
	}

	if !foundPeer {
		common.MessageWithType(common.MsgWarning, "Babble Gossip Endpoint is not in the Peers File: ", babbleListen)
		common.MessageWithType(common.MsgInformation, "This is not an issue if not one of the genesis peers.")
	}

	// Check .monet files for sanity.
	// Cross reference peers.json and monetd.toml for gossiping end point. They should be an exact match.
	return nil
}

func directoryChecks(dir string, descriptor string) bool {
	if dir == "" {
		common.MessageWithType(common.MsgError, descriptor+" is not set. There should be at least a default.")
		return false
	}

	if !common.CheckIfExists(dir) {
		common.MessageWithType(common.MsgError, descriptor+" does not exist.")
		return false
	}
	isDir, err := common.CheckIsDir(dir)
	if err != nil || !isDir {
		common.MessageWithType(common.MsgError, descriptor+" is not a directory.")
		return false
	}
	return true
}
