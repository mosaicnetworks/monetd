package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/contract"

	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"
	"github.com/mosaicnetworks/monetd/src/poa/types"
	"github.com/spf13/cobra"
)

//newClearCmd shows the config file path
func newContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [moniker]",
		Short: "displays poa contract",
		Long: `
monetd config contract

Outputs the standard monetd contract, configured with [moniker] as the initial
whitelist.
`,
		Args: cobra.ExactArgs(1),
		RunE: contractConfig,
	}
	return cmd
}

func contractConfig(cmd *cobra.Command, args []string) error {

	node := args[0]
	safeLabel := common.GetNodeSafeLabel(node)

	tomlfile := filepath.Join(configuration.Global.DataDir, common.KeyStoreDir, safeLabel+".toml")
	tree, err := files.LoadToml(tomlfile)
	if err != nil {
		common.MessageWithType(common.MsgError, "Cannot read peer configuration: ", tomlfile)
		return err
	}
	pubkey := tree.GetPath([]string{"node", "pubkey"}).(string)
	moniker := tree.GetPath([]string{"node", "moniker"}).(string)
	peersJSON := types.PeerRecordList{}
	peersJSON = append(peersJSON, &types.PeerRecord{NetAddr: "localhost:1337", PubKeyHex: pubkey, Moniker: moniker})

	solSource, err := contract.GetFinalSoliditySource(peersJSON)

	fmt.Print(solSource)
	fmt.Println("")
	return err
}
