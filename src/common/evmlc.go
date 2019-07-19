package common

import (
	"path/filepath"
	"strconv"
)

//SendKeyToEVMLC adds a keyfile.json to the evmlc keystore
func SendKeyToEVMLC(nodename string, keyFile string) error {
	// Generally this is non-critical so swallow minor errors

	tomlDir, _ := DefaultHomeDir(EvmlcTomlDir)
	keyDir := filepath.Join(tomlDir, "keystore")

	if !CheckIfExists(keyFile) {
		return nil
	}
	if CreateDirIfNotExists(tomlDir) != nil || CreateDirIfNotExists(keyDir) != nil {
		return nil
	}

	baseFileName := nodename + ".json"
	counter := 0

	for {
		if counter > 0 {
			baseFileName = nodename + "." + strconv.Itoa(counter) + ".json"
		}
		testFileName := filepath.Join(keyDir, baseFileName)
		if !CheckIfExists(testFileName) {
			CopyFileContents(keyFile, testFileName)
			return nil
		}
		counter++
	}

}
