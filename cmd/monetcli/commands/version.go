package commands

import (
	"fmt"

	monet "github.com/mosaicnetworks/monetd/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of evml being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Long:  `Monet-CLI Version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(monet.FullVersion())
	},
}
