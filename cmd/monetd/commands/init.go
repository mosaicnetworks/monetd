package commands

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"
	bkeys "github.com/mosaicnetworks/babble/src/crypto/keys"
	bpeers "github.com/mosaicnetworks/babble/src/peers"
	"github.com/spf13/cobra"
)

var (
	defaultKeyfile = "keyfile.json"
	privateKeyfile string
	passphrase     string
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive configuration wizard",
	RunE:  initialise,
}

/*
1) Private Key
2) Create Babble peers.json
3) Create genesis.json
4) Create monetd.toml
*/
func initialise(cmd *cobra.Command, args []string) error {

	fmt.Printf("\n Step 1) Private Key \n\n")

	if err := initPrivateKey(); err != nil {
		return err
	}

	fmt.Printf("\n Step 2) Create Babble peers.json \n\n")

	if err := initPeersJSON(); err != nil {
		return err
	}

	fmt.Printf("\n Step 3) Create genesis.json \n\n")

	if err := initGenesisJSON(); err != nil {
		return err
	}

	fmt.Printf("\n Step 4) Create default monetd.toml \n\n")

	if err := initMonetTOML(); err != nil {
		return err
	}

	fmt.Printf("\nConfiguration initialised. Use 'monetd run' to start the node.\n")

	return nil
}

// if [datadir]/keyfile.json AND [datadir]/babble/priv_key DO NOT exist:
// generate a new encrypted key in [datadir]/keyfile.json and extract the raw
// private key into [datadir]/babble/priv_key
func initPrivateKey() error {

	// Check that [datadir]/keyfile.json AND [datadir]/babble/priv_key DO NOT
	// exist.
	jsonKeyFilepath := fmt.Sprintf("%s/%s", config.DataDir, defaultKeyfile)
	rawKeyFilepath := fmt.Sprintf("%s/priv_key", config.Babble.DataDir)

	// return with no error ff a key already exists and the user does not wish
	// to generate a new one

	generateNew, err := shouldGenerateNew(config.DataDir, defaultKeyfile)
	if err != nil {
		return err
	}

	if !generateNew {
		return nil
	}

	// generate new private key.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate random private key: %v", err)
	}

	// Create the keyfile object with a random UUID
	// It would be preferable to create the key manually, rather than calling
	// this function, but we cannot use pborman/uuid directly because it is
	// vendored in go-ethereum. That package defines the type of keystore.Key.Id
	key := keystore.NewKeyForDirectICAP(rand.Reader)
	key.Address = crypto.PubkeyToAddress(privateKey.PublicKey)
	key.PrivateKey = privateKey

	// Encrypt key with passphrase.
	fmt.Printf("The private key will be encrypted and copied to %s\n", jsonKeyFilepath)

	passphrase, err := promptPassphrase(true)
	if err != nil {
		return err
	}

	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("Error encrypting key: %v", err)
	}

	// Store the keyfile to disk.
	if err := os.MkdirAll(filepath.Dir(jsonKeyFilepath), 0700); err != nil {
		return fmt.Errorf("Could not create directory %s: %v", filepath.Dir(jsonKeyFilepath), err)
	}
	if err := ioutil.WriteFile(jsonKeyFilepath, keyjson, 0600); err != nil {
		return fmt.Errorf("Failed to write keyfile to %s: %v", jsonKeyFilepath, err)
	}

	fmt.Printf("Succesfully created %s\n", jsonKeyFilepath)

	// extract priv_key to /babble
	simpleKeyfile := bkeys.NewSimpleKeyfile(rawKeyFilepath)
	if err := simpleKeyfile.WriteKey(privateKey); err != nil {
		return fmt.Errorf("Error saving private key: %s", err)
	}

	fmt.Printf("Extracted raw private key to %s\n", rawKeyFilepath)

	fmt.Println(`
	**********************************************************************
	* Please take all the necessary precautions to secure these files    * 
	* and remember the password, as it will be impossible to recover the *
	* key without them.                                                  *
	**********************************************************************
	`)

	return nil
}

// Create peers.json with a single peer corresponding to the key created in step
// 1, and copy it to [datadir]/babble/peers.json
func initPeersJSON() error {

	generateNew, err := shouldGenerateNew(config.Babble.DataDir, "peers.json")
	if err != nil {
		return err
	}

	if !generateNew {
		return nil
	}

	privateKey, err := getPrivateKey()
	if err != nil {
		return err
	}

	peer := &bpeers.Peer{
		NetAddr:   fmt.Sprintf(config.Babble.BindAddr),
		PubKeyHex: bkeys.PublicKeyHex(&privateKey.PublicKey),
		Moniker:   defaultMoniker(),
	}

	peers := []*bpeers.Peer{peer}

	store := bpeers.NewJSONPeerSet(config.Babble.DataDir, true)

	if err := store.Write(peers); err != nil {
		return fmt.Errorf("Error writing Babble peers.json file: %v", err)
	}

	fmt.Printf("Created peers.json file in %v\n", config.Babble.DataDir)

	return nil
}

