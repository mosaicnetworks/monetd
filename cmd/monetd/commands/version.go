package commands

import (
	"fmt"

	monet "github.com/mosaicnetworks/monetd/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of evml being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Long:  `Monetd Version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(monet.FullVersion())
	},
}
