package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"
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
	Address string `json:"address"`
	Abi     string `json:"abi"`
	Code    string `json:"code"`
}

// file is the structure that a Genesis file gets parsed into.
type file struct {
	Alloc *Alloc `json:"alloc"`
	Poa   *POA   `json:"poa"`
}

// GenerateGenesisJSON compiles the POA solitity smart-contract with the peers
// baked into the whitelist. It then creates a genesis file with the
// corresponding POA section, and if no alloc section is provided, creates one
// with all the keys in keystore. The file is written to outDir.
func GenerateGenesisJSON(outDir, keystore string, peers []*peers.Peer, alloc *Alloc, contractAddress string) error {
	var genesis file

	finalSource, err := contract.GetFinalSoliditySource(peers)
	if err != nil {
		return err
	}

	poaDir := filepath.Join(outDir, configuration.POADir)

	files.CreateDirsIfNotExists([]string{poaDir})

	genesispoa, err := BuildPOA(
		finalSource,
		contractAddress,
		poaDir)
	if err != nil {
		return err
	}
	genesis.Poa = &genesispoa

	if alloc == nil {
		alloctmp, err := BuildAlloc(keystore)
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

// BuildAlloc builds the alloc structure of the genesis file
func BuildAlloc(accountsDir string) (Alloc, error) {
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

// BuildPOA builds the poa section of the genesis file, and saves a compilation
// report to outDir.
func BuildPOA(solidityCode string, contractAddress string, outDir string) (POA, error) {
	var poagenesis POA

	// Retrieve and set the version number
	version, err := GetSolidityCompilerVersion()
	if err != nil {
		return poagenesis, err
	}

	contractInfo, err := CompileSolidityContract(solidityCode)
	if err != nil {
		common.ErrorMessage("Error compiling genesis contract:", err)
		return poagenesis, err
	}

	poagenesis, err = BuildCompilationReport(version, contractInfo, contractAddress, solidityCode, outDir)
	if err != nil {
		common.ErrorMessage("Error writing compilation output:", err)
		return poagenesis, err
	}

	return poagenesis, nil
}