// Create [datadir]/eth/genesis.json and use the key generated in step 1 to
// prefund the corresponding account and set it as a validator in the POA
// smart-contract.
func initGenesisJSON() error {

	generateNew, err := shouldGenerateNew(filepath.Dir(config.Eth.Genesis), "genesis.json")
	if err != nil {
		return err
	}

	if !generateNew {
		return nil
	}

	privateKey, err := getPrivateKey()
	if err != nil {
		return err
	}

	validatorAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	validatorMoniker := defaultMoniker()

	abi, code, err := compilePOA(validatorAddress, validatorMoniker)
	if err != nil {
		return err
	}

	tmpl, err := template.New("genesis").Parse(genesis)
	if err != nil {
		return fmt.Errorf("Error parsing genesis.json template: %v", err)
	}

	type genesisInfo struct {
		ValidatorAddress, ValidatorMoniker string
		ABI, Code                          string
	}

	genesisArgs := genesisInfo{
		ValidatorAddress: validatorAddress,
		ValidatorMoniker: validatorMoniker,
		ABI:              abi,
		Code:             code,
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, genesisArgs)
	if err != nil {
		return fmt.Errorf("Error executing genesis.json template: %v", err)
	}

	genesisJSON := buf.String()

	if err := os.MkdirAll(filepath.Dir(config.Eth.Genesis), 0700); err != nil {
		return fmt.Errorf("Could not create directory %s: %v", filepath.Dir(config.Eth.Genesis), err)
	}
	if err := ioutil.WriteFile(config.Eth.Genesis, []byte(genesisJSON), 0600); err != nil {
		return fmt.Errorf("Failed to write genesis.json to %s: %v", config.Eth.Genesis, err)
	}

	fmt.Printf("Created genesis.json file in %v\n", config.Eth.Genesis)

	return nil
}

// Generate default configuration file in [datadir]/monetd.toml
func initMonetTOML() error {

	generateNew, err := shouldGenerateNew(config.DataDir, "monetd.toml")
	if err != nil {
		return err
	}

	if !generateNew {
		return nil
	}

	configTmpl, err := template.New("monetd.toml").Parse(configTOML)
	if err != nil {
		return fmt.Errorf("Error parsing monetd.toml template: %v", err)
	}

	var buf bytes.Buffer
	err = configTmpl.Execute(&buf, config)
	if err != nil {
		return fmt.Errorf("Error executing monetd.toml template: %v", err)
	}
	configString := buf.String()

	configFile := fmt.Sprintf("%s/monetd.toml", config.DataDir)

	if err := ioutil.WriteFile(configFile, []byte(configString), 0600); err != nil {
		return fmt.Errorf("Failed to write monetd.toml to %s: %v", configFile, err)
	}

	fmt.Printf("Created monetd.toml file in %v\n", configFile)

	return nil
}

func shouldGenerateNew(path string, name string) (bool, error) {

	filePath := fmt.Sprintf("%s/%s", path, name)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("A %s file already exists in %s\n", name, path)

		generateNew, err := console.Stdin.PromptConfirm(
			fmt.Sprintf("Generate new %s? (this will override the existing one)", name))
		if err != nil {
			return false, fmt.Errorf("Error while asking whether to generate new %s: %v", name, err)
		}

		if !generateNew {
			fmt.Printf("Keeping existing %s\n", name)
			return false, nil
		}
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("Error checking if %s exists: %v", name, err)
	}

	return true, nil
}

// Read POA solidity smart-contract template, insert relevant code to preset the
// validator in the whitelist, and compile. Output the ABI and RuntimeCode.
func compilePOA(validatorAddress, validatorMoniker string) (abi string, code string, err error) {
	poaTmpl, err := template.New("poa").Parse(poa)
	if err != nil {
		return "", "", fmt.Errorf("Error creating POA smart-contract template: %v", err)
	}

	type validatorInfo struct {
		ValidatorAddress, ValidatorMoniker string
	}

	poaTemplateArgs := validatorInfo{
		ValidatorAddress: validatorAddress,
		ValidatorMoniker: validatorMoniker,
	}

	var buf bytes.Buffer
	err = poaTmpl.Execute(&buf, poaTemplateArgs)
	if err != nil {
		return "", "", fmt.Errorf("Error executing POA smart-contract template: %v", err)
	}

	poaSolidity := buf.String()

	err = ioutil.WriteFile(fmt.Sprintf("%s/poa.sol", config.DataDir),
		[]byte(poaSolidity),
		0600)
	if err != nil {
		return "", "", fmt.Errorf("Failed to write poa.sol: %v", err)
	}

	contractInfo, err := compiler.CompileSolidityString("solc", poaSolidity)

	poaContract, ok := contractInfo["<stdin>:POA_Genesis"]
	if !ok {
		return "", "", fmt.Errorf("POA_Genesis not found in compiler output")
	}

	jsonabi, err := json.Marshal(poaContract.Info.AbiDefinition)
	if err != nil {
		return "", "", fmt.Errorf("Error marshalling contract ABI: %v", err)
	}

	stringabi := string(jsonabi)

	abi = strings.Replace(stringabi, `"`, `\"`, -1)
	code = poaContract.RuntimeCode[2:]

	return abi, code, nil
}

// Read unencrypted private key from [datadir]/babble/priv_key
func getPrivateKey() (*ecdsa.PrivateKey, error) {
	rawKeyFilepath := fmt.Sprintf("%s/priv_key", config.Babble.DataDir)

	simpleKeyfile := bkeys.NewSimpleKeyfile(rawKeyFilepath)

	privKey, err := simpleKeyfile.ReadKey()
	if err != nil {
		return nil, fmt.Errorf("Error reading private key from %s: %s", rawKeyFilepath, err)
	}

	return privKey, nil
}

func defaultMoniker() string {
	moniker, _ := os.Hostname()
	return moniker
}
