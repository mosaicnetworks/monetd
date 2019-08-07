package config

import (
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/types"
	"github.com/spf13/cobra"
)

// newContractCommand returns the ContractCmd
func newContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [moniker]",
		Short: "display poa contract",
		Long: `
Display the PoA smart contract.

Outputs the standard monetd contract, configured with [moniker] in the initial
whitelist.
`,
		Args: cobra.ExactArgs(1),
		RunE: contractConfig,
	}
	return cmd
}

func contractConfig(cmd *cobra.Command, args []string) error {

	node := args[0]

	tomlfile := filepath.Join(configuration.Global.DataDir, configuration.KeyStoreDir, node+".toml")
	// For a simple change, tree is quicker and easier than unmarshalling the whole tree
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
