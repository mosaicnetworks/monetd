package genesis

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	toml "github.com/pelletier/go-toml"
)

// GetSolidityCompilerVersion gets the version of the solidity compiler that
// comes with Geth.
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

// CompileSolidityContract compiles a solitity smart-contract using the compiler
// that comes with Geth.
func CompileSolidityContract(soliditySource string) (map[string]*compiler.Contract, error) {
	contractInfo, err := compiler.CompileSolidityString("solc", soliditySource)
	if err != nil {
		common.ErrorMessage("Error compiling genesis contract:", err)
	}
	return contractInfo, err
}

// BuildCompilationReport outputs compiler results in a standard format and
// builds the poa structure that is written to the Genesis File
func BuildCompilationReport(version string,
	contractInfo map[string]*compiler.Contract,
	contractAddress string,
	solidityCode string,
	outputDir string) (POA, error) {

	var poagenesis POA

	eip55addr := eth_common.HexToAddress(contractAddress).Hex()

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

		files.WriteToFile(
			filepath.Join(outputDir, configuration.GenesisABI),
			string(jsonabi),
			files.OverwriteSilently,
		)

		tree.SetPath([]string{"poa", "bytecode"}, strings.TrimPrefix(v.RuntimeCode, "0x"))

		poagenesis.Abi = string(jsonabi)
		poagenesis.Address = eip55addr //EIP55 compliant
		poagenesis.Code = strings.TrimPrefix(v.RuntimeCode, "0x")

		break
		// We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	files.WriteToFile(
		filepath.Join(outputDir, configuration.GenesisContract),
		solidityCode,
		files.OverwriteSilently)

	files.SaveToml(tree, filepath.Join(outputDir, configuration.CompileResultFile))

	return poagenesis, nil
}
