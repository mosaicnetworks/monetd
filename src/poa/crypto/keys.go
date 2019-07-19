// Package crypto provides common functions for manipulating and generating keys
package crypto

import (
	"encoding/hex"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// PublicKeyHexToAddressHex takes a Hex string public key and returns a hex string
// Ethereum style address.
func PublicKeyHexToAddressHex(publicKey string) (string, error) {
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	pubKeyHash := crypto.Keccak256(pubBytes[1:])[12:]

	return ethcommon.BytesToAddress(pubKeyHash).Hex(), nil
}
