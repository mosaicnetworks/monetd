package network

import (
	"strconv"

	"github.com/spf13/cobra"
)

func generatekeypair(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	ip := args[1]
	isValidator, _ := strconv.ParseBool(args[2])

	return GenerateKeyPair(configDir, moniker, ip, isValidator)
}
