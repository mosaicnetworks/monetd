package commands

import (
	"github.com/mosaicnetworks/evm-lite/cmd/evml/commands/keys"
	"github.com/spf13/cobra"
)

//RootCmd is the root command for evml
var RootCmd = &cobra.Command{
	Use:   "monetcli",
	Short: "Monet-CLI",
}

func init() {
	RootCmd.AddCommand(
		keys.KeysCmd,
		VersionCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true
}
