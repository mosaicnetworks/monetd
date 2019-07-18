package config

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

func buildConfig(cmd *cobra.Command, args []string) error {

	// Some debug output confirming parameters
	common.MessageWithType(common.MsgDebug, "Building Config for: ", nodeParam)
	common.MessageWithType(common.MsgDebug, "Using Network Address: ", addressParam)
	common.MessageWithType(common.MsgDebug, "Using Password File: ", passwordFile)

	// Reject empty parameters with a helpful message
	if strings.TrimSpace(nodeParam) == "" {
		return errors.New("--node parameter is not set")
	}
	if strings.TrimSpace(addressParam) == "" {
		return errors.New("--address parameter is not set")
	}
	//	if strings.TrimSpace(passwordFile) == "" {
	//		return errors.New("--passfile parameter is not set")
	//	}

	safeLabel := common.GetNodeSafeLabel(nodeParam)
	keyFile := filepath.Join(networkConfigDir, common.MonetAccountsSubFolder, safeLabel, common.DefaultKeyfile)

	// Verify keyfile exists
	if !common.CheckIfExists(keyFile) {
		common.MessageWithType(common.MsgDebug, "Error opening keyfile: ", keyFile)
		return errors.New("Cannot open keyfile")
	}

	// Create Directories If they don't exist

	dirList := []string{
		monetConfigDir,
		filepath.Join(monetConfigDir, common.BabbleDir),
		filepath.Join(monetConfigDir, common.EthDir),
		filepath.Join(monetConfigDir, common.POADir),
		filepath.Join(monetConfigDir, common.EthDir, "keystore"),
	}

	for _, dir := range dirList {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			common.Message("Error creating directory: ", dir)
			return err
		}
		common.Message("Created directory: ", dir)
	}

	// Generate Peers List

	peerslist := []string{safeLabel} // TODO Add the extended peers list in here

	addr := addressParam + ":" + common.DefaultGossipPort

	peersJSON := common.PeerRecordList{}

	for idx, peer := range peerslist {
		if idx > 0 { // Not the node for this instance so source info from supplied list

		}

		tomlfile := filepath.Join(networkConfigDir, common.MonetAccountsSubFolder, peer, common.NodeFile)
		tree, err := common.LoadToml(tomlfile)
		if err != nil {
			common.MessageWithType(common.MsgError, "Cannot read peer configuration: ", peer)
			return err
		}
		pubkey := tree.GetPath([]string{"node", "pubkey"}).(string)
		moniker := tree.GetPath([]string{"node", "moniker"}).(string)

		peersJSON = append(peersJSON, &common.PeerRecord{NetAddr: addr, PubKeyHex: pubkey, Moniker: moniker})
	}

	// Write Peers.Json
	peersJSONOut, err := json.MarshalIndent(peersJSON, "", "\t")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(monetConfigDir, common.BabbleDir, common.PeersJSON)
	common.WriteToFile(jsonFileName, string(peersJSONOut))
	jsonFileName = filepath.Join(monetConfigDir, common.BabbleDir, common.PeersGenesisJSON)
	common.WriteToFile(jsonFileName, string(peersJSONOut))

	// Create genesis.json
	// Copy keyfile.json to eth/keystore

	common.CopyFileContents(keyFile, filepath.Join(monetConfigDir, common.EthDir, "keystore", safeLabel+".json"))

	// Derive Private Key and Write to Babble Config
	keyjson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyFile, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := common.GetPassphrase(passwordFile)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	// This may go
	privateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	common.WriteToFile(filepath.Join(monetConfigDir, common.BabbleDir, babble.DefaultKeyfile), privateKey)

	//TODO you check .monetcli/network.toml for an updated contract address

	err = common.BuildGenesisJSON(networkConfigDir, monetConfigDir, peersJSON, common.DefaultContractAddress)
	//	monetcliConfigDir string, monetdConfigDir string, peers PeerRecordList, contractAddress string)
	if err != nil {
		return err
	}

	tomlfile := filepath.Join(networkConfigDir, common.MonetcliTomlName+common.TomlSuffix)
	monettomlfile := filepath.Join(monetConfigDir, common.MonetdTomlName+common.TomlSuffix)

	var tree *toml.Tree

	if common.CheckIfExists(tomlfile) {
		tree, err = common.LoadToml(tomlfile)
	} else {
		tree, err = toml.Load("")
	}
	if err != nil {
		return err
	}

	err = common.TransformCliTomlToD(tree, monetConfigDir)
	if err != nil {
		return err
	}

	err = common.SaveToml(tree, monettomlfile)
	if err != nil {
		return err
	}

	err = common.SendKeyToEVMLC(safeLabel, keyFile)
	if err != nil {
		return err
	}

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
