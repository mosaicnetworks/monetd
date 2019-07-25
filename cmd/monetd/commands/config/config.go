package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

var (
	keyParam     = getDefaultKey() //get default keyfile
	addressParam = common.GetMyIP()
	passwordFile string
)

// ConfigCmd implements the config CLI subcommand
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage monetd configuration",
	Long: `
The monetd server reads configuration from --datadir.

There are two ways of initialising the configuration:

* config build - creates the configuration for a single-node network, based on 
                 one of the keys in [datadir]/keystore. This is a quick and easy 
                 way to get started with monetd. 

* config pull -  fetches the configuration from a running node. This is used to
                 join an existing network.

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

	keystore := filepath.Join(configuration.Global.DataDir, "keystore")

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
