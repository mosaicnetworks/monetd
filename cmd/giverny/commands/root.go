package commands

import (
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
		//		keys.KeysCmd,
		//		network.NetworkCmd,
		//		config.ConfigCmd,
		VersionCmd,
		//		network.WizardCmd,
		//		testnet.NewTestNetCmd(),
		//		testnet.NewTestJoinCmd(),
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true

	viper.BindPFlags(RootCmd.Flags())
}
