package network

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	com "github.com/mosaicnetworks/monetd/src/common"
)

func addValidator(cmd *cobra.Command, args []string) error {

	moniker := args[0]
	pubkey := args[1]
	ip := args[2]
	isValidator, _ := strconv.ParseBool(args[3])

	safeLabel := com.GetNodeSafeLabel(moniker)

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

	for _, node := range currentNodes {
		if node == safeLabel {
			com.Message("That Moniker has already been used", moniker)
			return errors.New("that moniker has already been used")
		}
	}

	return AddValidatorParamaterised(configDir, moniker, safeLabel, "", pubkey, ip, isValidator)
}

//AddValidatorParamaterised adds an entry to the peers list.
func AddValidatorParamaterised(configDir string, moniker string, labelsafe string, addr string, pubkey string, ip string, isValidator bool) error {

	tree, err := com.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	if tree.HasPath([]string{"validators", labelsafe}) {
		// Duplicate Node
		return errors.New("cannot add a node with a duplicate moniker")
	}

	// Derive address from the pubkey
	derivedAddr, err := com.PublicKeyHexToAddressHex(strings.TrimPrefix(strings.ToLower(pubkey), "0x"))
	if err != nil {
		message("invalid pubkey to address conversion: ", pubkey)
		return err
	}

	if addr == "" {
		addr = derivedAddr
	} else {
		if strings.TrimPrefix(strings.ToUpper(addr), "0X") != strings.TrimPrefix(strings.ToUpper(derivedAddr), "0X") {
			message("Address derived from public key does not match supplied address. Aborting.")
			message(addr)
			message(derivedAddr)
			return errors.New("derived address does not match supplied address")
		}
	}

	if !com.IsValidAddress(addr) {
		message("Invalid address: ", addr)
		return errors.New("Invalid Address")
	}

	tree.SetPath([]string{"validators", labelsafe, "label"}, labelsafe)
	tree.SetPath([]string{"validators", labelsafe, "moniker"}, moniker)
	tree.SetPath([]string{"validators", labelsafe, "address"}, addr)
	tree.SetPath([]string{"validators", labelsafe, "pubkey"}, pubkey)
	tree.SetPath([]string{"validators", labelsafe, "ip"}, ip)
	tree.SetPath([]string{"validators", labelsafe, "validator"}, isValidator)

	err = com.SaveTomlConfig(configDir, tree)
	if err != nil {
		return err
	}

	return nil
}
