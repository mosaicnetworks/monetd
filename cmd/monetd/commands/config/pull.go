package config

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	pconfig "github.com/mosaicnetworks/monetd/src/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [host:port]",
		Short: "pull the configuration files from a node",
		Long: `
The pull subcommand is used to join an existing Monet network. It takes the
address (host:port) of a running node, and downloads the following set of files 
into the configuration directory [datadir]:

- babble/peers.json         : The current validator-set 
- babble/peers.genesis.json : The initial validator-set
- eth/genesis.json          : The genesis file

Additionally, this command configures the validator-key and the network address
of the new node. The --key and --passfile options refer to the validator-key, 
while --address sets the network address of monetd.
`,
		Example: `  monetd config pull "192.168.5.1:8080"`,
		Args:    cobra.ExactArgs(1),
		RunE:    pullConfig,
	}
	addPullFlags(cmd)

	return cmd
}

func addPullFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&keyParam, "key", keyParam, "moniker of the key to use for this node")
	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func pullConfig(cmd *cobra.Command, args []string) error {
	peerAddr := args[0]

	common.InfoMessage(fmt.Sprintf("Pulling configuration from %s for key %s on %s",
		peerAddr, keyParam, addressParam))

	err := pconfig.PullConfig(configuration.Global.DataDir, keyParam, addressParam, peerAddr, passwordFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
