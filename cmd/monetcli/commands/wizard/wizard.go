package wizard

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Wizard to set up a Monet Network",
	Long:  `Wizard to set up a Monet Network`,
	Run: func(cmd *cobra.Command, args []string) {

		home, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
		requestFile("Configuration Directory Location: ", home)

	},
}
