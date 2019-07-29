package network

import (
	"encoding/hex"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables set by command line parameters
var (
	namesFile       string
	passFile        string
	initIP          string
	initPeers       = 0
	generatePassKey = false
	savePassKey     = false
)

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [network_name]",
		Short: "new configuration for a multi-node network",
		Long: `
giverny network build
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkNew,
	}

	addNewFlags(cmd)

	return cmd
}

func addNewFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&namesFile, "names", "", "filename of a file containing a list of node monikers")
	cmd.Flags().StringVar(&passFile, "pass", "", "filename of a file containing a passphrase")
	cmd.Flags().IntVar(&initPeers, "initial-peers", initPeers, "number of initial peers")
	cmd.Flags().StringVar(&initIP, "initial-ip", "", "initial IP address of range")
	cmd.Flags().BoolVar(&generatePassKey, "generate-pass", generatePassKey, "generate pass phrases")
	cmd.Flags().BoolVar(&savePassKey, "save-pass", savePassKey, "save pass phrase entered on command line")

	viper.BindPFlags(cmd.Flags())
}

func networkNew(cmd *cobra.Command, args []string) error {

	// First validate network name
	networkName = strings.TrimSpace(args[0])

	if (passFile != "") && (generatePassKey) {
		return errors.New("incompatible options --pass and --generate-pass")
	}

	if !common.CheckMoniker(networkName) {
		return errors.New("network name must only contains characters in the range 0-9 or A-Z or a-z")
	}

	// Check if already exists; if does, abort

	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	if files.CheckIfExists(networkDir) {
		return errors.New("network already exists: " + networkDir)
	}

	// Create folders for this node
	keystoreDir := filepath.Join(networkDir, givernyKeystoreDir)
	files.CreateDirsIfNotExists([]string{
		configuration.GivernyConfigDir,
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir),
		networkDir,
		keystoreDir,
	})

	common.InfoMessage("Generate " + strconv.Itoa(numberOfNodes) + " Nodes with " +
		strconv.Itoa(initPeers) + " initial peers")

	// Get Node list
	nodeList, err := getNodesWithNames(namesFile, numberOfNodes, initPeers, initIP)
	if err != nil {
		return err
	}

	// Generate Keys
	for j, n := range nodeList {

		var thisNodePassPhraseFile = filepath.Join(keystoreDir, n.Moniker+".txt")

		if generatePassKey {
			passphrase := crypto.RandomPassphrase(8)
			files.WriteToFile(thisNodePassPhraseFile, passphrase)
			common.DebugMessage("Written " + thisNodePassPhraseFile)
		} else {
			if savePassKey {
				if passFile != "" {
					files.CopyFileContents(passFile, thisNodePassPhraseFile)
					common.DebugMessage("copied " + passFile + " to " + thisNodePassPhraseFile)
				} else {
					passphrase, _ := crypto.GetPassphrase("", true)
					files.WriteToFile(thisNodePassPhraseFile, passphrase)
				}
			} else {
				thisNodePassPhraseFile = passFile
			}

		}
		common.InfoMessage("Generating Key for " + n.Moniker + " (" + strconv.FormatBool(n.Validator) + ") " + n.NetAddr)

		//TODO add a save passphrase option.
		key, err := crypto.NewKeyPair(networkDir, n.Moniker, thisNodePassPhraseFile)
		if err != nil {
			return err
		}
		nodeList[j].Address = key.Address.Hex()
		nodeList[j].PubKeyHex = hex.EncodeToString(ecrypto.FromECDSAPub(&key.PrivateKey.PublicKey))

		common.DebugMessage("   " + n.Address)
	}

	//TODO remove this loop, it is just debug verification code
	/*	for j, n := range nodeList {
			fmt.Println(strconv.Itoa(j) + " " + n.PubKeyHex)
		}
	*/
	return nil
}
