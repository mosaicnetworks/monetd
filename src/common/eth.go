package common

import (
	"encoding/hex"
	"log"
	"net"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
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
