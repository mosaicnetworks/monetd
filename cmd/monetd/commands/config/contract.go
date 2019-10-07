package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"
	"github.com/mosaicnetworks/monetd/src/crypto"

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

	jsonfile := filepath.Join(configuration.Global.DataDir, configuration.KeyStoreDir, node+".json")
	// For a simple change, tree is quicker and easier than unmarshalling the whole tree

	// Read key from file.
	keyjson, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", jsonfile, err)
	}

	k := new(crypto.EncryptedKeyJSONMonet)
	if err := json.Unmarshal(keyjson, k); err != nil {
		return err
	}

	peersJSON := types.PeerRecordList{}
	peersJSON = append(peersJSON, &types.PeerRecord{NetAddr: "localhost:1337", PubKeyHex: k.PublicKey, Moniker: node})

	solSource, err := contract.GetFinalSoliditySource(peersJSON)

	fmt.Print(solSource)
	fmt.Println("")
	return err
}
