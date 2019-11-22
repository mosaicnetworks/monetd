package parse

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

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

type storageRow struct {
	rawvalue    string
	value       string
	description string
	explained   bool
}

func parseGenesis(cmd *cobra.Command, args []string) error {
	genesisFile := args[0]

	// Load Genesis Data

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

	rows := make(map[string]*storageRow)
	addressMap := make(map[string]string)

	// Build our output structure
	for key1, val1 := range genesisJSON.Poa.Storage {
		rows[key1] = &storageRow{
			rawvalue:    val1,
			value:       "",
			explained:   false,
			description: "Unexplained",
		}
	}

	//  Slot 1 - whitelistCount
	slot := fmt.Sprintf("%064d", 1)
	whiteListCount := 0
	if rows[slot] == nil {
		rows[slot] = &storageRow{
			value:       "0",
			rawvalue:    "0",
			explained:   true,
			description: "whitelistCount is zero so no row set",
		}
	} else {
		rows[slot].value = rows[slot].rawvalue
		rows[slot].explained = true
		rows[slot].description = "whitelistCount is set"
		whiteListCount, _ = strconv.Atoi(rows[slot].value)
	}

	// Slot 2 - whiteListArray
	slot = fmt.Sprintf("%064d", 2)
	whiteListArrayLength := 0
	if rows[slot] == nil {
		rows[slot] = &storageRow{
			value:       "0",
			rawvalue:    "0",
			explained:   true,
			description: "whiteListArray length is zero so no row set",
		}
	} else {
		rows[slot].value = rows[slot].rawvalue
		rows[slot].explained = true
		rows[slot].description = "whiteListArray Length is set"
		whiteListArrayLength, _ = strconv.Atoi(rows[slot].value)
	}

	if whiteListArrayLength != whiteListCount {
		return fmt.Errorf("WhitelistCount (%d) and length of whiteListArray (%d) differ", whiteListCount, whiteListArrayLength)
	}

	slotBytes := eth_crypto.Keccak256(common.HexToHash(slot).Bytes())
	slotHex := strings.TrimPrefix(hexutil.Encode(slotBytes), "0x")
	fmt.Printf("Slot 2 Whitelist Array base address is at: %s \n", slotHex)

	baseSlot := new(big.Int)
	baseSlot, _ = baseSlot.SetString(slotHex, 16)

	for i := 0; i < whiteListCount; i++ {

		j := new(big.Int)
		j.SetInt64(int64(i))
		newSlot := j.Add(baseSlot, j)
		newHex := newSlot.Text(16)

		if rows[newHex] == nil {
			rows[newHex] = &storageRow{
				value:       "0",
				rawvalue:    "0",
				explained:   true,
				description: fmt.Sprintf("whiteListArray slot %d not set", i),
			}
		} else {
			rows[newHex].value = parseAddress(rows[newHex].rawvalue)
			rows[newHex].explained = true
			monikerSlot := getMonikerSlot(rows[newHex].value)

			addressMap[rows[newHex].value] = fmt.Sprintf("unknown peer %d", i)
			if rows[monikerSlot] != nil {
				moniker := rows[monikerSlot].rawvalue
				monikerBytes, err := hex.DecodeString(moniker[2:])
				if err == nil {
					addressMap[rows[newHex].value] = string(monikerBytes)
					if rows[monikerSlot].description != "" {
						rows[monikerSlot].value = string(monikerBytes)
						rows[monikerSlot].explained = true
						rows[monikerSlot].description =
							fmt.Sprintf("Moniker for %s", rows[newHex].value)

					}
				}
			}
			rows[newHex].description = fmt.Sprintf("whiteListArray slot %d set for %s", i, addressMap[rows[monikerSlot].value])

		}
		//		fmt.Printf("    %d %s \n", i, newHex)
	}

	// 		4	3	Slot 4 nomineeArray length – set t o 3

	// Slot 4 - nomineeListArray
	slot = fmt.Sprintf("%064d", 4)
	nomineeListCount := 0

	if rows[slot] == nil {
		rows[slot] = &storageRow{
			value:       "0",
			rawvalue:    "0",
			explained:   true,
			description: "Nominee Array length is zero so no row set",
		}
	} else {
		rows[slot].value = rows[slot].rawvalue
		rows[slot].explained = true
		rows[slot].description = "Nominee Array Length is set"
		nomineeListCount, _ = strconv.Atoi(rows[slot].value)
	}

	slotBytes = eth_crypto.Keccak256(common.HexToHash(slot).Bytes())
	slotHex = strings.TrimPrefix(hexutil.Encode(slotBytes), "0x")
	fmt.Printf("Slot 4 Nominee List Array base address is at: %s \n", slotHex)

	baseSlot = new(big.Int)
	baseSlot, _ = baseSlot.SetString(slotHex, 16)

	for i := 0; i < nomineeListCount; i++ {

		j := new(big.Int)
		j.SetInt64(int64(i))
		newSlot := j.Add(baseSlot, j)
		newHex := newSlot.Text(16)

		if rows[newHex] == nil {
			rows[newHex] = &storageRow{
				value:       "0",
				rawvalue:    "0",
				explained:   true,
				description: fmt.Sprintf("Nominee Array slot %d not set", i),
			}
		} else {
			rows[newHex].value = parseAddress(rows[newHex].rawvalue)
			rows[newHex].explained = true

			monikerSlot := getMonikerSlot(rows[newHex].value)

			addressMap[rows[newHex].value] = fmt.Sprintf("unknown peer %d", i)
			if rows[monikerSlot] != nil {
				moniker := rows[monikerSlot].rawvalue
				monikerBytes, err := hex.DecodeString(moniker[2:])
				if err == nil {
					addressMap[rows[newHex].value] = string(monikerBytes)
					if rows[monikerSlot].description != "" {
						rows[monikerSlot].value = string(monikerBytes)
						rows[monikerSlot].explained = true
						rows[monikerSlot].description =
							fmt.Sprintf("Moniker for %s", rows[newHex].value)

					}
				}
			}
			rows[newHex].description = fmt.Sprintf("Nominee Array slot %d set for %s", i, addressMap[rows[newHex].value])

		}
		//		fmt.Printf("    %d %s \n", i, newHex)
	}

	//	7	0	Slot 7, EvictionArray size, zero so empty

	// Slot 7 - Eviction ListArray
	slot = fmt.Sprintf("%064d", 7)
	evicteeListCount := 0

	if rows[slot] == nil {
		rows[slot] = &storageRow{
			value:       "0",
			rawvalue:    "0",
			explained:   true,
			description: "Eviction Array length is zero so no row set",
		}
	} else {
		rows[slot].value = rows[slot].rawvalue
		rows[slot].explained = true
		rows[slot].description = "Eviction Array Length is set"
		evicteeListCount, _ = strconv.Atoi(rows[slot].value)
	}

	slotBytes = eth_crypto.Keccak256(common.HexToHash(slot).Bytes())
	slotHex = strings.TrimPrefix(hexutil.Encode(slotBytes), "0x")
	fmt.Printf("Slot 7 Eviction List Array base address is at: %s \n", slotHex)

	baseSlot = new(big.Int)
	baseSlot, _ = baseSlot.SetString(slotHex, 16)

	for i := 0; i < evicteeListCount; i++ {

		j := new(big.Int)
		j.SetInt64(int64(i))
		newSlot := j.Add(baseSlot, j)
		newHex := newSlot.Text(16)

		if rows[newHex] == nil {
			rows[newHex] = &storageRow{
				value:       "0",
				rawvalue:    "0",
				explained:   true,
				description: fmt.Sprintf("Eviction Array slot %d not set", i),
			}
		} else {
			rows[newHex].value = parseAddress(rows[newHex].rawvalue)
			rows[newHex].explained = true

			monikerSlot := getMonikerSlot(rows[newHex].value)
			addressMap[rows[newHex].value] = fmt.Sprintf("unknown peer %d", i)
			if rows[monikerSlot] != nil {

				moniker := rows[monikerSlot].rawvalue
				monikerBytes, err := hex.DecodeString(moniker[2:])
				if err == nil {
					addressMap[rows[newHex].value] = string(monikerBytes)
					if rows[monikerSlot].description != "" {
						rows[monikerSlot].value = string(monikerBytes)
						rows[monikerSlot].explained = true
						rows[monikerSlot].description =
							fmt.Sprintf("Moniker for %s", rows[newHex].value)

					}
				}
			}
			rows[newHex].description = fmt.Sprintf("Eviction Array slot %d set for %s", i, addressMap[rows[newHex].value])

		}
		//		fmt.Printf("    %d %s \n", i, newHex)
	}

	//	5	0	Slot 5 monikerList mapping so empty
	//  Has been set as part of the logic above

	// For pragmatic reasons and because of potentially stale data we check the mappings for addresses that we know about
	// TODO add a command line option to be able to add extra addresses so we can check for any additional addresses that
	// the user knows about.

	SLOT0 := fmt.Sprintf("%064d", 0)
	slot0Bytes := common.HexToHash(SLOT0).Bytes()
	SLOT3 := fmt.Sprintf("%064d", 3)
	slot3Bytes := common.HexToHash(SLOT3).Bytes()
	SLOT6 := fmt.Sprintf("%064d", 6)
	slot6Bytes := common.HexToHash(SLOT6).Bytes()

	for key, moniker := range addressMap {
		//		fmt.Println(moniker)

		addrBytes := common.HexToHash(key).Bytes()
		//		0	0	Slot 0 whiteList mapping so empty
		// WhitelistPerson struct
		/*
		   struct WhitelistPerson {
		     address person;
		     uint  flags;
		   }
		*/
		// Currently no flags are set so the value will be 0 and thus the slot will be empty.

		//		fmt.Println(hex.EncodeToString(slot0Bytes))
		//		fmt.Println(hex.EncodeToString(addrBytes))

		addrHash := eth_crypto.Keccak256(append(addrBytes, slot0Bytes...))
		addrSlot := hex.EncodeToString(addrHash)

		//		fmt.Println(addrSlot)

		if rows[addrSlot] != nil {
			rows[addrSlot].value = parseAddress(rows[addrSlot].rawvalue)
			rows[addrSlot].explained = true
			rows[addrSlot].description = fmt.Sprintf("whiteList mapping set for %s", moniker)
		}
		//		3	0	Slot 3 nomineeList mapping so empty
		// NomineeElection struct
		/*
		   struct NomineeElection{
		     address nominee;
		     address proposer;
		     uint yesVotes;
		     uint noVotes;
		     mapping (address => NomineeVote) vote;
		     address[] yesArray;
		     address[] noArray;
		   }
		*/

		type NomineeArrayParse struct {
			slot        string
			description string
		}

		addrHash = eth_crypto.Keccak256(append(addrBytes, slot3Bytes...))
		addrSlot = hex.EncodeToString(addrHash)

		NomineeArrayParses := []NomineeArrayParse{
			NomineeArrayParse{
				slot:        addrSlot,
				description: "Nominee Mapping",
			},
		}
		//		6	0	Slot 6 evictionList mapping so empty
		// NomineeElection struct

		addrHash = eth_crypto.Keccak256(append(addrBytes, slot6Bytes...))
		addrSlot = hex.EncodeToString(addrHash)

		NomineeArrayParses = append(NomineeArrayParses,
			NomineeArrayParse{
				slot:        addrSlot,
				description: "Eviction Mapping",
			})

		for _, options := range NomineeArrayParses {
			modifier := "Orphaned "

			if rows[options.slot] != nil {
				rows[options.slot].value = parseAddress(rows[options.slot].rawvalue)
				rows[options.slot].explained = true
				rows[options.slot].description = fmt.Sprintf("%s Nominee set for %s", options.description, moniker)
				modifier = ""
			}

			//			     address proposer;
			bigBase := new(big.Int)
			bigBase, _ = baseSlot.SetString(options.slot, 16)
			bigOne := new(big.Int)
			bigOne.SetInt64(int64(1))
			newSlot := bigBase.Add(bigBase, bigOne)
			newHex := newSlot.Text(16)

			if rows[newHex] != nil {
				rows[newHex].value = parseAddress(rows[newHex].rawvalue)
				rows[newHex].explained = true
				rows[newHex].description = fmt.Sprintf("%s%s Proposer set for %s", modifier, options.description, moniker)
			}

			//			uint yesVotes;
			newSlot = newSlot.Add(newSlot, bigOne)
			newHex = newSlot.Text(16)

			if rows[newHex] != nil {
				rows[newHex].value = rows[newHex].rawvalue
				rows[newHex].explained = true
				rows[newHex].description = fmt.Sprintf("%s%s Yes Votes set for %s", modifier, options.description, moniker)
			}
			//			uint noVotes;
			newSlot = newSlot.Add(newSlot, bigOne)
			newHex = newSlot.Text(16)

			if rows[newHex] != nil {
				rows[newHex].value = rows[newHex].rawvalue
				rows[newHex].explained = true
				rows[newHex].description = fmt.Sprintf("%s%s No Votes set for %s", modifier, options.description, moniker)
			}

			//			mapping (address => NomineeVote) vote;
			newSlot = newSlot.Add(newSlot, bigOne)
			newHex = newSlot.Text(16)

			// We try all possible addresses for this mapping

			for voterAddr, voterMoniker := range addressMap {

				voterHash := eth_crypto.Keccak256(append(common.HexToHash(voterAddr).Bytes(), common.HexToHash(newHex).Bytes()...))
				voterSlot := hex.EncodeToString(voterHash)

				if rows[voterSlot] != nil {
					rows[voterSlot].value = parseAddress(rows[voterSlot].rawvalue)
					rows[voterSlot].explained = true
					rows[voterSlot].description = fmt.Sprintf("%s%s voter mapping %s voted on %s", modifier, options.description, voterMoniker, moniker)
				}

			}

			//			address[] yesArray;
			newSlot = newSlot.Add(newSlot, bigOne)
			newHex = newSlot.Text(16)

			if rows[newHex] != nil {
				rows[newHex].value = rows[newHex].rawvalue
				rows[newHex].explained = true
				rows[newHex].description = fmt.Sprintf("%s%s Yes Votes Array Length set for %s", modifier, options.description, moniker)

				slotBytes := eth_crypto.Keccak256(common.HexToHash(newHex).Bytes())
				slotHex := strings.TrimPrefix(hexutil.Encode(slotBytes), "0x")
				bigSlot := new(big.Int)
				bigSlot.SetString(slotHex, 16)

				yesVoteCount, _ := strconv.Atoi(rows[newHex].value)

				for k := 0; k < yesVoteCount; k++ {
					bigYes := new(big.Int)
					bigYes.SetInt64(int64(k))
					bigYes.Add(bigYes, bigSlot)
					yesHex := bigYes.Text(16)

					if rows[yesHex] != nil {
						rows[yesHex].value = parseAddress(rows[yesHex].rawvalue)
						rows[yesHex].explained = true
						rows[yesHex].description = fmt.Sprintf("%s%s Yes Votes Array Index %d set for %s", modifier, options.description, k, moniker)
					}

				}

			}

			//			address[] noArray;
			newSlot = newSlot.Add(newSlot, bigOne)
			newHex = newSlot.Text(16)

			if rows[newHex] != nil {
				rows[newHex].value = rows[newHex].rawvalue
				rows[newHex].explained = true
				rows[newHex].description = fmt.Sprintf("%s%s Yes Votes Array Length set for %s", modifier, options.description, moniker)

				slotBytes := eth_crypto.Keccak256(common.HexToHash(newHex).Bytes())
				slotHex := strings.TrimPrefix(hexutil.Encode(slotBytes), "0x")
				bigSlot := new(big.Int)
				bigSlot.SetString(slotHex, 16)

				yesVoteCount, _ := strconv.Atoi(rows[newHex].value)

				for k := 0; k < yesVoteCount; k++ {
					bigYes := new(big.Int)
					bigYes.SetInt64(int64(k))
					bigYes.Add(bigYes, bigSlot)
					yesHex := bigYes.Text(16)

					if rows[yesHex] != nil {
						rows[yesHex].value = parseAddress(rows[yesHex].rawvalue)
						rows[yesHex].explained = true
						rows[yesHex].description = fmt.Sprintf("%s%s No Votes Array Index %d set for %s", modifier, options.description, k, moniker)
					}

				}

			}

		}

		//		fmt.Println(addrSlot)
	}

	fmt.Printf("Whitelist Count : %d\n", whiteListCount)
	fmt.Printf("Nominee Count   : %d\n", nomineeListCount)
	fmt.Printf("Evictee Count   : %d\n", evicteeListCount)

	fmt.Printf("\n\nData\n====\n\n")

	for key1, rec := range rows {
		fmt.Printf("%s\t%s\t%s\t%s\t%v\n", key1, rec.rawvalue, rec.value, rec.description, rec.explained)
	}

	/*
		fmt.Printf("\n\nAddresses\n=========\n\n")

		for key := range addressMap {
			fmt.Println(key)
		}
	*/
	return nil

}

