package network

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	compile "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

func compileConfig(cmd *cobra.Command, args []string) error {
	var soliditySource string
	// Load the Current Config
	err := loadConfig()
	if err != nil {
		return err
	}

	// Retrieve and set the version number
	s, err := compile.SolidityVersion("")

	if err != nil {
		return err
	}

	message("Path         : ", s.Path)
	message("Full Version : \n", s.FullVersion)
	version := s.FullVersion
	re := regexp.MustCompile(`\r?\n`)
	version = re.ReplaceAllString(version, " ")

	networkViper.Set("poa.compilerversion", version) // version)

	// TODO - check if this has been overridden
	filename := filepath.Join(configDir, templateContract)
	message("Checking for file: ", filename)

	if _, err := os.Stat(filename); err == nil {
		message("Opening: ", filename)
		file, err := os.Open(filename)
		if err != nil {
			message("Error opening: ", filename)
			return err
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		soliditySource = string(b)
	} else {
		message("Loading: ", defaultSolidityContract)
		resp, err := http.Get(defaultSolidityContract)
		if err != nil {
			message("Error loading: ", defaultSolidityContract)
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			message("Error reading body of Solidity Contract")
			return err
		}

		soliditySource = string(body)
	}

	// message(soliditySource)

	// Get Peers
	monikers := networkViper.GetString("validators.monikers")
	addresses := networkViper.GetString("validators.addresses")
	isvalidators := networkViper.GetString("validators.isvalidator")
	ips := networkViper.GetString("validators.ips")

	// Reject if no peers set
	if monikers == "" || addresses == "" || isvalidators == "" || ips == "" {
		return errors.New("Peerset is empty")
	}

	// Parse Peers Config into Arrays
	monikerArray := strings.Split(monikers, ";")
	addressArray := strings.Split(addresses, ";")
	isvalidatorArray := strings.Split(isvalidators, ";")
	ipArray := strings.Split(ips, ";")

	// If any of the Peers arrays are of different lengths 
	if len(monikerArray) != len(addressArray) || len(addressArray) != len(isvalidatorArray) || len(isvalidatorArray) != len(ipArray) || len(ipArray) != len(monikerArray) {
		return errors.New("peers configutation is inconsistent")
	}



	var consts, addTo, checks []string

	for  i, value := range addressArray
	{
		consts = append(consts, "    address constant initWhitelist"+i+" = "+addressArray[i] +";" )
		consts = append(consts, "    bytes32 constant initWhitelistMoniker"+i+" = \""+monikerArray[i] +"\";" )


		addTo = append(addTo,   "     addToWhitelist(initWhitelist"+i+", initWhitelistMoniker"+i+");")
		checks = append(checks, " ( initWhitelist"+i+" == _address ) ")

	}


	generatedSol := " //GENERATED GENESIS BEGIN \n " +
	" \n" +	

	strings.Join(consts, "\n")+
	" \n" +	
	" \n" +	
	" \n" +	
  "    function processGenesisWhitelist() private \n" +
  "    { \n" +
  strings.Join(addTo, "\n")+
	" \n" +	
  "    } \n" +
	" \n" +	
	" \n" +	
  "    function isGenesisWhitelisted(address _address) pure private returns (bool) \n"+
  "    { \n"+
  "        return ( "+strings.Join(checks,"||") + "); \n"+
  "    } \n"+

	" \n" +	
  " //GENERATED GENESIS END \n " ;


// replace 


	finalSoliditySource := soliditySource // TODO Prepopulate the genesis whitelist from code
	writeToFile(filepath.Join(configDir, genesisContract), finalSoliditySource)

	contractInfo, err := compile.CompileSolidityString("solc", finalSoliditySource)

	for k, v := range contractInfo {

		jsonabi, err := json.MarshalIndent(v.Info.AbiDefinition, "", "\t")
		if err != nil {
			message("ABI error:", err)
			return err
		}

		networkViper.Set("poa.contractclass", strings.TrimPrefix(k, "<stdin>:"))
		networkViper.Set("poa.abi", string(jsonabi))

		writeToFile(filepath.Join(configDir, genesisABI), string(jsonabi))
		networkViper.Set("poa.bytecode", v.RuntimeCode)
		break // We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	writeConfig()
	return nil
}
