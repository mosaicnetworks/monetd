// Package crypto provides common functions for manipulating and generating keys
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	com "github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type outputGenerate struct {
	Address      string
	AddressEIP55 string
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
//
func NewKeyPair(configDir, moniker, passwordFile string) (*keystore.Key, error) {

	safeLabel := com.GetNodeSafeLabel(moniker)

	dirlist := []string{ configDir,
		filepath.Join(configDir, com.KeyStoreDir)
	}

	keyfilepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".json")
	tomlfilepath := filepath.Join(configDir, com.KeyStoreDir, safeLabel+".toml")

	key, err := GenerateKeyPair(keyfilepath, passwordFile, "", false)

	return key, err
}
