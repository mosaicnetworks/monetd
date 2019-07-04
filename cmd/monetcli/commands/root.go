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
	Long: `Monet-CLI
	
Monetcli is the swiss army knife of tools for the Monet Hub. The README.md in the github repository is a good starting place in the documentation. For quicker access and a handy reference of flags and options.:
	
monetcli help [subcommand]
	
The best starting points are the commands testnet, testjoin or wizard. `,
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
	RootCmd.PersistentFlags().BoolVarP(&common.HideBanners, "hide-banners", "q", false, "hide banners")

	viper.BindPFlags(RootCmd.Flags())
}
