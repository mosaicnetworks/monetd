package network

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [network_name] [node_name] [moniker]",
		Short: "Add a node",
		Long: `
giverny network add
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkAdd,
	}

	addAddFlags(cmd)

	return cmd
}

func addAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&passFile, "pass", "", "filename of a file containing a passphrase")
	cmd.Flags().StringVar(&initIP, "initial-ip", "", "IP")
	cmd.Flags().BoolVar(&generatePassKey, "generate-pass", generatePassKey, "generate pass phrases")
	cmd.Flags().BoolVar(&noSavePassKey, "no-save-pass", noSavePassKey, "don't save pass phrase entered on command line")
	viper.BindPFlags(cmd.Flags())
}

func networkAdd(cmd *cobra.Command, args []string) error {
	if (passFile != "") && (generatePassKey) {
		return errors.New("incompatible options --pass and --generate-pass")
	}

	return addNodeToNetwork(args[0], args[1])
}

func addNodeToNetwork(networkName, moniker string) error {

	// Check moniker valid
	if !common.CheckMoniker(moniker) {
		return errors.New("the moniker, " + moniker + ", is invalid")
	}

	// Check network exists
	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	networkTomlFile := filepath.Join(networkDir, networkTomlFileName)

	if !files.CheckIfExists(networkDir) {
		return errors.New("network " + networkName + " does not exist")
	}

	// Check nodes does not exist
	keyFile := filepath.Join(networkDir, givernyKeystoreDir, moniker+".json")
	if files.CheckIfExists(keyFile) {
		return errors.New("node " + moniker + " already exists")
	}

	var conf = Config{}

	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}

	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return nil
	}

	// Get IP address for this node
	// fall through without setting initIP if anything goes wrong
	if initIP == "" {
		if conf.Docker.LastIP != "" {
			lastip := strings.Split(conf.Docker.LastIP, ".")
			if len(lastip) == 4 {
				lastOctet, err := strconv.Atoi(lastip[3])
				if err == nil {
					initIP = strings.Join(lastip[:3], ".") + "." + strconv.Itoa(lastOctet+1)
				}
			}

		}
	}

	if initIP == "" {
		return errors.New("no address specified, and none in docker config")
	}

	newNodeConf := node{
		Moniker:   moniker,
		NetAddr:   initIP,
		Validator: false,
		Tokens:    "0",
	}

	// Prompt for passphrase if nothing suitable found.

	if !noSavePassKey && passFile == "" && !generatePassKey {
		passphrase, _ := crypto.GetPassphrase("", true)
		passFile = filepath.Join(networkDir, "pwd.txt")
		files.WriteToFile(passFile, passphrase)
	}

	// Create a Key
	var thisNodePassPhraseFile = filepath.Join(networkDir, givernyKeystoreDir, moniker+".txt")

	if generatePassKey {
		passphrase := crypto.RandomPassphrase(8)
		files.WriteToFile(thisNodePassPhraseFile, passphrase)
		common.DebugMessage("Written " + thisNodePassPhraseFile)
	} else {
		if !noSavePassKey {
			if passFile != "" {
				files.CopyFileContents(passFile, thisNodePassPhraseFile)
				common.DebugMessage("copied " + passFile + " to " + thisNodePassPhraseFile)
			} else {
				passphrase, _ := crypto.GetPassphrase("", true)
				files.WriteToFilePrivate(thisNodePassPhraseFile, passphrase)
			}
		} else {
			thisNodePassPhraseFile = passFile
		}

	}
	common.InfoMessage("Generating Key for " + moniker)

	key, err := crypto.NewKeyPair(networkDir, moniker, thisNodePassPhraseFile)
	if err != nil {
		return err
	}
	newNodeConf.Address = key.Address.Hex()
	newNodeConf.PubKeyHex = hex.EncodeToString(ecrypto.FromECDSAPub(&key.PrivateKey.PublicKey))

	// Add node section
	conf.Docker.LastIP = initIP

	conf.Nodes = append(conf.Nodes, newNodeConf)

	tomlBytes, err := toml.Marshal(conf)
	if err != nil {
		return err
	}

	err = files.WriteToFile(filepath.Join(networkDir, networkTomlFileName), string(tomlBytes))
	if err != nil {
		return err
	}

	dockerDir := filepath.Join(networkDir, givernyDockerDir)
	if !files.CheckIfExists(dockerDir) {
		common.DebugMessage("Other docker containers not yet created")
		return nil
	}

	// Docker processing
	if err := exportDockerNodeConfig(networkDir, dockerDir, &newNodeConf); err != nil {
		return err
	}

	return nil
}
