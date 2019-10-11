package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	bpeers "github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/monetd/src/genesis"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [network_name]",
		Short: "build monetd configuration files based on a giverny network file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkBuild,
	}
	return cmd
}

func networkBuild(cmd *cobra.Command, args []string) error {
	return buildNetwork(strings.TrimSpace(args[0]))
}

// buildNetwork builds the network. It is called directly from the "new" command
// as well.
func buildNetwork(networkName string) error {
	if !common.CheckMoniker(networkName) {
		return errors.New("network name, " + networkName + ", is invalid")
	}

	// Check all the files and directories we expect actually exist
	thisNetworkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)
	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	var conf = Config{}

	tomlbytes, err := ioutil.ReadFile(networkTomlFile)
	if err != nil {
		return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
	}

	err = toml.Unmarshal(tomlbytes, &conf)
	if err != nil {
		return fmt.Errorf("Failed to parse the toml file at '%s': %v", networkTomlFile, err)
	}

	common.DebugMessage("Building network " + networkName)

	err = generateMonetConfig(&conf, thisNetworkDir)
	if err != nil {
		common.ErrorMessage("Error writing peers json file")
		return err
	}

	return nil
}

func generateMonetConfig(conf *Config, thisNetworkDir string) error {

	var peers []*bpeers.Peer
	var alloc = make(genesis.Alloc)

	for _, n := range conf.Nodes {

		netaddr := n.NetAddr
		if !strings.Contains(netaddr, ":") {
			netaddr += ":" + monetconfig.DefaultGossipPort
		}

		rec := genesis.AllocRecord{Moniker: n.Moniker, Balance: n.Tokens}
		alloc[n.Address] = &rec

		if !n.Validator || n.NonNode {
			continue
		}

		peers = append(peers, bpeers.NewPeer(n.PubKeyHex, netaddr, n.Moniker))
	}

	err := generateBabbleFiles(thisNetworkDir, peers)
	if err != nil {
		return err
	}

<<<<<<< HEAD
	jsonFileName := filepath.Join(thisNetworkDir, monetconfig.PeersJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut), files.PromptIfExisting|files.BackupExisting)
	if err != nil {
		return err
	}

	// Write copy of peers.json to peers.genesis.json
	jsonFileName = filepath.Join(thisNetworkDir, monetconfig.PeersGenesisJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut), files.OverwriteSilently)
	if err != nil {
		return err
	}

	err = BuildGenesisJSON(thisNetworkDir, peers, monetconfig.DefaultContractAddress, alloc)
=======
	err = genesis.GenerateGenesisJSON(thisNetworkDir,
		"",
		peers,
		&alloc,
		monetconfig.DefaultContractAddress)
>>>>>>> filesystem-layout
	if err != nil {
		return err
	}

	return err
}

<<<<<<< HEAD
func createKeyFileIfNotExists(configDir string, moniker string, addr string, pubkey string) error {
	keyfile := filepath.Join(configDir, monetconfig.KeyStoreDir, moniker+".json")
	if files.CheckIfExists(keyfile) {
		return nil
	} // If exists, nothing to do

	type minjson struct {
		Address string `json:"address"`
		Pub     string `json:"pub"`
	}

	j := minjson{Address: addr, Pub: pubkey}
	out, err := json.Marshal(j)
	if err != nil {
		return err
	}

	err = files.WriteToFile(keyfile, string(out), files.BackupExisting|files.PromptIfExisting)
	if err != nil {
		return err
	}

	return nil
}

//BuildGenesisJSON compiles and build a genesis.json file
func BuildGenesisJSON(configDir string, peers types.PeerRecordList, contractAddress string, alloc config.GenesisAlloc) error {
	var genesis config.GenesisFile

	common.DebugMessage("buildGenesisJSON")

	finalSource, err := contract.GetFinalSoliditySource(peers)
=======
func generateBabbleFiles(configDir string, peers []*bpeers.Peer) error {
	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
>>>>>>> filesystem-layout
	if err != nil {
		return err
	}

	// write peers.json
	jsonFileName := filepath.Join(configDir, monetconfig.PeersJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut))
	if err != nil {
		return err
	}

	// Write peers.genesis.json
	jsonFileName = filepath.Join(configDir, monetconfig.PeersGenesisJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut))
	if err != nil {
		return err
	}

<<<<<<< HEAD
	common.DebugMessage("Write Genesis.json")
	jsonFileName := filepath.Join(configDir, monetconfig.GenesisJSON)
	files.WriteToFile(jsonFileName, string(genesisjson), files.BackupExisting)

=======
>>>>>>> filesystem-layout
	return nil
}
