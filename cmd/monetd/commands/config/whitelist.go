package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/contract"
	"github.com/mosaicnetworks/monetd/src/crypto"

	"github.com/spf13/cobra"
)

// newContractCommand returns the ContractCmd
func newWhiteListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist network [validator] [validator] ...",
		Short: "generate the POA Storage value",
		Long: `
Generate the POA Storage value.

Generates the storage values to be put into the Storage section of the Genesis
file.`,
		Args: cobra.MinimumNArgs(2),
		RunE: whiteListConfig,
	}
	return cmd
}

func whiteListConfig(cmd *cobra.Command, args []string) error {
	var err error

	var minPeers []*contract.MinimalPeerRecord

	//	var network = args[0]

	for _, peer := range args[1:] {
		splitRec := strings.Split(peer, "=")
		if len(splitRec) > 1 {
			minPeers = append(minPeers, &contract.MinimalPeerRecord{Address: splitRec[1], Moniker: splitRec[0]})
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
			minPeers = append(minPeers, &contract.MinimalPeerRecord{Address: k.Address, Moniker: peer})
		}
	}

	storage, err := contract.GetStorage(minPeers)

	js, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	fmt.Print(string(js))
	//	fmt.Println("")
	return nil
}
