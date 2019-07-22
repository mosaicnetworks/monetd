// Package crypto provides common functions for manipulating and generating keys
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	com "github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/pelletier/go-toml"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type outputGenerate struct {
	Address      string
	AddressEIP55 string
}

type outputInspect struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

// PublicKeyHexToAddressHex takes a Hex string public key and returns a hex string
// Ethereum style address.
func PublicKeyHexToAddressHex(publicKey string) (string, error) {
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	pubKeyHash := crypto.Keccak256(pubBytes[1:])[12:]

	return ethcommon.BytesToAddress(pubKeyHash).Hex(), nil
}

//GenerateKeyPair generates an Ethereum key pair. keyfilepath is the path to write the new keyfile to.
//passwordFile is a plain text file containing the passphrase to use for the keyfile. privateKeyfile is the
//path to a private key. If specified, this function does not generate a new keyfile, it instead
//generates a keyfile from the private key. outputJSON controls whether the output to stdio is in
//JSON format or not.  The function returns a key object which can be used to retrive public or private
//keys or the address.
func GenerateKeyPair(keyfilepath, passwordFile, privateKeyfile string, outputJSON bool) (*keystore.Key, error) {
	if keyfilepath == "" {
		keyfilepath = com.DefaultKeyfile
	}
	if _, err := os.Stat(keyfilepath); err == nil {
		return nil, fmt.Errorf("Keyfile already exists at %s", keyfilepath)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error checking if keyfile exists: %v", err)
	}

	var privateKey *ecdsa.PrivateKey
	var err error
	if file := privateKeyfile; file != "" {
		// Load private key from file.
		privateKey, err = crypto.LoadECDSA(file)
		if err != nil {
			return nil, fmt.Errorf("Can't load private key: %v", err)
		}
	} else {
		// If not loaded, generate random.
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("Failed to generate random private key: %v", err)
		}
	}

	//Create the keyfile object with a random UUID
	//It would be preferable to create the key manually, rather by calling this
	//function, but we cannot use pborman/uuid directly because it is vendored
	//in go-ethereum. That package defines the type of keystore.Key.Id.
	key := keystore.NewKeyForDirectICAP(rand.Reader)
	key.Address = crypto.PubkeyToAddress(privateKey.PublicKey)
	key.PrivateKey = privateKey

	// Encrypt key with passphrase.
	passphrase, err := GetPassphrase(passwordFile, true)
	if err != nil {
		return nil, err
	}

	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting key: %v", err)
	}

	// Store the file to disk.
	if err := os.MkdirAll(filepath.Dir(keyfilepath), 0700); err != nil {
		return nil, fmt.Errorf("Could not create directory %s: %v", filepath.Dir(keyfilepath), err)
	}
	if err := ioutil.WriteFile(keyfilepath, keyjson, 0600); err != nil {
		return nil, fmt.Errorf("Failed to write keyfile to %s: %v", keyfilepath, err)
	}

	// Output some information.
	out := outputGenerate{
		Address: key.Address.Hex(),
	}

	if outputJSON {
		com.MustPrintJSON(out)
	} else {
		fmt.Println("Address:", out.Address)
	}

	return key, nil
}

//NewKeyPair is a wrapper to GenerateKeyPair. It does not support setting a private key.
//Additionally it does not support outputting to JSON format - if required, that can be
//achieved calling GenerateKeyPair directly.
func NewKeyPair(configDir, moniker, passwordFile string) (*keystore.Key, error) {

	if strings.TrimSpace(moniker) == "" {
		return nil, errors.New("moniker is not set")
	}

	safeLabel := com.GetNodeSafeLabel(moniker)

	dirlist := []string{configDir,
		filepath.Join(configDir, com.KeyStoreDir),
	}

	err := files.CreateDirsIfNotExists(dirlist)

	if err != nil {
		com.ErrorMessage("cannot create keystore directories")
		return nil, err
	}

	keyfilepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".json")
	tomlfilepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".toml")

	if files.CheckIfExists(keyfilepath) {
		return nil, errors.New("key for node " + safeLabel + " already exists")
	}

	key, err := GenerateKeyPair(keyfilepath, passwordFile, "", false)

	if err != nil {
		return key, err
	}

	err = WriteTomlForKey(moniker, safeLabel, tomlfilepath, key)
	// Skip err check as we are returning the error anyway.
	return key, err
}

