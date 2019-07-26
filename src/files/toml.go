package files

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/pelletier/go-toml"
)

//LoadToml loads a Toml file and returns a tree object.
func LoadToml(tomlFile string) (*toml.Tree, error) {
	config, err := toml.LoadFile(tomlFile)

	if err != nil {
		common.ErrorMessage("Error loading toml file: ", tomlFile)
		return nil, err
	}

	return config, nil
}

//SaveToml writes a tree object (back) to a toml file
func SaveToml(tree *toml.Tree, tomlFile string) error {

	// Open Writer
	tomlStr, err := tree.ToTomlString()
	if err != nil {
		common.ErrorMessage("Cannot parse toml output file: ", tomlFile)
		return err
	}

	err = WriteToFile(tomlFile, tomlStr)

	if err != nil {
		common.ErrorMessage("Failed to write toml output file", tomlFile)
		return err
	}

	common.DebugMessage("Written toml file: ", tomlFile)
	common.DebugMessage("Characters written ", len(tomlStr))

	return nil
}
