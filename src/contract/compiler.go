package contract

import (
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/mosaicnetworks/monetd/src/common"
)

//CompileSolidityContract is a simple wrapper for the Geth compiler function.
//It takes the Solidity source as a string input parameter.
func CompileSolidityContract(soliditySource string) (map[string]*compiler.Contract, error) {
	contractInfo, err := compiler.CompileSolidityString("solc", soliditySource)
	if err != nil {
		common.ErrorMessage("Error compiling genesis contract:", err)
	}
	return contractInfo, err
}
