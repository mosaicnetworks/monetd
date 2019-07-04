package commands

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/common"
	monet "github.com/mosaicnetworks/monetd/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of evml being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Long: `Monet-CLI Version information
	
	This command returns the version number to the monetcli app itself, 
	and the version of the EVM-Lite, Babble and Geth librarys used to 
	build it. The suffix (if shown) on the Monet version if the github 
	commit for this version.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		common.Banner("monetd")
		common.BlankLine()
		fmt.Print(monet.FullVersion())
	},
}
