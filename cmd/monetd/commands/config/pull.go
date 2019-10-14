package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type downloadItem struct {
	URL  string
	Dest string
}

func newPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [host:port]",
		Short: "pull the configuration files from a node",
		Long: `
The pull subcommand is used to join an existing network. It takes the address
(host:port) of a running node, and downloads the following set of files into the
configuration directory <config>:

- babble/peers.json         : The current validator-set 
- babble/peers.genesis.json : The initial validator-set
- eth/genesis.json          : The genesis file

Additionally, this command configures the key and network address of the new
node. The --key flag identifies a keyfile by moniker, which is expected to be in 
the <keystore>. If --passfile is not specified, the user will be prompted to
enter the passphrase manually. If the --address flag is omitted, the first 
non-loopback address is used.
`,
		Example: `  monetd config pull "192.168.5.1:8080"`,
		Args:    cobra.ExactArgs(1),
		RunE:    pullConfig,
	}
	addPullFlags(cmd)

	return cmd
}

func addPullFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&_keystore, "keystore", _keystore, "keystore directory")
	cmd.Flags().StringVar(&_configDir, "config", _configDir, "output directory")
	cmd.Flags().StringVar(&_addressParam, "address", _addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&_keyParam, "key", _keyParam, "moniker of the key to use for this node")
	cmd.Flags().StringVar(&_passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func pullConfig(cmd *cobra.Command, args []string) error {
	peerAddr := args[0]

	// Helpful debugging output
	common.MessageWithType(common.MsgDebug, "Pulling from         : ", peerAddr)
	common.MessageWithType(common.MsgDebug, "Using Network Address: ", _addressParam)
	common.MessageWithType(common.MsgDebug, "Using Key            : ", _keyParam)
	common.MessageWithType(common.MsgDebug, "Using Password File  : ", _passwordFile)

	// Retrieve the keyfile corresponding to moniker
	privateKey, err := keystore.GetKey(_keystore, _keyParam, _passwordFile)
	if err != nil {
		return err
	}

	// Create Directories if they don't exist
	files.CreateMonetConfigFolders(_configDir)

	// Copy the key to babble directory with appropriate permissions
	err = keystore.DumpPrivKey(
		filepath.Join(_configDir, configuration.BabbleDir),
		privateKey)
	if err != nil {
		return err
	}

	rootURL := "http://" + peerAddr

	filesList := []*downloadItem{
		&downloadItem{URL: rootURL + "/genesispeers",
			Dest: filepath.Join(_configDir, configuration.BabbleDir, configuration.PeersGenesisJSON)},
		&downloadItem{URL: rootURL + "/peers",
			Dest: filepath.Join(_configDir, configuration.BabbleDir, configuration.PeersJSON)},
		&downloadItem{URL: rootURL + "/genesis",
			Dest: filepath.Join(_configDir, configuration.EthDir, configuration.GenesisJSON)},
	}

	for _, item := range filesList {
		err := files.DownloadFile(item.URL, item.Dest, files.OverwriteSilently)
		if err != nil {
			common.ErrorMessage(fmt.Sprintf("Error downloading %s", item.URL))
			return err
		}
		common.DebugMessage("Downloaded ", item.Dest)
	}

	// Write TOML file for monetd based on global config object
	return configuration.DumpGlobalTOML(_configDir, configuration.MonetTomlFile)
}
