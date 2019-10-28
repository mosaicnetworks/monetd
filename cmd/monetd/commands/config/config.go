package config

import (
	"errors"
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
	_keyParam     = ""
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
func getDefaultKey(keystore string) (string, error) {

	files, err := ioutil.ReadDir(keystore)
	if err != nil {
		return "", err
	}

	var monikers []string

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			monikers = append(monikers, strings.TrimSuffix(
				file.Name(),
				filepath.Ext(file.Name()),
			))
		}
	}

	if len(monikers) == 0 {
		return "", errors.New("No keys found. Use 'monet keys new' to generate keys ")
	}

	if len(monikers) == 1 {
		return monikers[0], nil
	}

	common.ErrorMessage("You have multiple available keys. Specify one using the --key parameter.")
	common.InfoMessage(monikers)

	return "", errors.New("key to use is ambiguous")

}
