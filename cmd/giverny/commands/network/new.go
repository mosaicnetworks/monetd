package network

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"

	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/evm-lite/src/currency"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"

	mconfiguration "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const nodeNamePrefix = "node"

// Variables set by command line parameters
var (
	_namesFile       string
	_passFile        string
	_initIP          string
	_numberOfNodes   = -1
	_initPeers       = 0
	_generatePassKey = false
	_noSavePassKey   = false
	_noBuild         = false
)

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [network_name]",
		Short: "new configuration for a multi-node network",
		Args:  cobra.ExactArgs(1),
		RunE:  networkNew,
	}

	addNewFlags(cmd)

	return cmd
}

func addNewFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&_namesFile, "names", "", "file containing node configurations")
	cmd.Flags().StringVar(&_passFile, "pass", "", "file containing a passphrase")
	cmd.Flags().IntVar(&_initPeers, "initial-peers", _initPeers, "number of initial peers")
	cmd.Flags().StringVar(&_initIP, "initial-ip", "", "initial IP address of range")
	cmd.Flags().BoolVar(&_generatePassKey, "generate-pass", _generatePassKey, "generate pass phrases")
	cmd.Flags().BoolVar(&_noSavePassKey, "no-save-pass", _noSavePassKey, "don't save passphrase entered on command line")
	cmd.Flags().BoolVar(&_noBuild, "no-build", _noBuild, "disables the automatic build of a new network")
	cmd.Flags().IntVarP(&_numberOfNodes, "nodes", "n", _numberOfNodes, "number of nodes in this configuration")

	viper.BindPFlags(cmd.Flags())
}

func networkNew(cmd *cobra.Command, args []string) error {

	// First validate network name
	networkName = strings.TrimSpace(args[0])

	if (_passFile != "") && (_generatePassKey) {
		return errors.New("incompatible options --pass and --generate-pass")
	}

	if _namesFile == "" && _numberOfNodes < 1 {
		return errors.New("incompatible options you must specify --nodes or --names")
	}

	if !common.CheckMoniker(networkName) {
		return errors.New("network name must only contains characters in the range 0-9 or A-Z or a-z")
	}

	// Check if already exists; if does, abort
	networkDir := filepath.Join(configuration.GivernyConfigDir, configuration.GivernyNetworkDir, networkName)
	if files.CheckIfExists(networkDir) {
		return errors.New("network already exists: " + networkDir)
	}

	// Get node list
	nodeList, err := createNodes(_namesFile, _numberOfNodes, _initPeers, _initIP)
	if err != nil {
		return err
	}

	// Generate keys and populate keystore. This updates the nodeList with
	// addresses and public keys.
	err = generateKeys(
		filepath.Join(networkDir, givernyKeystoreDir),
		nodeList)
	if err != nil {
		return err
	}

	// Build and write network.toml
	err = generateGivernyConfig(
		filepath.Join(networkDir, networkTomlFileName),
		networkName,
		_initIP,
		nodeList,
	)
	if err != nil {
		return err
	}

	// Write default monetd.toml file
	mconfiguration.DumpGlobalTOML(networkDir, mconfiguration.MonetTomlFile)

	if _noBuild {
		return nil
	}

	return buildNetwork(networkName)
}

// Generate a list of nodes based on a configuration file (names) or just
// iterating from a base
func createNodes(srcFile string, numNodes int, numValidators int, initialIP string) ([]node, error) {
	var rtn []node
	ipStem := ""
	lastDigit := 0
	var err error

	if initialIP != "" {
		splitIP := strings.Split(initialIP, ".")

		if len(splitIP) != 4 {
			return rtn, errors.New("malformed initial IP: " + initialIP)
		}

		lastDigit, err = strconv.Atoi(splitIP[3])
		if err != nil {
			fmt.Println("lastDigit Set to Zero")
			lastDigit = 0
		} else {
			ipStem = strings.Join(splitIP[:3], ".") + "."
		}
	}

	if srcFile == "" {
		return createNodesIteratively(numNodes, numValidators, ipStem, lastDigit)
	}

	return createNodesFromNamesFile(srcFile, numNodes, numValidators, ipStem, lastDigit)
}

// Create nodes with names and IPs generated iteratively from a base.
func createNodesIteratively(numNodes int, numValidators int, ipStem string, lastDigit int) ([]node, error) {
	var rtn []node

	for i := 0; i < numNodes; i++ {
		node := node{
			Moniker:   nodeNamePrefix + strconv.Itoa(i),
			Validator: (numValidators < 1 || i < numValidators),
			Tokens:    defaultTokens,
		}

		if ipStem != "" {
			node.NetAddr = ipStem + strconv.Itoa(lastDigit+i)
		}

		rtn = append(rtn, node)
	}

	return rtn, nil
}

