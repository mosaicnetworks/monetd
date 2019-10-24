package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

var (
	_keystore     = configuration.DefaultKeystoreDir()
	_configDir    = configuration.DefaultConfigDir()
	_keyParam     = getDefaultKey() //get default keyfile
	_addressParam = common.GetMyIP()
	_passwordFile string
	_force        = false
)

// ConfigCmd implements the config CLI subcommand
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "manage configuration",
	Long: `
Manage monetd configuration.

* config build - creates the configuration for a single-node network, based on 
                 one of the keys in <keystore>. This is a quick and easy way to
                 get started with monetd. 

* config pull -  fetches the configuration from a running node. This is used to
                 join an existing network.

For more complex scenarios, please refer to 'giverny', which is a specialised 
Monet configuration tool. 
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

// getDefaultKey returns the moniker of the the first keyfile in the default
// keystore
func getDefaultKey() string {

	keystore := configuration.DefaultKeystoreDir()

	files, err := ioutil.ReadDir(keystore)
	if err != nil {
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
