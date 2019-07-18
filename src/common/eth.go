package common

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"
)

//IsValidAddress checks if is a valid ethereum style address
func IsValidAddress(v string) bool {

	re := regexp.MustCompile("^(0[xX]){0,1}[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

//GetNodeSafeLabel converts a free format string into a node label friendly format
//Anything other than an alphanumeric is converted to _
func GetNodeSafeLabel(moniker string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(moniker, "_")
}

//PublicKeyHexToAddressHex takes a Hex string public key and returns a hex string Ethereum style address
func PublicKeyHexToAddressHex(publicKey string) (string, error) {
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	pubKeyHash := crypto.Keccak256(pubBytes[1:])[12:]

	return common.BytesToAddress(pubKeyHash).Hex(), nil

}

//GetMyIP returns the IP address of this instance as a string.
func GetMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// PromptPassphrase prompts the user for a passphrase.  Set confirmation to true
// to require the user to confirm the passphrase.
func PromptPassphrase(confirmation bool) (string, error) {
	passphrase, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		return "", fmt.Errorf("Failed to read passphrase: %v", err)
	}

	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			return "", fmt.Errorf("Failed to read passphrase confirmation: %v", err)
		}
		if passphrase != confirm {
			return "", fmt.Errorf("Passphrases do not match")
		}
	}

	return passphrase, nil
}

// GetPassphrase obtains a passphrase given by the user.  It first checks the
// --passfile command line flag and ultimately prompts the user for a
// passphrase.
func GetPassphrase(passwordFile string) (string, error) {
	// Look for the --passfile flag.
	if passwordFile != "" {
		content, err := ioutil.ReadFile(passwordFile)
		if err != nil {
			return "", fmt.Errorf("Failed to read passphrase file '%s': %v", passwordFile, err)
		}
		return strings.TrimRight(string(content), "\r\n"), nil
	}

	// Otherwise prompt the user for the passphrase.
	return PromptPassphrase(false)
}
