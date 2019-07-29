package crypto

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/ethereum/go-ethereum/console"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,_"

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
func GetPassphrase(passwordFile string, confirmation bool) (string, error) {
	// Look for the --passfile flag.
	if passwordFile != "" {
		content, err := ioutil.ReadFile(passwordFile)
		if err != nil {
			return "", fmt.Errorf("Failed to read passphrase file '%s': %v", passwordFile, err)
		}
		return strings.TrimRight(string(content), "\r\n"), nil
	}

	// Otherwise prompt the user for the passphrase.
	return PromptPassphrase(confirmation)
}

//RandomPassphrase generates a random passphrase
func RandomPassphrase(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,_"

	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
