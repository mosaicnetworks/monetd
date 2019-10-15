package contract

import (
	"bytes"
	"text/template"

	eth_common "github.com/ethereum/go-ethereum/common"
)

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



		///@notice This function sets the current POA Contract Address
		///@param contractAddress is the address of the POA smart contract
		function UNSAFESetPOAContractAddress(address contractAddress) public 
		{
			poaContract = contractAddress;
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
