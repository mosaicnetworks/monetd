package config

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/mosaicnetworks/monetd/src/configuration"
)

//ImportZip imports a monetd config zip file (src) and writes to the appropriate
//subfolder in dest
func ImportZip(src string, dest string) error {

	var filenames []string
	var keyfile string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		var fpath string

		switch f.Name {
		case configuration.GenesisJSON:
			fpath = filepath.Join(dest, configuration.EthDir, configuration.GenesisJSON)
		case configuration.MonetTomlFile:
			fpath = filepath.Join(dest, configuration.MonetTomlFile)
		case configuration.PeersJSON:
			fpath = filepath.Join(dest, configuration.BabbleDir, configuration.PeersJSON)
		case configuration.PeersGenesisJSON:
			fpath = filepath.Join(dest, configuration.BabbleDir, configuration.PeersGenesisJSON)

		default:
			fpath = filepath.Join(dest, configuration.KeyStoreDir, f.Name)
			if strings.ToLower(filepath.Ext(f.Name)) == ".json" {
				keyfile = f.Name[0 : len(f.Name)-5] // 5 is length of .json
			}
		}

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		common.DebugMessage("Writing file " + fpath)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	// Files have been copied into place

	// Need to get listen parameter from peers.json

	listen, err := getListenForPeer(keyfile, filepath.Join(dest, configuration.BabbleDir, configuration.PeersJSON))
	if err != nil {
		return err
	}

	// Need to edit monetd.toml and set datadir and listen
	err = setLocalParamsInToml(dest, filepath.Join(dest, configuration.MonetTomlFile), listen)
	if err != nil {
		return err
	}
	// Need to generate private key

	err = generateBabblePrivateKey(dest, keyfile)
	if err != nil {
		return err
	}

	return nil
}

// getListenForPeer opens the peers.json file and reads it. If it finds a
// moniker match, it uses the local value. Otherwise it uses the current IP
func getListenForPeer(moniker string, peersfile string) (string, error) {

	jsonFile, err := os.Open(peersfile)
	if err != nil {
		common.ErrorMessage("Error opening " + peersfile)
		return "", err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []interface{}
	json.Unmarshal([]byte(byteValue), &result)

	for _, peer := range result {
		if peer.(map[string]interface{})["Moniker"].(string) == moniker {
			netaddr := peer.(map[string]interface{})["NetAddr"].(string)
			if strings.Index(netaddr, ":") < 0 {
				netaddr = netaddr + ":" + configuration.DefaultGossipPort
			}
			common.DebugMessage("Set listen from peers: " + netaddr)
			return netaddr, nil
		}

	}

	netaddr := common.GetMyIP() + ":" + configuration.DefaultGossipPort
	common.DebugMessage("Set listen from peers: " + netaddr)
	return netaddr, nil
}

func setLocalParamsInToml(datadir string, toml string, listen string) error {
	tree, err := files.LoadToml(toml)
	if err != nil {
		return err
	}
	tree.SetPath([]string{"datadir"}, datadir)
	tree.SetPath([]string{"babble", "listen"}, listen)
	files.SaveToml(tree, toml)
	if err != nil {
		return err
	}

	return nil
}

func generateBabblePrivateKey(datadir string, basename string) error {

	privateKeyfile := filepath.Join(datadir, configuration.BabbleDir, configuration.DefaultPrivateKeyFile)
	jsonfile := filepath.Join(datadir, configuration.KeyStoreDir, basename+".json")
	pwdfile := filepath.Join(datadir, configuration.KeyStoreDir, basename+".txt")

	if !files.CheckIfExists(jsonfile) {
		return errors.New("cannot read key file: " + jsonfile)
	}

	if !files.CheckIfExists(pwdfile) {
		common.DebugMessage("No passphrase file available")
		pwdfile = ""
	}

	keyjson, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", jsonfile, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := crypto.GetPassphrase(pwdfile, false)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	addr := key.Address.Hex()
	privkey := hex.EncodeToString(eth_crypto.FromECDSA(key.PrivateKey))

	files.WriteToFile(privateKeyfile, privkey)

	common.DebugMessage("Written Private Key for " + addr)

	return nil
}
