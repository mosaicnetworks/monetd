package config

import (
	"path/filepath"
	"time"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/pelletier/go-toml"
)

type keyValue struct {
	key      string
	value    interface{}
	override bool
}

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

	delKeys := []string{"poa.bytecode", "poa.abi", "validators", "config.datadir"}

	setKeys := []keyValue{
		keyValue{key: "datadir", value: monetConfigDir, override: true},
		keyValue{key: "log", value: "debug", override: false},
		keyValue{key: "eth.datadir", value: filepath.Join(monetConfigDir, "eth"), override: true},
		keyValue{key: "eth.genesis", value: filepath.Join(monetConfigDir, "eth", "genesis.json"), override: true},
		keyValue{key: "eth.keystore", value: filepath.Join(monetConfigDir, "eth", "keystore"), override: true},
		keyValue{key: "eth.pwd", value: filepath.Join(monetConfigDir, "eth", "pwd.txt"), override: true},
		keyValue{key: "eth.db", value: filepath.Join(monetConfigDir, "eth", "chaindata"), override: true},
		keyValue{key: "eth.listen", value: ":8080", override: false},
		keyValue{key: "eth.cache", value: "128", override: false},

		keyValue{key: "babble.datadir", value: filepath.Join(monetConfigDir, "babble"), override: true},
		keyValue{key: "babble.listen", value: ":1337", override: false},
		keyValue{key: "babble.service-listen", value: ":8000", override: false},
		keyValue{key: "babble.heartbeat", value: time.Duration(500 * time.Millisecond).String(), override: false},
		keyValue{key: "babble.timeout", value: time.Duration(1000 * time.Millisecond).String(), override: false},
		keyValue{key: "babble.cache-size", value: "50000", override: false},
		keyValue{key: "babble.sync-limit", value: "1000", override: false},

		keyValue{key: "babble.fast-sync", value: false, override: false},
		keyValue{key: "babble.max-pool", value: "2", override: false},
		keyValue{key: "babble.store", value: true, override: true},
	}

	// First we delete the extraneous keys
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

	// Then we set default values and overrides
	common.Message("Setting default values")
	for _, keys := range setKeys {
		if keys.override || (!tree.Has(keys.key)) {
			common.Message("Setting "+keys.key+" to: ", keys.value)
			tree.Set(keys.key, keys.value)
		}
	}

	return nil
}
