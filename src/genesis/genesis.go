package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/src/common"
)

// AllocRecord is an object that contains information about a pre-funded acount.
type AllocRecord struct {
	Balance string `json:"balance"`
	Moniker string `json:"moniker"`
}

// Alloc is the section of a genesis file that contains the list of pre-funded
// accounts.
type Alloc map[string]*AllocRecord

// POA is the section of a genesis file that contains information about
// the POA smart-contract.
type POA struct {
	Address string            `json:"address"`
	Abi     string            `json:"abi"`
	Code    string            `json:"code"`
	Storage map[string]string `json:"storage,omitempty"`
}

// GenesisFile is the structure that a Genesis file gets parsed into.
type GenesisFile struct {
	Alloc      *Alloc `json:"alloc"`
	Poa        *POA   `json:"poa"`
	Controller *POA   `json:"controller"`
}

// MinimalPeerRecord is used where only an Address and Moniker are required.
// The standard Peer datatypes us PubKeyHex not address.
type MinimalPeerRecord struct {
	Address string
	Moniker string
}

// GenerateGenesisJSON uses a precompiled POA contract
func GenerateGenesisJSON(outDir, keystore string, peers []*peers.Peer, alloc *Alloc, contractAddress string, controllerAddress string) error {

	var genesis GenesisFile
	var miniPeers []*MinimalPeerRecord

	for _, peer := range peers {
		addr, err := crypto.PublicKeyHexToAddressHex(peer.PubKeyHex)
		if err != nil {
			return err
		}
		miniPeers = append(miniPeers, &MinimalPeerRecord{Address: addr, Moniker: peer.Moniker})
	}

	storageJSON, err := GetStorage(miniPeers)
	if err != nil {
		return err
	}

	genesispoa := POA{Code: StandardPOAContractByteCode,
		Address: controllerAddress,
		Abi:     StandardPOAContractABI,
		Storage: storageJSON,
	}

	//TODO set genesiscontroller

	genesis.Poa = &genesispoa
	//	genesis.Controller = &genesiscontroller

	if alloc == nil {
		alloctmp, err := buildAlloc(keystore)
		if err != nil {
			return err
		}
		alloc = &alloctmp
	}

	genesis.Alloc = alloc

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.DebugMessage("Write Genesis.json")
	jsonFileName := filepath.Join(outDir, configuration.GenesisJSON)
	files.WriteToFile(jsonFileName, string(genesisjson), files.OverwriteSilently)

	return nil
}

// buildAlloc builds the alloc structure of the genesis file
func buildAlloc(accountsDir string) (Alloc, error) {
	var alloc = make(Alloc)

	tfiles, err := ioutil.ReadDir(accountsDir)
	if err != nil {
		return alloc, err
	}

	for _, f := range tfiles {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}

		path := filepath.Join(accountsDir, f.Name())

		// Read key from file.
		keyjson, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to read the keyfile at '%s': %v", path, err)
		}

		k := new(crypto.EncryptedKeyJSONMonet)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return nil, err
		}

		moniker := strings.TrimSuffix(f.Name(), ".json")
		balance := configuration.DefaultAccountBalance
		addr := k.Address

		rec := AllocRecord{Moniker: moniker, Balance: balance}
		alloc[addr] = &rec
	}

	return alloc, nil
}
