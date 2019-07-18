package config

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

func pullConfig(cmd *cobra.Command, args []string) error {
	common.MessageWithType(common.MsgInformation, "The Monet Configuration files are located at:")
	common.MessageWithType(common.MsgInformation, monetConfigDir)
	return nil
}
