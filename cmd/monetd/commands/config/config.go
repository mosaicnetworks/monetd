package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/network"
	"github.com/spf13/cobra"
)

var (
	keyParam     = getDefaultKey() //get default keyfile
	addressParam = network.GetMyIP()
	passwordFile string
)

// ConfigCmd implements the config CLI subcommand
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage monetd configuration",
	Long: `
The config subcommand initialises the configuration for a monetd node in the
folder specified by [datadir] (~/.monet by default on Linux). The configuration
creates all the files necessary for a node to join an existing network or 
to create a new one.

There are two ways of initialising the configuration:

* config build - config build creates the configuration for a single-node 
                 network, based on one of the keys in [datadir]/keystore. 
                 This is a quick and easy way to get started with monetd. 

* config pull - config pull is used to join an existing network. It fetches the 
                configuration from one of the existing nodes.

For more complex scenarios, please refer to 'giverny', which is a specialised 
monet configuration tool. 
`,
	TraverseChildren: true,
}

func init() {
	// Subcommands
	ConfigCmd.AddCommand(
		newLocationCmd(),
		newClearCmd(),
		newPullCmd(),
		newBuildCmd(),
		newContractCmd(),
	)
}

// getDefaultKey returns the moniker of the the first keyfile in
// [datadir]/keystore
func getDefaultKey() string {

	keystore := filepath.Join(configuration.Configuration.DataDir, "keystore")

	files, err := ioutil.ReadDir(keystore)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			return strings.TrimSuffix(
				file.Name(),
				filepath.Ext(file.Name()),
			)
		}
	}

	return ""
}
