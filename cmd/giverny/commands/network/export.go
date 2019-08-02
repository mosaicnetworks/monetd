package network

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
)

//set by command line flags
// var includePassPhrase = false

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export [network] [node]",
		Short: "export the configuration for a node on the named network",
		Long: `
giverny network export
		`,
		Args: cobra.ArbitraryArgs,
		RunE: networkExport,
	}

	return cmd
}

func networkExport(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return errors.New("you must specify a network name")
	}

	networkName := args[0]
	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	if !files.CheckIfExists(networkDir) {
		return errors.New("cannot read the configuration for network " + networkName)
	}

	if len(args) == 1 { // only specified a network name - so iterate through all nodes
		keystore := filepath.Join(networkDir, givernyKeystoreDir)
		files, err := ioutil.ReadDir(keystore)
		if err != nil {
			return err
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".json" {
				nodeName := strings.TrimSuffix(
					file.Name(),
					filepath.Ext(file.Name()))
				err := buildZip(configuration.GivernyConfigDir, networkName, nodeName)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for i := 1; i < len(args); i++ {
		nodeName := args[i]
		err := buildZip(configuration.GivernyConfigDir, networkName, nodeName)
		if err != nil {
			return err
		}
	}

	return nil
}
