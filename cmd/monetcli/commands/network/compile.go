package network

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"

	types "github.com/ethereum/go-ethereum/common"
	compile "github.com/ethereum/go-ethereum/common/compiler"

	"github.com/spf13/cobra"
)

func compileConfig(cmd *cobra.Command, args []string) error {
	return CompileConfigWithParam(configDir)
}

//CompileConfigWithParam "finishes" the monetcli configuration, compiling the POA smart contract
//in preparation for a call to monetcli config publish
func CompileConfigWithParam(configDir string) error {
	var soliditySource string
	// Load the Current Config

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	// Retrieve and set the version number
	s, err := compile.SolidityVersion("")

	if err != nil {
		return err
	}

	common.Message("Path         : ", s.Path)
	message("Full Version : \n", s.FullVersion)
	version := s.FullVersion
	re := regexp.MustCompile(`\r?\n`)
	version = re.ReplaceAllString(version, " ")

	tree.SetPath([]string{"poa", "compilerverison"}, version)

	//When contracts are "set" for a network, the solidity source is copied into the monetcli config directory
	//with a name of template.sol (defined by constant common.TemplateContract). Thus we can check just for that file.
	//If not found, then we download a fresh contract.
	filename := filepath.Join(configDir, common.TemplateContract)
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
	} else { // NB, we do not write the downloaded template to file. Preferable to get fresh is regenerating.
		message("Loading: ", common.DefaultSolidityContract)
		resp, err := http.Get(common.DefaultSolidityContract)
		if err != nil {
			message("Error loading: ", common.DefaultSolidityContract)
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

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

	if len(currentNodes) < 1 {
		return errors.New("Peerset is empty")
	}

	var consts, addTo, checks []string

	var alloc = make(genesisAlloc)
	var peers peerRecordList
	var genesisPeers peerRecordList

	for i, value := range currentNodes {

		rawaddr := tree.GetPath([]string{"validators", value, "address"}).(string)
		rawmoniker := tree.GetPath([]string{"validators", value, "moniker"}).(string)
		rawpubkey := tree.GetPath([]string{"validators", value, "pubkey"}).(string)
		rawisvalidator := tree.GetPath([]string{"validators", value, "validator"}).(bool)
		//	rawip := tree.GetPath([]string{"validators", value, "ip"}).(string)

		// Convert Hex to Address and back out to get a EIP55 compliant address
		addr := types.HexToAddress(rawaddr).Hex()

		/*	val, err := strconv.ParseBool(rawisvalidator)
			if err != nil {
				return err
			} */
		// Non-validators are added to the peer set, but not to the genesis peer set.
		peer := peerRecord{NetAddr: rawaddr, PubKeyHex: rawpubkey, Moniker: rawmoniker}
		peers = append(peers, &peer)

		if rawisvalidator {
			consts = append(consts, "    address constant initWhitelist"+strconv.Itoa(i)+" = "+addr+";")
			consts = append(consts, "    bytes32 constant initWhitelistMoniker"+strconv.Itoa(i)+" = \""+rawmoniker+"\";")

			addTo = append(addTo, "     addToWhitelist(initWhitelist"+strconv.Itoa(i)+", initWhitelistMoniker"+strconv.Itoa(i)+");")
			checks = append(checks, " ( initWhitelist"+strconv.Itoa(i)+" == _address ) ")
			genesisPeers = append(genesisPeers, &peer)
		}

		rec := genesisAllocRecord{Moniker: rawmoniker, Balance: common.DefaultAccountBalance}
		alloc[addr] = &rec

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

	err = common.WriteToFile(filepath.Join(configDir, common.GenesisContract), finalSoliditySource)
	if err != nil {
		message("Error writing genesis contract:", err)
		return err
	}

	contractInfo, err := compile.CompileSolidityString("solc", finalSoliditySource)
	var poagenesis genesisPOA

	// message("Contract Compiled: ", contractInfo)

	for k, v := range contractInfo {
		message("Processing Contract: ", k)
		jsonabi, err := json.MarshalIndent(v.Info.AbiDefinition, "", "\t")
		if err != nil {
			message("ABI error:", err)
			return err
		}

		tree.SetPath([]string{"poa", "contractclass"}, strings.TrimPrefix(k, "<stdin>:"))
		tree.SetPath([]string{"poa", "abi"}, string(jsonabi))

		common.WriteToFile(filepath.Join(configDir, common.GenesisABI), string(jsonabi))
		tree.SetPath([]string{"poa", "bytecode"}, strings.TrimPrefix(v.RuntimeCode, "0x"))

		poagenesis.Abi = string(jsonabi)
		poagenesis.Address = types.HexToAddress(tree.Get("poa.contractaddress").(string)).Hex() //EIP55 compliant
		poagenesis.Code = strings.TrimPrefix(v.RuntimeCode, "0x")

		message("Set Contract Items")
		break // We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	err = common.SaveTomlConfig(configDir, tree)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Cannot save TOML file")
		return err
	}

	var genesis genesisFile

	genesis.Alloc = &alloc
	genesis.Poa = &poagenesis

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(configDir, common.GenesisJSON)
	common.WriteToFile(jsonFileName, string(genesisjson))

	peersjson, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}
	jsonFileName = filepath.Join(configDir, common.PeersJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	peersjson, err = json.MarshalIndent(genesisPeers, "", "\t")
	if err != nil {
		return err
	}
	jsonFileName = filepath.Join(configDir, common.PeersGenesisJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	return nil
}
