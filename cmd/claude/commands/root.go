package commands

import (
	"github.com/mosaicnetworks/monetd/cmd/claude/commands/trans"
	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RootCmd is the root command for giverny
var RootCmd = &cobra.Command{
	Use:   "claude",
	Short: "Claude",
	Long: `Claude
	
Claude is a staging area for new potential giverny commands.
`,
}

func init() {

	RootCmd.AddCommand(
		VersionCmd,
		trans.TransCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true

	RootCmd.PersistentFlags().BoolVarP(&common.VerboseLogging, "verbose", "v", false, "verbose messages")
	viper.BindPFlags(RootCmd.Flags())
}
