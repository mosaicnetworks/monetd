package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"
	"github.com/mosaicnetworks/monetd/src/crypto"

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

	jsonfile := filepath.Join(configuration.DefaultKeystoreDir(), node+".json")

	// Read key from file.
	keyjson, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", jsonfile, err)
	}

	k := new(crypto.EncryptedKeyJSONMonet)
	if err := json.Unmarshal(keyjson, k); err != nil {
		return err
	}

	whitelistPeers := []*peers.Peer{
		peers.NewPeer(k.PublicKey, "localhost:1337", node),
	}

	solSource, err := contract.GetFinalSoliditySource(whitelistPeers)

	fmt.Print(solSource)
	fmt.Println("")
	return err
}
