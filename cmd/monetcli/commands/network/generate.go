package network

import (
	"encoding/hex"
	"errors"
	"path/filepath"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	keys "github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"
	"github.com/spf13/cobra"
)

func generatekeypair(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	ip := args[1]
	isValidator, _ := strconv.ParseBool(args[2])

	message("Generating key pair for: ", moniker)

	targetDir := filepath.Join(configDir, moniker)

	message("Generate to :", targetDir)

	if checkIfExists(targetDir) {
		message("Key Pair for " + moniker + " already exists. Aborting.")
		return errors.New("key pair for " + moniker + " already exists")
	}

	targetFile := filepath.Join(targetDir, keys.DefaultKeyfile)

	/*   // Not required, handled by GenerateKeyPair
	message("Creating dir: ", targetDir)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}
	*/

	key, err := keys.GenerateKeyPair(targetFile)

	if err != nil {
		return err
	}

	pubkey := hex.EncodeToString(
		crypto.FromECDSAPub(&key.PrivateKey.PublicKey))

	return addValidatorParamaterised(moniker, key.Address.Hex(),
		pubkey, ip, isValidator)
	//	return nil
}
