package commands

import (
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/network"
	"github.com/spf13/cobra"
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
		VersionCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true
}
