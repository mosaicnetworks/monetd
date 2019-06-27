package network

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	keys "github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/viper"
)

const tomlDir = ".monetcli"
const (
	tomlName = "network"
)

// We declare our own instance of viper to avoid any possibility of a clash
var (
	networkViper *viper.Viper
)

func setUpConfigFile() {
	networkViper = viper.New()
	networkViper.SetConfigName(tomlName) // name of config file (without extension)
	networkViper.SetConfigType("toml")
	defaultConfig()
}

// Write configure to file
func writeConfig() {

	message("Writing toml file")
	err := networkViper.WriteConfig()
	if err != nil {
		message("writeConfig error: ", err)
	}
}

func safeWriteConfig() {

	message("Writing toml file")
	err := networkViper.SafeWriteConfig()
	if err != nil {
		message("safeWriteConfig error: ", err)
		message(networkViper.AllSettings())
		networkViper.Debug()
	}
}

func createEmptyFile(f string) {
	emptyFile, err := os.Create(f)

	if err != nil {
		message("Create empty file: ", f, err)
		return
	}
	emptyFile.Close()
}

func loadConfig() error {

	message("Starting to load configuration")
	setUpConfigFile()
	networkViper.AddConfigPath(configDir)
	message("Added viper config path: ", configDir)
	err := networkViper.ReadInConfig() // Find and read the config file

	if err != nil {
		message("loadConfig: ", err)
		return err
	}
	message("Loaded Config")
	return nil
}

func isEmptyDir(dir string) (bool, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func GenerateKeyPair(configDir string, moniker string, ip string, isValidator bool) error {
	message("Generating key pair for: ", moniker)

	targetDir := filepath.Join(configDir, moniker)

	message("Generate to :", targetDir)

	if common.CheckIfExists(targetDir) {
		message("Key Pair for " + moniker + " already exists. Aborting.")
		return errors.New("key pair for " + moniker + " already exists")
	}

	targetFile := filepath.Join(targetDir, keys.DefaultKeyfile)

	/*   // Not required, handled by GenerateKeyPair
	message("Creating dir: ", targetDir)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}
	*/

	key, err := keys.GenerateKeyPair(targetFile)

	if err != nil {
		return err
	}

	pubkey := hex.EncodeToString(
		crypto.FromECDSAPub(&key.PrivateKey.PublicKey))

	return addValidatorParamaterised(moniker, key.Address.Hex(),
		pubkey, ip, isValidator)
	//	return nil
}