// WriteTomlForKey takes a key and writes the .toml file for it with Address,
// ID, pubkey etc. It allows basic information about a key file to be assertained
// without having to decrypt the keyfile each time.
func WriteTomlForKey(monikerParam, safeLabel, tomlfilepath string, key *keystore.Key) error {

	com.DebugMessage("Generated Address      : ", key.Address.Hex())
	com.DebugMessage("Generated PubKey       : ", hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)))
	com.DebugMessage("Generated ID           : ", key.Id)

	//TODO Remove this line
	com.DebugMessage("Generated Private Key  : ", hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))

	// write peers config
	tree, err := toml.Load("")
	if err != nil {
		return err
	}
	tree.SetPath([]string{"node", "moniker"}, monikerParam)
	tree.SetPath([]string{"node", "label"}, safeLabel)
	tree.SetPath([]string{"node", "address"}, key.Address.Hex())
	tree.SetPath([]string{"node", "pubkey"}, hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)))

	//TODO Remove this line
	tree.SetPath([]string{"node", "privatekey"}, hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))

	// pass any errors back up the function tree
	return common.SaveToml(tree, tomlfilepath)
}

func InspectKeyMoniker(configDir string, moniker string, PasswordFile string, showPrivate bool, outputJSON bool) error {
	safeLabel := com.GetNodeSafeLabel(moniker)
	filepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".json")

	if !files.CheckIfExists(filepath) {
		return errors.New("cannot find keyfile for that moniker")
	}

	return InspectKey(filepath, PasswordFile, showPrivate, outputJSON)
}

// InspectKey inspects an encrypted keyfile
func InspectKey(keyfilepath string, PasswordFile string, showPrivate bool, outputJSON bool) error {

	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := common.GetPassphrase(PasswordFile, false)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	// Output all relevant information we can retrieve.
	out := outputInspect{
		Address: key.Address.Hex(),
		PublicKey: hex.EncodeToString(
			crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
	}
	if showPrivate {
		out.PrivateKey = hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	}

	if outputJSON {
		com.MustPrintJSON(out)
	} else {
		fmt.Println("Address:       ", out.Address)
		fmt.Println("Public key:    ", out.PublicKey)
		if showPrivate {
			fmt.Println("Private key:   ", out.PrivateKey)
		}
	}

	return nil
}

// UpdateKeysMoniker wraps UpdateKeys adding moniker support
func UpdateKeysMoniker(configDir string, moniker string, PasswordFile string, newPasswordFile string) error {
	safeLabel := com.GetNodeSafeLabel(moniker)
	filepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".json")

	if !files.CheckIfExists(filepath) {
		return errors.New("cannot find keyfile for that moniker")
	}

	return UpdateKeys(filepath, PasswordFile, newPasswordFile)
}

// UpdateKeys changes the passphrase on an encrypted keyfile
func UpdateKeys(keyfilepath string, PasswordFile string, newPasswordFile string) error {
	//	keyfilepath := args[0]

	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := common.GetPassphrase(PasswordFile, false)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	// Get a new passphrase.
	fmt.Println("Please provide a new passphrase")
	var newPhrase string
	if newPasswordFile != "" {
		content, err := ioutil.ReadFile(newPasswordFile)
		if err != nil {
			return fmt.Errorf("Failed to read new passphrase file '%s': %v", newPasswordFile, err)
		}
		newPhrase = strings.TrimRight(string(content), "\r\n")
	} else {
		newPhrase, err = common.PromptPassphrase(true)
		if err != nil {
			return err
		}
	}

	// Encrypt the key with the new passphrase.
	newJSON, err := keystore.EncryptKey(key, newPhrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("Error encrypting with new passphrase: %v", err)
	}

	// Then write the new keyfile in place of the old one.
	if err := ioutil.WriteFile(keyfilepath, newJSON, 600); err != nil {
		return fmt.Errorf("Error writing new keyfile to disk: %v", err)
	}

	// Don't print anything.  Just return successfully,
	// producing a positive exit code.
	return nil
}
