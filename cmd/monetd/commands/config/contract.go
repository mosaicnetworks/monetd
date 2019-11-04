package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/mosaicnetworks/monetd/src/genesis"

	"github.com/spf13/cobra"
)

// newContractCommand returns the ContractCmd
func newContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [validator] [validator] ...",
		Short: "display poa contract",
		Long: `
Display the PoA smart contract.

Outputs the standard monetd contract, configured with the list of validators
in the initial whitelist. The validator arguments are either a moniker that
is in the keystore, or in the form moniker=address. The two forms can be mixed
in the same command line.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: contractConfig,
	}
	return cmd
}

func contractConfig(cmd *cobra.Command, args []string) error {
	var err error

	var minPeers []*genesis.MinimalPeerRecord

	for _, peer := range args {
		splitRec := strings.Split(peer, "=")
		if len(splitRec) > 1 {
			minPeers = append(minPeers, &genesis.MinimalPeerRecord{Address: splitRec[1], Moniker: splitRec[0]})
		} else {
			jsonfile := filepath.Join(configuration.DefaultKeystoreDir(), peer+".json")

			// Read key from file.
			keyjson, err := ioutil.ReadFile(jsonfile)
			if err != nil {
				return fmt.Errorf("Failed to read the keyfile at '%s': %v", jsonfile, err)
			}

			k := new(crypto.EncryptedKeyJSONMonet)
			if err := json.Unmarshal(keyjson, k); err != nil {
				return err
			}
			minPeers = append(minPeers, &genesis.MinimalPeerRecord{Address: k.Address, Moniker: peer})
		}
	}

	solSource, err := genesis.GetFinalSoliditySourceFromAddress(minPeers)

	fmt.Print(solSource)
	fmt.Println("")
	return err
}
