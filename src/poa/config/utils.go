package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/mosaicnetworks/monetd/src/poa/types"
)

// getKer looks in the keystore for a keyfile corresponding to the provided
// moniker. If it exists, it decrypts it and returns the private key. Otherwise,
// it returns an error
func getKey(configDir, moniker, passwordFile string) (*ecdsa.PrivateKey, error) {
	// Find keyfile based on moniker
	safeLabel := common.GetNodeSafeLabel(moniker)

	keyfile := filepath.Join(configDir, configuration.KeyStoreDir, safeLabel+".json")

	privateKey, err := crypto.GetPrivateKey(keyfile, passwordFile)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// dumpConfigTOML takes the global Config object defined in the configuration
// package, encodes it into a TOML string, and writes it to a file.
func dumpConfigTOML(configDir, fileName string) error {
	tomlString, err := configuration.GlobalTOML()
	if err != nil {
		return err
	}

	tomlPath := filepath.Join(configDir, fileName)

	if err := ioutil.WriteFile(tomlPath, []byte(tomlString), 0644); err != nil {
		return fmt.Errorf("Failed to write %s: %v", tomlPath, err)
	}

	return nil
}

// dumpPrivKey converts an ecdsa private key into a hex string and writes it to
// a file with UNIX permissions 600.
func dumpPrivKey(configDir string, privKey *ecdsa.PrivateKey) error {

	keyString := hex.EncodeToString(eth_crypto.FromECDSA(privKey))

	// The private key is writte with 600 permissions because Babble would
	// complain otherwise
	err := ioutil.WriteFile(
		filepath.Join(configDir,
			configuration.BabbleDir,
			configuration.DefaultPrivateKeyFile,
		),
		[]byte(keyString), 0600)
	if err != nil {
		return err
	}

	return nil
}

// dumpPeers takes PeerRecordList and dumps it into peers.json and
// peers.genesis.json in the babble directory
func dumpPeers(configDir string, peers types.PeerRecordList) error {
	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}

	// peers.json
	jsonFileName := filepath.Join(configDir, configuration.BabbleDir, configuration.PeersJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut))

	// peers.genesis.json
	jsonFileName = filepath.Join(configDir, configuration.BabbleDir, configuration.PeersGenesisJSON)
	files.WriteToFile(jsonFileName, string(peersJSONOut))

	return nil
}

// createSoloPeers creates PeerRecordList with a single node
func createSoloPeerRecordList(moniker, selfAddress, pubKey string) (types.PeerRecordList, error) {
	addr := selfAddress + ":" + configuration.DefaultGossipPort

	peers := types.PeerRecordList{
		&types.PeerRecord{NetAddr: addr, PubKeyHex: pubKey, Moniker: moniker},
	}

	return peers, nil
}
