package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type urlList struct {
	URL  string
	Dest string
}

func pullConfig(cmd *cobra.Command, args []string) error {
	var err error
	// Helpful debugging output
	common.MessageWithType(common.MsgDebug, "Building Config for  : ", nodeParam)
	common.MessageWithType(common.MsgDebug, "Using Network Address: ", addressParam)
	common.MessageWithType(common.MsgDebug, "Existing Peer        : ", existingPeer)
	common.MessageWithType(common.MsgDebug, "Using Password File  : ", passwordFile)

	// Reject empty parameters with a helpful message
	if strings.TrimSpace(nodeParam) == "" {
		return errors.New("--node parameter is not set")
	}
	if strings.TrimSpace(addressParam) == "" {
		return errors.New("--address parameter is not set")
	}
	if strings.TrimSpace(existingPeer) == "" {
		return errors.New("--peer parameter is not set")
	}

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

	filesList := []*urlList{
		&urlList{URL: "http://" + existingPeer + ":" + common.DefaultBabblePort + "/genesispeers",
			Dest: filepath.Join(monetConfigDir, common.BabbleDir, common.PeersGenesisJSON)},
		&urlList{URL: "http://" + existingPeer + ":" + common.DefaultBabblePort + "/peers",
			Dest: filepath.Join(monetConfigDir, common.BabbleDir, common.PeersJSON)},
		&urlList{URL: "http://" + existingPeer + ":" + common.DefaultEVMLitePort + "/genesis",
			Dest: filepath.Join(monetConfigDir, common.EthDir, common.GenesisJSON)},
	}

	for _, filemap := range filesList {

		err := common.DownloadFile(filemap.URL, filemap.Dest)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error downloading genesis peers")
			return err
		}
		common.MessageWithType(common.MsgDebug, "Downloaded ", filemap.Dest)
	}

	common.CopyFileContents(keyFile, filepath.Join(monetConfigDir, common.EthDir, "keystore", safeLabel+".json"))
	// Decrypt just to confirm this address
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

	common.MessageWithType(common.MsgInformation, key.Address.Hex())

	// This may go
	privateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	common.WriteToFile(filepath.Join(monetConfigDir, common.BabbleDir, babble.DefaultKeyfile), privateKey)
	common.WriteToFile(filepath.Join(monetConfigDir, common.EthDir, common.PwdFile), passphrase)

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
