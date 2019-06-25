package network

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func addValidator(cmd *cobra.Command, args []string) error {

	moniker := args[0]
	pubkey := args[1]
	ip := args[2]
	isValidator, _ := strconv.ParseBool(args[3])

	return addValidatorParamaterised(moniker, "", pubkey, ip, isValidator)
}

func addValidatorParamaterised(moniker string, addr string, pubkey string, ip string, isValidator bool) error {
	var config configurationRecord

	//	message("addr:   ", addr)
	//	message("pubkey: ", pubkey)

	err := loadConfig()
	if err != nil {
		return err
	}

	err = networkViper.Unmarshal(&config)

	if err != nil {
		message("Cannot unmarshall config: ", err, moniker)
		return err
	}

	// Derive address from the pubkey
	derivedAddr, err := PublicKeyHexToAddressHex(strings.TrimPrefix(strings.ToLower(pubkey), "0x"))
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

	if !isValidAddress(addr) {
		message("Invalid address: ", addr)
		return errors.New("Invalid Address")
	}

	message("config", config)

	monikers := networkViper.GetString("validators.monikers")
	addresses := networkViper.GetString("validators.addresses")
	pubkeys := networkViper.GetString("validators.pubkeys")
	isvalidators := networkViper.GetString("validators.isvalidator")
	ips := networkViper.GetString("validators.ips")

	if (strings.TrimSpace(monikers)) != "" {
		monikers += ";"
	}

	if (strings.TrimSpace(addresses)) != "" {
		addresses += ";"
	}

	if (strings.TrimSpace(pubkeys)) != "" {
		pubkeys += ";"
	}

	if (strings.TrimSpace(isvalidators)) != "" {
		isvalidators += ";"
	}

	if (strings.TrimSpace(ips)) != "" {
		ips += ";"
	}

	monikers += strings.Replace(moniker, ";", ":", -1)
	addresses += addr
	isvalidators += strconv.FormatBool(isValidator)
	ips += ip
	pubkeys += pubkey

	networkViper.Set("validators.monikers", monikers)
	networkViper.Set("validators.addresses", addresses)
	networkViper.Set("validators.isvalidator", isvalidators)
	networkViper.Set("validators.ips", ips)
	networkViper.Set("validators.pubkeys", pubkeys)

	//	validator := validatorRecord{address: addr, isInitialValidator: true}
	//	config.validators[moniker] = &validator

	writeConfig()
	return nil
}

func PublicKeyHexToAddressHex(publicKey string) (string, error) {

	//	message("Pub Key: ", publicKey)

	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	pubKeyHash := crypto.Keccak256(pubBytes[1:])[12:]

	return common.BytesToAddress(pubKeyHash).Hex(), nil

}
