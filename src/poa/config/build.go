package config

import (
	"encoding/hex"

	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"
)

// BuildConfig provides the functionality for monetd config build.
// configDir: monetd datadir
// moniker: moniker of key to use as validator
// selfAddress: network address for monetd
func BuildConfig(configDir, moniker, selfAddress, passwordFile string) error {

	// Some debug output confirming parameters
	common.DebugMessage("Building Config for: ", moniker)
	common.DebugMessage("Using Network Address: ", selfAddress)
	common.DebugMessage("Using Password File: ", passwordFile)

	// Retrieve the keyfile corresponding to moniker
	privateKey, err := getKey(configDir, moniker, passwordFile)
	if err != nil {
		return err
	}

	// Create Directories if they don't exist
	files.CreateMonetConfigFolders(configDir)

	// Copy the key to babble directory with appropriate permissions
	if err := dumpPrivKey(configDir, privateKey); err != nil {
		return err
	}

	pubKey := hex.EncodeToString(eth_crypto.FromECDSAPub(&privateKey.PublicKey))

	// Create a peer-set whith a single node
	peers, err := createSoloPeerRecordList(moniker, selfAddress, pubKey)
	if err != nil {
		return err
	}

	// Write peers.json and peers.genesis.json
	if err := dumpPeers(configDir, peers); err != nil {
		return err
	}

	// Create the eth/genesis.json file
	err = BuildGenesisJSON(configDir, peers, common.DefaultContractAddress)
	if err != nil {
		return err
	}

	// Write TOML file for monetd based on global config object
	if err := dumpConfigTOML(configDir, common.MonetTomlFile); err != nil {
		return err
	}

	return nil
}