func parseGenesisOld(cmd *cobra.Command, args []string) error {
	genesisFile := args[0]

	verbose := true

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
	nominees := make(map[uint64]genesis.MinimalPeerRecord)
	evictees := make(map[uint64]genesis.MinimalPeerRecord)

	SLOT2 := fmt.Sprintf("%064d", 2)
	slot2Bytes := eth_crypto.Keccak256(common.HexToHash(SLOT2).Bytes())
	SLOT4 := fmt.Sprintf("%064d", 4)
	slot4Bytes := eth_crypto.Keccak256(common.HexToHash(SLOT4).Bytes())
	SLOT7 := fmt.Sprintf("%064d", 7)
	slot7Bytes := eth_crypto.Keccak256(common.HexToHash(SLOT7).Bytes())

	if verbose {
		fmt.Printf("Array slots:\n  whiteListArray\n    Slot 2: %s\n  nomineeArray\n    Slot 4: %s\n  evictionArray\n    Slot 7: %s\n\n",
			hexutil.Encode(slot2Bytes), hexutil.Encode(slot4Bytes), hexutil.Encode(slot7Bytes))
	}
	// Populate a mapping peers, from the solidity array

	for key, val := range genesisJSON.Poa.Storage {
		// The hardcoded values corresponds to the solidity array.
		// Easier to pick up than the mapping.

		// Evictees
		if key[0:59] == "a66cc928b5edb82af9bd49922954155ab7b0942694bea4ce44661d9a873" {
			// Parse last 5 hex digits, and rebase to zero
			// For all practical purposes 10,000 nominee entries is sufficient
			intkey, err := strconv.ParseUint(key[59:], 16, 64)
			//	intkey -= 774555 // Rebase intkey to 1 for convenience
			if err != nil {
				fmt.Println("Error parsing hex " + key[60:])
				return err
			}
			fmt.Printf("%d\n", intkey)
			// Actually populate the nominee array
			peer, ok := evictees[intkey]
			if !ok {
				evictees[intkey] = genesis.MinimalPeerRecord{
					Address: parseAddress(val),
					Moniker: fmt.Sprintf("Peer %03d", intkey),
				}
			} else {
				peer.Address = parseAddress(val)
				evictees[intkey] = peer
			}

		}
		// Nominees
		if key[0:59] == "8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36" {
			// Parse last 5 hex digits, and rebase to zero
			// For all practical purposes 10,000 nominee entries is sufficient
			intkey, err := strconv.ParseUint(key[59:], 16, 64)
			intkey -= 774555 // Rebase intkey to 1 for convenience
			if err != nil {
				fmt.Println("Error parsing hex " + key[60:])
				return err
			}
			//		fmt.Printf("%d\n", intkey)
			// Actually populate the nominee array
			peer, ok := nominees[intkey]
			if !ok {
				nominees[intkey] = genesis.MinimalPeerRecord{
					Address: parseAddress(val),
					Moniker: fmt.Sprintf("Peer %03d", intkey),
				}
			} else {
				peer.Address = parseAddress(val)
				nominees[intkey] = peer
			}

		}

		// Whitelist
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
					Address: parseAddress(val),
					Moniker: fmt.Sprintf("Peer %03d", intkey),
				}
			} else {
				peer.Address = parseAddress(val)
				peers[intkey] = peer
			}
		}
	}

	// Slot 5 is the moniker mapping.
	SLOT5 := fmt.Sprintf("%064d", 5)
	slot5Bytes := common.HexToHash(SLOT5).Bytes()

	fmt.Printf("%d peers found \n\n", len(peers))
	outputAddresses(peers, slot5Bytes, genesisJSON, verbose)

	fmt.Printf("\n%d nominees found \n\n", len(nominees))
	outputAddresses(nominees, slot5Bytes, genesisJSON, verbose)

	fmt.Printf("\n%d evictees found \n\n", len(evictees))
	outputAddresses(evictees, slot5Bytes, genesisJSON, verbose)

	fmt.Println("")

	if verbose {
		fmt.Println("Slot 0 - whiteList mapping")
		fmt.Println("")

		SLOT0 := fmt.Sprintf("%064d", 0)
		slot0Bytes := common.HexToHash(SLOT0).Bytes()
		for _, peer := range peers {
			addr := strings.TrimPrefix(strings.ToLower(peer.Address), "0x")
			addrBytes := common.HexToHash(addr).Bytes()
			// Handle the Moniker mapping
			addrHash := eth_crypto.Keccak256(append(addrBytes, slot0Bytes...))
			addrSlot := hex.EncodeToString(addrHash)
			fmt.Printf("%s  %s  %s\n", addrSlot, peer.Address, "whiteList")
		}

		fmt.Println("")
		fmt.Println("Slot 3 - nomineeList mapping")
		fmt.Println("")

		SLOT3 := fmt.Sprintf("%064d", 3)
		slot3Bytes := common.HexToHash(SLOT3).Bytes()
		for _, peer := range nominees {
			addr := strings.TrimPrefix(strings.ToLower(peer.Address), "0x")
			addrBytes := common.HexToHash(addr).Bytes()
			// Handle the Moniker mapping
			addrHash := eth_crypto.Keccak256(append(addrBytes, slot3Bytes...))
			addrSlot := hex.EncodeToString(addrHash)

			fmt.Printf("%s  %s  %s\n", addrSlot, peer.Address, "nominees")
		}

		fmt.Println("")
		fmt.Println("Slot 6 - evictionList mapping")
		fmt.Println("")

		SLOT6 := fmt.Sprintf("%064d", 6)
		slot6Bytes := common.HexToHash(SLOT6).Bytes()
		for _, peer := range evictees {
			addr := strings.TrimPrefix(strings.ToLower(peer.Address), "0x")
			addrBytes := common.HexToHash(addr).Bytes()
			// Handle the Moniker mapping
			addrHash := eth_crypto.Keccak256(append(addrBytes, slot6Bytes...))
			addrSlot := hex.EncodeToString(addrHash)

			fmt.Printf("%s  %s  %s\n", addrSlot, peer.Address, "evictees")
		}

	}

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

// Returns slot for Moniker
func getMonikerSlot(address string) string {
	addr := strings.TrimPrefix(strings.ToLower(address), "0x")
	addrBytes := common.HexToHash(addr).Bytes()
	SLOT5 := fmt.Sprintf("%064d", 5)
	slot5Bytes := common.HexToHash(SLOT5).Bytes()

	addrHash := eth_crypto.Keccak256(append(addrBytes, slot5Bytes...))
	addrSlot := hex.EncodeToString(addrHash)

	return addrSlot
}

func outputAddresses(peers map[uint64]genesis.MinimalPeerRecord, slot5Bytes []byte,
	genesisJSON genesis.JSONGenesisFile, verbose bool) {
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
		if verbose {
			fmt.Printf("%s  %s  %s\n", addrSlot, peer.Address, moniker)

		} else {
			fmt.Printf("%s  %s\n", peer.Address, moniker)
		}
	}

}

func parseAddress(in string) string {

	inlen := len(in)

	if inlen == 42 {
		return fmt.Sprintf("0x%s", in[2:])
	}

	if inlen == 44 {
		return fmt.Sprintf("0x%s", in[4:])
	}

	return "Unknown Address"

}
