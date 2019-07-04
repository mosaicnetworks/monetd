package network

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

func setParams(cmd *cobra.Command, args []string) error {
	// Call a wrapper function to ease calling from outside cobra
	common.BannerTitle("Params")
	return common.SetParamsWithParams(configDir)
}
