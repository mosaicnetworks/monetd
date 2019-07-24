package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/pelletier/go-toml"
)

type urlList struct {
	URL  string
	Dest string
}

//PullConfig pulls a monet configuration from an existing monet node
func PullConfig(configDir, nodeParam, addressParam, existingPeer, passwordFile string) error {
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
	keyFile := filepath.Join(configDir, common.KeyStoreDir, safeLabel+".json")

	// Verify keyfile exists
	if !files.CheckIfExists(keyFile) {
		common.DebugMessage("Error opening keyfile: ", keyFile)
		return errors.New("Cannot open keyfile")
	}

	// Create Directories If they don't exist
	files.CreateMonetConfigFolders(configDir)

	babbleRootURL := "http://" + existingPeer + ":" + common.DefaultBabblePort
	ethRootURL := "http://" + existingPeer + ":" + common.DefaultEVMLitePort

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

	//	files.CopyFileContents(keyFile, filepath.Join(configDir, "keystore", safeLabel+".json"))

	privateKey, err := crypto.GetPrivateKey(keyFile, passwordFile)
	if err != nil {
		return err
	}
	files.WriteToFile(filepath.Join(configDir, common.BabbleDir, common.DefaultPrivateKeyFile), privateKey)

	monettomlfile := filepath.Join(configDir, common.MonetTomlFile)

	/*
		tomlfile := filepath.Join(configDir, common.MonetcliTomlName+common.TomlSuffix)

		var tree *toml.Tree

		if files.CheckIfExists(tomlfile) {
			tree, err = files.LoadToml(tomlfile)
		} else {
			tree, err = toml.Load("")
		}
	*/
	// With the removal of the monetcli structures, we have just populate
	// monet.toml with the defaults.

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

	/*
		err = common.SendKeyToEVMLC(safeLabel, keyFile)
		if err != nil {
			return err
		}
	*/

	return nil
}
