package commands

import (
	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/keys"
	"github.com/mosaicnetworks/monetd/cmd/monetd/config"
	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RootCmd is the root command for monetcli
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
		//		network.NetworkCmd,
		//		config.ConfigCmd,
		VersionCmd,
		//		network.WizardCmd,
		//		testnet.NewTestNetCmd(),
		//		testnet.NewTestJoinCmd(),
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true

	RootCmd.PersistentFlags().StringVarP(&config.Config.DataDir, "datadir", "d", config.Config.DataDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().BoolVarP(&common.VerboseLogging, "verbose", "v", false, "verbose messages")
	viper.BindPFlags(RootCmd.Flags())
}