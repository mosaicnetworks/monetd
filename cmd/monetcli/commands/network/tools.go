package network

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
	bkeys "github.com/mosaicnetworks/babble/src/crypto/keys"
	keys "github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"
	"github.com/mosaicnetworks/monetd/src/common"
	com "github.com/mosaicnetworks/monetd/src/common"
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
)

const tomlDir = ".monetcli"
const (
	tomlName = "network"
)

func createEmptyFile(f string) {
	emptyFile, err := os.Create(f)

	if err != nil {
		message("Create empty file: ", f, err)
		return
	}
	emptyFile.Close()
}

func isEmptyDir(dir string) (bool, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

//GenerateKeyPair wraps monetcli keys functionality, adding more interactivity
func GenerateKeyPair(configDir string, moniker string, ip string, isValidator bool, passwordFile string, privateKeyFile string) error {
	message("Generating key pair for: ", moniker)

	//Enhancement - pass safeLabel into this function. Validation should have happened further up the stack
	safeLabel := common.GetNodeSafeLabel(moniker)

	tree, err := com.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	if tree.HasPath([]string{"validators", safeLabel}) {
		// Duplicate Node
		return errors.New("cannot generate a node with a duplicate moniker")
	}

	targetDir := filepath.Join(configDir, moniker)
	message("Generate to :", targetDir)
	targetFile := filepath.Join(configDir, moniker, common.DefaultKeyfile)

	if common.CheckIfExists(targetFile) {
		message("Key Pair for " + moniker + " already exists. Aborting.")
		return errors.New("key pair for " + moniker + " already exists")
	}

	key, err := keys.GenerateKeyPair(targetFile, passwordFile)
	if err != nil {
		return err
	}

	pubkey := hex.EncodeToString(
		crypto.FromECDSAPub(&key.PrivateKey.PublicKey))

	privateKey := key.PrivateKey

	if privateKeyFile != "" {
		simpleKeyfile := bkeys.NewSimpleKeyfile(privateKeyFile)
		if err := simpleKeyfile.WriteKey(privateKey); err != nil {
			return err
		}
	}

	return AddValidatorParamaterised(configDir, moniker, safeLabel, key.Address.Hex(),
		pubkey, ip, isValidator)

}

//GetPeersLabelsListFromToml processes a monetcli toml file and returns
//a string slice
func GetPeersLabelsListFromToml(configDir string) ([]string, error) {
	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return nil, err
	}

	return GetPeersLabelsList(tree), nil
}

//GetPeersLabelsList takes a tree and extracts a peer list from it
func GetPeersLabelsList(tree *toml.Tree) []string {
	var rtn []string

	validators, err := query.CompileAndExecute("$.validators", tree)
	if err != nil {
		common.Message("Error Getting Peers Labels")
		return rtn
	}

	for _, value := range validators.Values() {

		if reflect.TypeOf(value).String() == "*toml.Tree" {
			v := value.(*toml.Tree)

			keys := v.Keys()
			return keys
		}
	}

	return rtn
}
