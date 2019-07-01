package commands

import (
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/config"
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/network"
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/testnet"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RootCmd is the root command for monetcli
var RootCmd = &cobra.Command{
	Use:   "monetcli",
	Short: "Monet-CLI",
}

func init() {

	RootCmd.AddCommand(
		keys.KeysCmd,
		network.NetworkCmd,
		config.ConfigCmd,
		VersionCmd,
		network.WizardCmd,
		testnet.NewTestNetCmd(),
		testnet.NewTestJoinCmd(),
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true

	RootCmd.PersistentFlags().BoolVarP(&common.VerboseLogging, "verbose", "v", false, "verbose messages")
	viper.BindPFlags(RootCmd.Flags())
}
