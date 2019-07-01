package common

import (
	"path/filepath"

	"github.com/pelletier/go-toml"
)

func LoadToml(tomlFile string) (*toml.Tree, error) {
	config, err := toml.LoadFile(tomlFile)

	if err != nil {
		Message("Error loading toml file: ", tomlFile)
		return nil, err
	}

	return config, nil
}

func SaveToml(tree *toml.Tree, tomlFile string) error {

	// Open Writer

	tomlStr, err := tree.ToTomlString()
	if err != nil {
		Message("Cannot parse toml output file: ", tomlFile)
		return err
	}

	err = WriteToFile(tomlFile, tomlStr)

	if err != nil {
		Message("Failed to write toml output file", tomlFile)
		return err
	}

	Message("Written toml file: ", tomlFile)
	Message("Characters written ", len(tomlStr))

	return nil
}

func TransformCliTomlToD(tree *toml.Tree, monetConfigDir string) error {

	delKeys := []string{"poa.bytecode", "poa.abi", "validators", "config.datadir"}

	setKeys := GetMonetDefaultConfigKeys(monetConfigDir)

	// First we delete the extraneous keys
	for _, key := range delKeys {

		if tree.Has(key) {
			err := tree.Delete(key)
			if err != nil {
				Message("Error deleting "+key+": ", err)
				return err
			}
		} else {
			Message("Key " + key + " is not set. Continuing.")
		}
	}

	// Then we set default values and overrides
	Message("Setting default values")
	for _, keys := range setKeys {
		if keys.Override || (!tree.Has(keys.Key)) {
			Message("Setting "+keys.Key+" to: ", keys.Value)
			tree.Set(keys.Key, keys.Value)
		}
	}

	return nil
}

func LoadTomlConfig(configDir string) (*toml.Tree, error) {

	Message("Starting to load configuration")

	tree, err := LoadToml(filepath.Join(configDir, MonetcliTomlName+TomlSuffix))

	if err != nil {
		Message("loadConfig: ", err)
		return nil, err
	}
	Message("Loaded Config")
	return tree, nil
}

func SaveTomlConfig(configDir string, tree *toml.Tree) error {

	Message("Starting to save configuration")

	err := SaveToml(tree, filepath.Join(configDir, MonetcliTomlName+TomlSuffix))

	if err != nil {
		Message("saveConfig: ", err)
		return err
	}
	Message("Saved Config")
	return nil
}
