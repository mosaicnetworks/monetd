package config

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mosaicnetworks/monetd/src/configuration"
	pconfig "github.com/mosaicnetworks/monetd/src/config"
)

// newBuildCmd initialises a bare-bones configuration for monetd
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [moniker]",
		Short: "create the configuration for a single-node network",
		Long: `
The build subcommand initialises the bare-bones configuration to get started 
with monetd. It uses one of the accounts from the keystore to define a network 
consisting of a unique node, which is automatically added to the PoA whitelist.
Additionally, all the accounts in [datadir]/keystore are credited with a large
amount of tokens in the genesis file. This command is mostly used for testing.

If the --address flag is omitted, the first non-loopback address for this 
instance is used.
`,
		Args: cobra.ExactArgs(1),
		RunE: buildConfig,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func buildConfig(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	common.InfoMessage(fmt.Sprintf("Builing configuration for key %s on %s", moniker, addressParam))

	err := pconfig.BuildConfig(configuration.Global.DataDir, moniker, addressParam, passwordFile)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
