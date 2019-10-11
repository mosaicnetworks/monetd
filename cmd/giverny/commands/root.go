package commands

import (
	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/keys"
	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/network"
	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/server"
	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/transactions"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RootCmd is the root command for giverny
var RootCmd = &cobra.Command{
	Use:   "giverny",
	Short: "Giverny",
	Long: `Giverny
	
Giverny is the swiss army knife of advanced tools for the Monet Hub. For most users, 
you should not need to use this command. The inbuild commands in monetd will suffice for
most use cases.`,
}

func init() {

	RootCmd.AddCommand(
		keys.KeysCmd,
		server.ServerCmd,

		network.NetworkCmd,
		transactions.TransCmd,
		//		config.ConfigCmd,
		VersionCmd,
		//		network.WizardCmd,
		//		testnet.NewTestNetCmd(),
		//		testnet.NewTestJoinCmd(),
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true

	RootCmd.PersistentFlags().StringVarP(&monetconfig.Global.DataDir, "monet-data-dir", "m", monetconfig.Global.DataDir, "Top-level monetd directory for configuration and data")
	RootCmd.PersistentFlags().StringVarP(&configuration.GivernyConfigDir, "giverny-data-dir", "g", configuration.GivernyConfigDir, "Top-level giverny directory for configuration and data")
	RootCmd.PersistentFlags().BoolVar(&monetconfig.NonInteractive, "non-interactive", false, "non-interactive")
	RootCmd.PersistentFlags().BoolVarP(&common.VerboseLogging, "verbose", "v", false, "verbose messages")
	viper.BindPFlags(RootCmd.Flags())
}
