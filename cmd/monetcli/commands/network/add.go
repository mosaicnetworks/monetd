package network

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func addValidator(cmd *cobra.Command, args []string) error {

	moniker := args[0]
	addr := args[1]
	ip := args[2]
	isValidator, _ := strconv.ParseBool(args[3])

	return addValidatorParamaterised(moniker, addr, ip, isValidator)
}

func addValidatorParamaterised(moniker string, addr string, ip string, isValidator bool) error {
	var config configurationRecord

	err := loadConfig()
	if err != nil {
		return err
	}

	err = networkViper.Unmarshal(&config)

	if err != nil {
		message("Cannot unmarshall config: ", err, moniker)
		return err
	}

	if !isValidAddress(addr) {
		message("Invalid address: ", addr)
		return errors.New("Invalid Address")
	}

	message("config", config)

	monikers := networkViper.GetString("validators.monikers")
	addresses := networkViper.GetString("validators.addresses")
	isvalidators := networkViper.GetString("validators.isvalidator")
	ips := networkViper.GetString("validators.ips")

	if (strings.TrimSpace(monikers)) != "" {
		monikers += ";"
	}

	if (strings.TrimSpace(addresses)) != "" {
		addresses += ";"
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

	networkViper.Set("validators.monikers", monikers)
	networkViper.Set("validators.addresses", addresses)
	networkViper.Set("validators.isvalidator", isvalidators)
	networkViper.Set("validators.ips", ips)

	//	validator := validatorRecord{address: addr, isInitialValidator: true}
	//	config.validators[moniker] = &validator

	writeConfig()
	return nil
}
