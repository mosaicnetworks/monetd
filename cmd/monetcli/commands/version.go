package commands

import (
	"fmt"

	geth "github.com/ethereum/go-ethereum/params"
	_babble "github.com/mosaicnetworks/babble/src/version"
	evm "github.com/mosaicnetworks/evm-lite/src/version"
	monet "github.com/mosaicnetworks/monetd/src/version"
	"github.com/spf13/cobra"
)

// VersionCmd displays the version of evml being used
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Long:  `Monet-CLI Version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Monet Version: " + monet.Version)
		fmt.Println("     EVM-Lite Version: " + evm.Version)
		fmt.Println("     Babble Version: " + _babble.Version)
		fmt.Println("     Geth Version: " + geth.Version)
	},
}
