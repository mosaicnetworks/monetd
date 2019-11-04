package genesis

import (
	"bytes"
	"text/template"

	eth_common "github.com/ethereum/go-ethereum/common"
)

//GetControllerSoliditySource returns the source for the contract that returns
//where the POA contract is hosted.
func GetControllerSoliditySource(contractAddress string) (string, error) {

	const templateSol = `pragma solidity >=0.4.22;

	/// @title Proof of Authority Whitelist Proof of Concept
	/// @author Jon Knight
	/// @author Mosaic Networks
	/// @notice Copyright Mosaic Networks 2019, released under the MIT license
	
	contract Monet_Controller {
		
		
		address poaContract = {{.ContractAddress}};
		
		
		
		
		///@notice This function returns the current POA Contract Address
		///@return returns the current POA Contract Address
		function POAContractAddress() public view returns (address contractAddress)
		{
			return poaContract;
		}    
			
	} `

	type controllerData struct {
		ContractAddress string
	}

	eip55Controlleraddr := eth_common.HexToAddress(contractAddress).Hex()

	conStruct := controllerData{
		ContractAddress: eip55Controlleraddr}

	templ, err := template.New("controllersolidity").Parse(templateSol)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	templ.Execute(buf, conStruct)

	return buf.String(), nil
}
