package network

import (
	"errors"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/common"
	com "github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
)

func generatekeypair(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	ip := args[1]
	isValidator, _ := strconv.ParseBool(args[2])

	safeLabel := com.GetNodeSafeLabel(moniker)
	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

	for _, node := range currentNodes {
		if node == safeLabel {
			common.Message("That Moniker has already been used", moniker)
			return errors.New("that moniker has already been used")
		}
	}

	return GenerateKeyPair(configDir, moniker, ip, isValidator, passwordFile)
}
