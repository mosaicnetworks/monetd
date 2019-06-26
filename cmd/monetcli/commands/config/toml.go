package config

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/pelletier/go-toml"
)

func LoadToml(tomlFile string) (*toml.Tree, error) {
	config, err := toml.LoadFile(tomlFile)

	if err != nil {
		common.Message("Error loading toml file: ", tomlFile)
		return nil, err
	}

	return config, nil
}

func SaveToml(tree *toml.Tree, tomlFile string) error {

	// Open Writer

	tomlStr, err := tree.ToTomlString()
	if err != nil {
		common.Message("Cannot parse toml output file: ", tomlFile)
		return err
	}

	err = common.WriteToFile(tomlFile, tomlStr)

	if err != nil {
		common.Message("Failed to write toml output file", tomlFile)
		return err
	}

	common.Message("Written toml file: ", tomlFile)
	common.Message("Characters written ", len(tomlStr))

	return nil
}

func transformCliTomlToD(tree *toml.Tree) error {

	delKeys := []string{"poa.bytecode", "poa.abi"}

	for _, key := range delKeys {

		if tree.Has(key) {
			err := tree.Delete(key)
			if err != nil {
				common.Message("Error deleting "+key+": ", err)
				return err
			}
		} else {
			common.Message("Key " + key + " is not set. Continuing.")
		}

	}

	return nil
}
