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
	Long: `Giverny Version information
	
The version command outputs the version number for Monet, EVM-Lite, 
Babble and Geth. 

If you compile your own tools, the suffices are the GIT branch and the GIT
commit hash. 
`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print(monet.FullVersion())
	},
}
