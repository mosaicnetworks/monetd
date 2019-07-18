package common

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/compiler"
)

//SolidityResponse is the return structure from BuildGenesisJson()
type SolidityResponse struct {
	SolidityVersion string
}

//GetSolidityCompilerVersion does exactly that
func GetSolidityCompilerVersion() (string, error) {

	s, err := compiler.SolidityVersion("")

	if err != nil {
		return "", err
	}

	MessageWithType(MsgDebug, "Path         : ", s.Path)
	MessageWithType(MsgDebug, "Full Version : \n", s.FullVersion)
	version := s.FullVersion
	re := regexp.MustCompile(`\r?\n`)
	version = re.ReplaceAllString(version, " ")
	return version, nil
}

//GetSoliditySource ...
func GetSoliditySource(filename string) (string, error) {
	var soliditySource string

	if _, err := os.Stat(filename); err == nil {
		MessageWithType(MsgDebug, "Opening: ", filename)
		file, err := os.Open(filename)
		if err != nil {
			MessageWithType(MsgError, "Error opening: ", filename)
			return "", err
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			MessageWithType(MsgError, "Error reading: ", filename)
			return "", err
		}

		soliditySource = string(b)
	} else { // NB, we do not write the downloaded template to file. Preferable to get fresh is regenerating.
		MessageWithType(MsgDebug, "Loading: ", DefaultSolidityContract)
		resp, err := http.Get(DefaultSolidityContract)
		if err != nil {
			MessageWithType(MsgError, "Could not load the standard poa smart contract from GitHub. Aborting.")
			MessageWithType(MsgError, "You can specify the contract explicitly using the standard one from [...monetd]/smart-contract/genesis.sol")
			MessageWithType(MsgInformation, "monetcli network contract [...monetd]/smart-contract/genesis.sol")

			MessageWithType(MsgError, "Error loading: ", DefaultSolidityContract)

			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			MessageWithType(MsgError, "Error reading body of Solidity Contract")
			return "", err
		}

		soliditySource = string(body)
	}

	return soliditySource, nil
}

//CompileSolidityContract ...
func CompileSolidityContract(soliditySource string) (map[string]*compiler.Contract, error) {
	contractInfo, err := compiler.CompileSolidityString("solc", soliditySource)
	if err != nil {
		MessageWithType(MsgError, "Error compiling genesis contract:", err)
	}
	return contractInfo, err
}

//ApplyInitialWhitelistToSoliditySource ...
func ApplyInitialWhitelistToSoliditySource(soliditySource string, peers PeerRecordList) (string, error) {

	var consts, addTo, checks []string

	for i, peer := range peers {
		addr, err := PublicKeyHexToAddressHex(peer.PubKeyHex)
		if err != nil {
			return "", err
		}

		consts = append(consts, "    address constant initWhitelist"+strconv.Itoa(i)+" = "+addr+";")
		consts = append(consts, "    bytes32 constant initWhitelistMoniker"+strconv.Itoa(i)+" = \""+peer.Moniker+"\";")

		addTo = append(addTo, "     addToWhitelist(initWhitelist"+strconv.Itoa(i)+", initWhitelistMoniker"+strconv.Itoa(i)+");")
		checks = append(checks, " ( initWhitelist"+strconv.Itoa(i)+" == _address ) ")
	}

	generatedSol := "GENERATED GENESIS BEGIN \n " +
		" \n" +
		strings.Join(consts, "\n") +
		" \n" +
		" \n" +
		" \n" +
		"    function processGenesisWhitelist() private \n" +
		"    { \n" +
		strings.Join(addTo, "\n") +
		" \n" +
		"    } \n" +
		" \n" +
		" \n" +
		"    function isGenesisWhitelisted(address _address) pure private returns (bool) \n" +
		"    { \n" +
		"        return ( " + strings.Join(checks, "||") + "); \n" +
		"    } \n" +

		" \n" +
		" //GENERATED GENESIS END \n "

	// replace

	reg := regexp.MustCompile(`(?s)GENERATED GENESIS BEGIN.*GENERATED GENESIS END`)
	finalSoliditySource := reg.ReplaceAllString(soliditySource, generatedSol)

	return finalSoliditySource, nil
}

//BuildGenesisJSON ...
func BuildGenesisJSON(solidityCode string) (SolidityResponse, error) {
	rtn := SolidityResponse{}
	// Retrieve and set the version number
	version, err := GetSolidityCompilerVersion()
	if err != nil {
		return rtn, err
	}

	rtn.SolidityVersion = version

	//	MessageWithType(MsgDebug, "checking for File : ", solContractFileName)

	return rtn, nil
}
