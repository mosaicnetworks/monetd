package config

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/mosaicnetworks/monetd/src/poa/types"
)

//BuildConfig provides the functionality for monet config build.
func BuildConfig(configDir, nodeParam, addressParam, passwordFile string) error {

	// Some debug output confirming parameters
	common.DebugMessage("Building Config for: ", nodeParam)
	common.DebugMessage("Using Network Address: ", addressParam)
	common.DebugMessage("Using Password File: ", passwordFile)

	// Reject empty parameters with a helpful message
	if strings.TrimSpace(nodeParam) == "" {
		return errors.New("--node parameter is not set")
	}
	if strings.TrimSpace(addressParam) == "" {
		return errors.New("--address parameter is not set")
	}

	safeLabel := common.GetNodeSafeLabel(nodeParam)
	keyFile := filepath.Join(configDir, common.KeyStoreDir, safeLabel+".json")

	// Verify keyfile exists
	if !files.CheckIfExists(keyFile) {
		common.DebugMessage("Error opening keyfile: ", keyFile)
		return errors.New("Cannot open keyfile")
	}

	// Create Directories If they don't exist
	files.CreateMonetConfigFolders(configDir)

	// Generate Peers List

	peerslist := []string{safeLabel} // TODO Add the extended peers list in here

	addr := addressParam + ":" + common.DefaultGossipPort

	peersJSON := types.PeerRecordList{}

	for idx, peer := range peerslist {
		if idx > 0 { // Not the node for this instance so source info from supplied list

		}

		tomlfile := filepath.Join(configDir, common.KeyStoreDir, peer+".toml")
		tree, err := files.LoadToml(tomlfile)
		if err != nil {
			common.MessageWithType(common.MsgError, "Cannot read peer configuration: ", peer)
			return err
		}
		pubkey := tree.GetPath([]string{"node", "pubkey"}).(string)
		moniker := tree.GetPath([]string{"node", "moniker"}).(string)

		peersJSON = append(peersJSON, &types.PeerRecord{NetAddr: addr, PubKeyHex: pubkey, Moniker: moniker})
	}

	// Write Peers.Json
	peersJSONOut, err := json.MarshalIndent(peersJSON, "", "\t")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(configDir, common.BabbleDir, common.PeersJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut))
	jsonFileName = filepath.Join(configDir, common.BabbleDir, common.PeersGenesisJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut))

	// Copy keyfile.json to eth/keystore  - redundant as we now have a common keystore
	//	files.CopyFileContents(keyFile, filepath.Join(monetConfigDir, common.EthDir, "keystore", safeLabel+".json"))

	privateKey, err := crypto.GetPrivateKey(keyFile, passwordFile)
	if err != nil {
		return err
	}
	files.WriteToFile(filepath.Join(configDir, common.BabbleDir, common.DefaultPrivateKeyFile), privateKey)

	//  pwd.txt is deprecated
	//	files.WriteToFile(filepath.Join(monetConfigDir, common.EthDir, common.PwdFile), passphrase)

	err = BuildGenesisJSON(configDir, peersJSON, common.DefaultContractAddress)
	//	monetcliConfigDir string, monetdConfigDir string, peers PeerRecordList, contractAddress string)
	if err != nil {
		return err
	}

	monettomlfile := filepath.Join(configDir, common.MonetTomlFile)
	tree, err := toml.Load("")

	if err != nil {
		return err
	}

	err = TransformCliTomlToD(tree, configDir)
	if err != nil {
		return err
	}

	err = files.SaveToml(tree, monettomlfile)
	if err != nil {
		return err
	}

	//TODO - write to wallet.toml to ensure gas is set appropriately highly
	// Set the default from address

	return nil
}

/*
//BuildGenesisJSON ...
func BuildGenesisJSON(peers common.PeerRecordList, networkConfigDir string) error {

	version, err := common.GetSolidityCompilerVersion()
	if err != nil {
		return err
	}

	// Attempts to load contract from standard location, falls back to
	// github
	filename := filepath.Join(networkConfigDir, common.TemplateContract)
	soliditySource, err := common.GetSoliditySource(filename)

	if err != nil || strings.TrimSpace(soliditySource) == "" {
		return errors.New("no valid solidity contract source found")
	}

	// We now have the solidity template in a string. We need to apply
	// the peers to create the initial whitelist.

	finalSoliditySource, err := common.ApplyInitialWhitelistToSoliditySource(soliditySource, peers)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error building genesis contract:", err)
		return err
	}

	// We write out the final source to file

	err = common.WriteToFile(filepath.Join(networkConfigDir, common.GenesisContract), finalSoliditySource)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error writing genesis contract:", err)
		return err
	}

	// And we compile the final source
	contractInfo, err := common.CompileSolidityContract(finalSoliditySource)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error compiling genesis contract:", err)
		return err
	}

	return nil
}
*/

//TransformCliTomlToD process monetcli config to become a monetd config - largely by
//removing unused keys which are used for peers and genesis files.
func TransformCliTomlToD(tree *toml.Tree, configDir string) error {

	delKeys := []string{"poa.bytecode", "poa.abi", "validators", "config.datadir"}

	setKeys := GetMonetDefaultConfigKeys(configDir)

	// First we delete the extraneous keys
	for _, key := range delKeys {

		if tree.Has(key) {
			err := tree.Delete(key)
			if err != nil {
				common.ErrorMessage("Error deleting "+key+": ", err)
				return err
			}
		} else {
			common.DebugMessage("Key " + key + " is not set. Continuing.")
		}
	}

	// Then we set default values and overrides
	common.DebugMessage("Setting default values")
	for _, keys := range setKeys {
		if keys.Override || !tree.Has(keys.Key) {
			common.DebugMessage("Setting "+keys.Key+" to: ", keys.Value)
			tree.Set(keys.Key, keys.Value)
		}
	}

	return nil
}
