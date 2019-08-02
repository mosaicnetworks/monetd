package network

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"

	"github.com/spf13/cobra"
)

func newDumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump [network_name]",
		Short: "Dump the network settings",
		Long: `
giverny network dump
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkDump,
	}

	return cmd
}

func networkDump(cmd *cobra.Command, args []string) error {

	networkName = strings.TrimSpace(args[0])

	// Sanity check the network
	if !common.CheckMoniker(networkName) {
		return errors.New("the network name, " + networkName + ", is invalid")
	}

	if !files.CheckIfExists(filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)) {
		return errors.New("the network, " + networkName + " has not been created")
	}

	networkTomlFile := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName, networkTomlFileName)
	// Load Toml file to a tree
	tree, err := files.LoadToml(networkTomlFile)
	if err != nil {
		common.ErrorMessage("Cannot load network.toml file: ", networkTomlFile)
		return err
	}

	nodesquery, err := query.CompileAndExecute("$.nodes", tree)
	if err != nil {
		common.ErrorMessage("Error loading nodes")
		return err
	}

	var dumpOut []string

	for _, value := range nodesquery.Values() {
		if reflect.TypeOf(value).String() == "[]*toml.Tree" {
			nodes := value.([]*toml.Tree)

			for _, tr := range nodes { // loop around nodes
				// Data wrangling
				var addr, moniker, netaddr string
				var validator bool

				if tr.HasPath([]string{"moniker"}) {
					moniker = tr.GetPath([]string{"moniker"}).(string)
				}
				if tr.HasPath([]string{"netaddr"}) {
					netaddr = tr.GetPath([]string{"netaddr"}).(string)
					if idx := strings.Index(netaddr, ":"); idx > -1 {
						netaddr = netaddr[:idx]
					}
				}
				if tr.HasPath([]string{"address"}) {
					addr = tr.GetPath([]string{"address"}).(string)
				}

				if tr.HasPath([]string{"validator"}) {
					validator = tr.GetPath([]string{"validator"}).(bool)
				}

				dumpOut = append(dumpOut, moniker+"|"+netaddr+"|"+addr+"|"+strconv.FormatBool(validator))

			}
		}

	}

	for _, o := range dumpOut {
		fmt.Println(o)
	}

	return nil
}
