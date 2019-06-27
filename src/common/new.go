package common

import (
	"os"

	"github.com/pelletier/go-toml"
)

func CreateNewConfig(configDir string) error {

	/*
		tree, err := LoadTomlConfig(configDir)
		if err != nil {
			common.Message("Error Loading TOML")
		}
	*/

	if !CheckIfExists(configDir) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			Message("Error creating empty config folder: ", err)
			return err
		}
	}

	tree, err := toml.Load("")

	if err != nil {
		Message("Error in CreateNewConfig: ", err)
		return err
	}

	tree.Set("poa.contractaddress", DefaultContractAddress)

	err = SaveTomlConfig(configDir, tree)
	if err != nil {
		Message("Error saving in CreateNewConfig: ", err)
		return err
	}

	return nil
}
