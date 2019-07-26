package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/src/common"
	mtypes "github.com/mosaicnetworks/monetd/src/types"

	"github.com/pelletier/go-toml"

	"github.com/ethereum/go-ethereum/common/compiler"

	types "github.com/ethereum/go-ethereum/common"
)

// GenesisAllocRecord ...
type GenesisAllocRecord struct {
	Balance string `json:"balance"`
	Moniker string `json:"moniker"`
}

// GenesisAlloc ...
type GenesisAlloc map[string]*GenesisAllocRecord

// GenesisPOA ...
type GenesisPOA struct {
	Address string `json:"address"`
	Abi     string `json:"abi"`
	Code    string `json:"code"`
}

// GenesisFile ...
type GenesisFile struct {
	Alloc *GenesisAlloc `json:"alloc"`
	Poa   *GenesisPOA   `json:"poa"`
}

//GetSolidityCompilerVersion does exactly that
func GetSolidityCompilerVersion() (string, error) {

	s, err := compiler.SolidityVersion("")

	if err != nil {
		return "", err
	}

	common.DebugMessage("Path         : ", s.Path)
	common.DebugMessage("Full Version : \n", s.FullVersion)
	version := s.FullVersion
	re := regexp.MustCompile(`\r?\n`)
	version = re.ReplaceAllString(version, " ")
	return version, nil
}

//CompileSolidityContract ...
func CompileSolidityContract(soliditySource string) (map[string]*compiler.Contract, error) {
	contractInfo, err := compiler.CompileSolidityString("solc", soliditySource)
	if err != nil {
		common.ErrorMessage("Error compiling genesis contract:", err)
	}
	return contractInfo, err
}

type solidityFields struct {
	Constants string
	AddTo     string
	Checks    string
}

//BuildGenesisJSON ...
func BuildGenesisJSON(configDir string, peers mtypes.PeerRecordList, contractAddress string) error {
	var genesis GenesisFile

	finalSource, err := contract.GetFinalSoliditySource(peers)
	if err != nil {
		return err
	}

	genesispoa, err := BuildGenesisPOAJSON(finalSource, configDir, contractAddress)

	alloc, err := BuildGenesisAlloc(filepath.Join(configDir, configuration.KeyStoreDir))

	if err != nil {
		return err
	}

	genesis.Alloc = &alloc
	genesis.Poa = &genesispoa

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.DebugMessage("Write Genesis.json")
	jsonFileName := filepath.Join(configDir, configuration.EthDir, configuration.GenesisJSON)
	files.WriteToFile(jsonFileName, string(genesisjson))

	return nil
}

//BuildGenesisAlloc builds the alloc structure of the Genesis File
func BuildGenesisAlloc(accountsDir string) (GenesisAlloc, error) {
	var alloc = make(GenesisAlloc)

	tfiles, err := ioutil.ReadDir(accountsDir)
	if err != nil {
		return alloc, err
	}

	for i, f := range tfiles {
		splits := strings.Split(f.Name(), ".")
		if splits[len(splits)-1] != "toml" {
			continue
		}

		tomlFile := filepath.Join(accountsDir, f.Name())

		tree, err := files.LoadToml(tomlFile)
		if err != nil {
			return alloc, err
		}
		if !tree.HasPath([]string{"node", "address"}) {
			continue
		} // Need a address
		addr := tree.GetPath([]string{"node", "address"}).(string)

		// Set defaults then overwrite if set
		balance := configuration.DefaultAccountBalance
		moniker := "node" + strconv.Itoa(i)

		if tree.HasPath([]string{"node", "moniker"}) {
			moniker = tree.GetPath([]string{"node", "moniker"}).(string)
		}
		if tree.HasPath([]string{"node", "balance"}) {
			balance = tree.GetPath([]string{"node", "balance"}).(string)
		}

		rec := GenesisAllocRecord{Moniker: moniker, Balance: balance}
		alloc[addr] = &rec
	}

	return alloc, nil
}

//BuildGenesisPOAJSON ...
func BuildGenesisPOAJSON(solidityCode string, monetdConfigDir string, contractAddress string) (GenesisPOA, error) {
	var poagenesis GenesisPOA
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

	poagenesis, err = BuildCompilationReport(version, contractInfo, filepath.Join(monetdConfigDir, configuration.EthDir, configuration.POADir), contractAddress, solidityCode)
	if err != nil {
		common.ErrorMessage("Error writing compilation output:", err)
		return poagenesis, err
	}

	return poagenesis, nil
}

//BuildCompilationReport outputs compiler results in a standard format and
//builds the poa structure that is written to the Genesis File
func BuildCompilationReport(version string, contractInfo map[string]*compiler.Contract, outputDir string, contractAddress string, solidityCode string) (GenesisPOA, error) {

	var poagenesis GenesisPOA

	eip55addr := types.HexToAddress(contractAddress).Hex()

	// Create empty tree by loading an empty string
	tree, err := toml.Load("")
	if err != nil {
		return poagenesis, errors.New("cannot create compiler results tree")
	}

	tree.SetPath([]string{"solc", "compilerversion"}, version)
	tree.SetPath([]string{"solc", "os"}, runtime.GOOS)
	tree.SetPath([]string{"solc", "arch"}, runtime.GOARCH)

	for k, v := range contractInfo {
		common.DebugMessage("Processing Contract: ", k)
		jsonabi, err := json.MarshalIndent(v.Info.AbiDefinition, "", "\t")
		if err != nil {
			common.ErrorMessage("ABI error:", err)
			return poagenesis, err
		}

		tree.SetPath([]string{"poa", "contractclass"}, strings.TrimPrefix(k, "<stdin>:"))
		tree.SetPath([]string{"poa", "abi"}, string(jsonabi))
		tree.SetPath([]string{"poa", "address"}, eip55addr)

		files.WriteToFile(filepath.Join(outputDir, configuration.GenesisABI), string(jsonabi))

		tree.SetPath([]string{"poa", "bytecode"}, strings.TrimPrefix(v.RuntimeCode, "0x"))

		poagenesis.Abi = string(jsonabi)
		poagenesis.Address = eip55addr //EIP55 compliant
		poagenesis.Code = strings.TrimPrefix(v.RuntimeCode, "0x")

		break
		// We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	files.WriteToFile(filepath.Join(outputDir, configuration.GenesisContract), solidityCode)

	files.SaveToml(tree, filepath.Join(outputDir, configuration.CompileResultFile))

	return poagenesis, nil
}
