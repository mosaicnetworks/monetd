package commands

import (
	"github.com/spf13/cobra"
)

//RootCmd is the root command for monetd
var RootCmd = &cobra.Command{
	Use:   "monetd",
	Short: "Monet-Daemon",
}

func init() {
	RootCmd.AddCommand(
		RunCmd,
		VersionCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true
}

