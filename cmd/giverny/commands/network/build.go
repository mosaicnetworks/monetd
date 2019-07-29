package network

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mosaicnetworks/monetd/src/config"
	"github.com/mosaicnetworks/monetd/src/contract"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/types"
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [network_name]",
		Short: "create the configuration for a multi-node network",
		Long: `
giverny network build
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkBuild,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkBuild(cmd *cobra.Command, args []string) error {
	return buildNetwork(strings.TrimSpace(args[0]))
}

//buildNetwork builds the network. It is called directly from new command as well.
func buildNetwork(networkName string) error {
	if !common.CheckMoniker(networkName) {
		return errors.New("network name, " + networkName + ", is invalid")
	}

	// Check all the files and directories we expect actually exist
	thisNetworkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)

	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	tree, err := files.LoadToml(networkTomlFile)
	if err != nil {
		common.ErrorMessage("Cannot load network.toml file: ", networkTomlFile)
		return err
	}

	common.DebugMessage("Building network " + networkName)

	err = dumpPeersJSON(tree, thisNetworkDir)
	if err != nil {
		common.ErrorMessage("Error writing peers json file")
		return err
	}

	return nil
}

func dumpPeersJSON(tree *toml.Tree, thisNetworkDir string) error {

	var peers types.PeerRecordList

	if tree.HasPath([]string{"Name"}) {
		netName := tree.GetPath([]string{"Name"}).(string)
		common.DebugMessage("Network Name ", netName)
	}

	common.DebugMessage("dumpPeersJSON")

	nodesquery, err := query.CompileAndExecute("$.nodes", tree)
	if err != nil {
		common.ErrorMessage("Error loading nodes")
		return err
	}

	for _, value := range nodesquery.Values() {

		//		common.DebugMessage(reflect.TypeOf(value).String())
		//		common.DebugMessage("Found a value: "+strconv.Itoa(i), value)

		if reflect.TypeOf(value).String() == "[]*toml.Tree" {
			nodes := value.([]*toml.Tree)

			for _, tr := range nodes {
				var moniker, netaddr, pubkey string

				if tr.HasPath([]string{"validator"}) && (!tr.GetPath([]string{"validator"}).(bool)) {
					continue
				}

				if tr.HasPath([]string{"moniker"}) {
					moniker = tr.GetPath([]string{"moniker"}).(string)
				}
				if tr.HasPath([]string{"netaddr"}) {
					netaddr = tr.GetPath([]string{"netaddr"}).(string)
				}
				if tr.HasPath([]string{"pubkey"}) {
					pubkey = tr.GetPath([]string{"pubkey"}).(string)
				}

				peers = append(peers, &types.PeerRecord{Moniker: moniker,
					NetAddr:   netaddr,
					PubKeyHex: pubkey})
			}

		}
	}

	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(thisNetworkDir, monetconfig.PeersJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut))
	if err != nil {
		return err
	}

	err = BuildGenesisJSON(thisNetworkDir, peers, monetconfig.DefaultContractAddress)
	if err != nil {
		return err
	}

	return err
}

//BuildGenesisJSON compiles and build a genesis.json file
func BuildGenesisJSON(configDir string, peers types.PeerRecordList, contractAddress string) error {
	var genesis config.GenesisFile

	common.DebugMessage("buildGenesisJSON")

	finalSource, err := contract.GetFinalSoliditySource(peers)
	if err != nil {
		return err
	}

	common.DebugMessage("Source Loaded")

	genesispoa, err := config.BuildGenesisPOAJSON(finalSource, configDir, contractAddress, false)
	if err != nil {
		return err
	}
	genesis.Poa = &genesispoa

	common.DebugMessage("POA Section Build")

	//TODO source the Token values from the genesis file.

	alloc, err := config.BuildGenesisAlloc(filepath.Join(configDir, monetconfig.KeyStoreDir))
	if err != nil {
		return err
	}
	genesis.Alloc = &alloc

	common.DebugMessage("Alloc Built")

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.DebugMessage("Write Genesis.json")
	jsonFileName := filepath.Join(configDir, monetconfig.GenesisJSON)
	files.WriteToFile(jsonFileName, string(genesisjson))

	return nil
}
