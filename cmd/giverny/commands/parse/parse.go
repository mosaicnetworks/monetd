package parse

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/version"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mosaicnetworks/monetd/src/genesis"

	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/spf13/cobra"
)

//ParseCmd is an Ethereum key manager
var ParseCmd = &cobra.Command{
	Use:   "parse [genesis file]",
	Short: "parse genesis file",
	Long: `
The parse command parses a genesis file. 
`,
	Args: cobra.ExactArgs(1),
	RunE: parseGenesis,
}

func parseGenesis(cmd *cobra.Command, args []string) error {
	genesisFile := args[0]

	// Check the file exists
	if !files.CheckIfExists(genesisFile) {
		return errors.New("cannot find the file " + genesisFile)
	}

	// Read Genesis file and load into genesisJSON struct
	genesisJSON := genesis.JSONGenesisFile{}

	file, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		fmt.Println("Error loading " + genesisFile)
		return err
	}

	err = json.Unmarshal([]byte(file), &genesisJSON)
	if err != nil {
		fmt.Println("Error parsing " + genesisFile)
		return err
	}

	fmt.Println("")

	fmt.Printf("POA Address:  0x%s \n", genesisJSON.Poa.Address)
	fmt.Println("")
	peers := make(map[uint64]genesis.MinimalPeerRecord)

	// Populate a mapping peers, from the solidity array

	for key, val := range genesisJSON.Poa.Storage {
		// The hardcoded values corresponds to the solidity array.
		// Easier to pick up than the mapping.
		if key[0:59] == "405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3b" {
			// Parse last 5 hex digits, and rebase to zero
			// For all practical purposes 10,000 whitelist entries is sufficient
			intkey, err := strconv.ParseUint(key[59:], 16, 64)
			intkey -= 744141 // Rebase intkey to 1 for convenience
			if err != nil {
				fmt.Println("Error parsing hex " + key[60:])
				return err
			}

			// Actually populate the peer array
			peer, ok := peers[intkey]
			if !ok {
				peers[intkey] = genesis.MinimalPeerRecord{
					Address: fmt.Sprintf("0x%s", strings.TrimPrefix(val, "94")),
					Moniker: fmt.Sprintf("Peer %03d", intkey),
				}
			} else {
				peer.Address = fmt.Sprintf("0x%s", strings.TrimPrefix(val, "94"))
				peers[intkey] = peer
			}
		}
	}

	fmt.Printf("%d peers found \n\n", len(peers))

	// Slot 5 is the moniker mapping.
	SLOT5 := fmt.Sprintf("%064d", 5)
	slot5Bytes := common.HexToHash(SLOT5).Bytes()

	for _, peer := range peers {
		addr := strings.TrimPrefix(strings.ToLower(peer.Address), "0x")
		addrBytes := common.HexToHash(addr).Bytes()
		// Handle the Moniker mapping
		addrHash := eth_crypto.Keccak256(append(addrBytes, slot5Bytes...))
		addrSlot := hex.EncodeToString(addrHash)

		moniker, ok := genesisJSON.Poa.Storage[addrSlot]
		if ok {
			monikerBytes, err := hex.DecodeString(moniker)
			if err == nil {
				moniker = string(monikerBytes)
			} else {
				moniker = peer.Moniker
			}
		} else {
			moniker = peer.Moniker
		}

		fmt.Printf("%s  %s\n", peer.Address, moniker)

	}

	fmt.Println("")

	if genesisJSON.Poa.Code == genesis.StandardPOAContractByteCode {
		fmt.Println("POA bytecode matches the standard contract")
	} else {
		fmt.Println("Contract does not match the POA bytecode")
		fmt.Println("This may not be an issue if a different release of Monetd " +
			"was used to generate the genesis.json file.")
		fmt.Println("Your version of Monetd is:")
		fmt.Print(version.FullVersion())
		fmt.Printf("Solc: %s \n      %s\n", genesis.SolcCompilerVersion, genesis.SolcOSVersion)
		fmt.Printf("      %s\n", genesis.GitVersion)

	}

	return nil
}