// Create nodes based on names-file definiition
func createNodesFromNamesFile(srcFile string, numNodes int, numValidators int, ipStem string, lastDigit int) ([]node, error) {
	var rtn []node

	file, err := os.Open(srcFile)
	if err != nil {
		common.ErrorMessage("failed opening file: ", err)
		return rtn, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Read file line by line.
	i := 1
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// Ignore blank lines
		if line == "" {
			continue
		}

		// Ignore comments
		if line[:1] == "#" {
			continue
		}

		var moniker string
		var netaddr string
		var validator bool
		var tokens string
		var nonnode bool

		if ipStem != "" {
			netaddr = ipStem + strconv.Itoa(lastDigit+i-1)
		}

		validator = (numValidators < 1 || i <= numValidators)
		tokens = defaultTokens
		nonnode = false

		if strings.Contains(line, ",") {
			arrLine := strings.Split(line, ",")

			moniker = arrLine[0]
			common.DebugMessage("Setting moniker to " + moniker)

			if len(arrLine) > 1 {
				if strings.TrimSpace(arrLine[1]) != "" {
					netaddr = arrLine[1]
				}
			}

			if len(arrLine) > 2 {
				if strings.TrimSpace(arrLine[2]) != "" {
					tokens = currency.ExpandCurrencyString(arrLine[2])
				}
			}

			if len(arrLine) > 3 && len(strings.TrimSpace(arrLine[3])) > 0 {
				validator, _ = strconv.ParseBool(arrLine[3])
			}

			if len(arrLine) > 4 && len(strings.TrimSpace(arrLine[4])) > 0 {
				nonnode, _ = strconv.ParseBool(arrLine[4])
				if nonnode {
					netaddr = ""
				}
			}
		} else {
			moniker = line
		}

		if !common.CheckMoniker(moniker) {
			return rtn, errors.New("node name " + moniker + " contains invalid characters")
		}

		rtn = append(rtn,
			node{
				Moniker:   moniker,
				NetAddr:   netaddr,
				Validator: validator,
				Tokens:    tokens,
				NonNode:   nonnode,
			})

		if i >= numNodes && numNodes > 0 {
			break
		}

		i++
	}

	return rtn, nil
}

// Generate keys for each node, update the nodeList
func generateKeys(keystore string, nodeList []node) error {
	// Create folders for this node
	files.CreateDirsIfNotExists([]string{
		configuration.GivernyConfigDir,
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir),
		keystore,
	})

	// prompt password and write it to a file, if necessary
	if !_noSavePassKey && _passFile == "" && !_generatePassKey {
		passphrase, _ := crypto.GetPassphrase("", true)
		_passFile = filepath.Join(keystore, "pwd.txt")
		files.WriteToFile(_passFile, passphrase)
	}

	// Generate Keys
	for j, n := range nodeList {

		var thisNodePassPhraseFile = filepath.Join(keystore, n.Moniker+".txt")

		if _generatePassKey {
			passphrase := crypto.RandomPassphrase(8)
			files.WriteToFile(thisNodePassPhraseFile, passphrase, files.BackupExisting)
			common.DebugMessage("Written " + thisNodePassPhraseFile)
		} else {
			if !_noSavePassKey {
				if _passFile != "" {
					files.CopyFileContents(_passFile, thisNodePassPhraseFile)
					common.DebugMessage("copied " + _passFile + " to " + thisNodePassPhraseFile)
				} else {
					passphrase, _ := crypto.GetPassphrase("", true)
					files.WriteToFilePrivate(thisNodePassPhraseFile, passphrase)
				}
			} else {
				thisNodePassPhraseFile = _passFile
			}
		}

		common.InfoMessage(fmt.Sprintf("Generating key for %s (validator? %v) %s",
			n.Moniker,
			n.Validator,
			n.NetAddr))

		//TODO add a save passphrase option.
		key, err := crypto.NewKeyfile(keystore, n.Moniker, thisNodePassPhraseFile)
		if err != nil {
			return err
		}

		nodeList[j].Address = key.Address.Hex()
		nodeList[j].PubKeyHex = hex.EncodeToString(ecrypto.FromECDSAPub(&key.PrivateKey.PublicKey))

		common.DebugMessage("   " + n.Address)
	}

	return nil
}

// Build and write the network.toml file, which is used by giverny to manage
// networks.
func generateGivernyConfig(configFile, networkName, initIP string, nodeList []node) error {
	conf := Config{
		Network: networkConfig{Name: networkName},
		Nodes:   nodeList,
		Docker:  dockerConfig{Name: networkName},
	}

	if initIP != "" {
		conf.Docker.BaseIP = initIP

		arrIP := strings.Split(initIP, ".")
		if len(arrIP) > 3 {
			conf.Docker.Subnet = strings.Join(arrIP[:2], ".") + ".0.0/16"
			conf.Docker.IPRange = conf.Docker.Subnet
			conf.Docker.Gateway = strings.Join(arrIP[:3], ".") + ".254"
		}
	}

	tomlBytes, err := toml.Marshal(conf)
	if err != nil {
		return err
	}

<<<<<<< HEAD
	err = files.WriteToFile(filepath.Join(networkDir, networkTomlFileName), string(tomlBytes), files.BackupExisting)
=======
	err = files.WriteToFile(configFile, string(tomlBytes))
>>>>>>> filesystem-layout
	if err != nil {
		return err
	}

	return nil
}
