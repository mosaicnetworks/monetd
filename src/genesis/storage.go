package genesis

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	eth_crypto "github.com/ethereum/go-ethereum/crypto"
)

//GetStorage creates a mapping for the storage records for the genesis file
//for a given set of initial validators
func GetStorage(peers []*MinimalPeerRecord) (map[string]string, error) {

	storage := make(map[string]string)

	// Set the array and mapping lengths
	SLOT0 := fmt.Sprintf("%064d", 0)
	SLOT1 := fmt.Sprintf("%064d", 1)
	SLOT2 := fmt.Sprintf("%064d", 2)
	SLOT5 := fmt.Sprintf("%064d", 5)

	peerslength := len(peers)
	hexLength := fmt.Sprintf("%x", peerslength)
	if len(hexLength)%2 == 1 { // Force hex pairs
		hexLength = "0" + hexLength
	}

	// Slot 1 is the Whitelist Count Int
	// Slot 2 is the Whitelist Count Dynamically Sized Array
	storage[SLOT1] = hexLength
	storage[SLOT2] = hexLength

	// Calculate the Hash of (slot) 2. It is the first element of the array
	pubBytes := common.HexToHash(SLOT2).Bytes()
	pubKeyHash := eth_crypto.Keccak256(pubBytes)
	ARRAYSLOT := hex.EncodeToString(pubKeyHash)

	// Set up a big Int, as we need to increment the slot for each address in
	// the list
	arraySlotCounter := new(big.Int)
	arraySlotCounter.SetString(ARRAYSLOT, 16)

	// Set up a big int of value 1 to increment the slot for each array address.
	inc := new(big.Int)
	inc.SetUint64(1)

	// Prebuild the slot 0 byte array
	slot0Bytes := common.HexToHash(SLOT0).Bytes()
	// Prebuild the slot 5 byte array
	slot5Bytes := common.HexToHash(SLOT5).Bytes()

	for _, peer := range peers {
		addr := peer.Address

		// Set the Array and increment the array slot by one
		storage[fmt.Sprintf("%064x", arraySlotCounter)] = "94" + addr //TODO - why is this 94?
		arraySlotCounter.Add(arraySlotCounter, inc)

		// Handle the Whitelist mapping
		addrBytes := common.HexToHash(addr).Bytes()
		addrHash := eth_crypto.Keccak256(append(addrBytes, slot0Bytes...))
		addrSlot := hex.EncodeToString(addrHash)
		storage[addrSlot] = "94" + addr

		// Handle the Moniker mapping
		addrHash = eth_crypto.Keccak256(append(addrBytes, slot5Bytes...))
		addrSlot = hex.EncodeToString(addrHash)
		monikerString := peer.Moniker
		if len(monikerString) > 64 {
			monikerString = monikerString[0:63]
		}

		monikerString = "0a" + strings.TrimLeft(hex.EncodeToString([]byte(monikerString)), "726a")
		storage[addrSlot] = monikerString + strings.Repeat("0", 64-len(monikerString))
	}

	return storage, nil
}
