package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
)

var (
	configFileList = []*configFile{
		&configFile{sourcefilename: common.MonetcliTomlName + ".toml",
			targetfilename: common.MonetdTomlName + ".toml", label: "toml", required: true, transformation: true},
		&configFile{sourcefilename: common.PeersJSON,
			targetfilename: common.PeersJSONTarget, label: "peers", required: true, transformation: false},
		&configFile{sourcefilename: common.PeersGenesisJSON,
			targetfilename: common.PeersGenesisJSONTarget, label: "genesispeers", required: true, transformation: false},
		&configFile{sourcefilename: common.GenesisJSON,
			targetfilename: common.GenesisJSONTarget, label: "genesis", required: true, transformation: false},
	}
)

func publishConfig(cmd *cobra.Command, args []string) error {

	// Check that we have a valid monetcli config

	//First check that we have a file location
	if networkConfigDir == "" {
		common.Message("networkConfigDir is empty")
		return errors.New("config path not set. use --config-dir parameter")
	}

	// Check the location actually exists
	if !common.CheckIfExists(networkConfigDir) {
		common.Message("Network Configuration not found", networkConfigDir)
		return errors.New("network configuration not found")
	}

	// Check the location is a directory
	isDir, err := common.CheckIsDir(networkConfigDir)
	if err != nil {
		common.Message("Failed directory check", networkConfigDir)
		return err
	}

	if !isDir {
		common.Message("Failed directory check", networkConfigDir)
		return errors.New("configuration folder is a flat file")
	}

	// Check the Output Dir

	//First check that we have a file location
	if monetConfigDir == "" {
		common.Message("monetConfigDir is empty")
		return errors.New("monet config path not set. use --monet-config-dir parameter")
	}

	// Check the location actually exists
	if common.CheckIfExists(monetConfigDir) {
		// config directory exists.
		if !force {
			common.Message("directory already exists. ", monetConfigDir)
			return errors.New("output config directory already exists. Use --force to rename existing config")
		} else {
			err := common.SafeRenameDir(monetConfigDir)
			if err != nil {
				common.Message("Cannot rename existing configuration: ", monetConfigDir)
				return errors.New("cannot rename existing config. Try manually renaming it")
			}
		}

	}

	dirList := []string{monetConfigDir, filepath.Join(monetConfigDir, "babble"), filepath.Join(monetConfigDir, "eth")}

	for _, dir := range dirList {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			common.Message("Error creating directory: ", dir)
			return err
		}
		common.Message("Created directory: ", dir)
	}

	// Check than all mandatory files are actually present
	// Currently this working on a single pass. Split to have a validation pass then an action pass.
	for _, file := range configFileList {
		fileWithPath := filepath.Join(networkConfigDir, file.sourcefilename)
		outFileWithPath := filepath.Join(monetConfigDir, file.targetfilename)

		if !common.CheckIfExists(fileWithPath) {
			if file.required {
				common.Message("Incomplete configuration: ", fileWithPath)
				return errors.New("missing file " + file.label + ". Try running: monetcli network compile")
			}
			common.Message("Optional config file " + file.sourcefilename + " not found.")
			continue
		}

		if file.transformation {
			switch file.label {
			case "toml":
				tr, err := LoadToml(fileWithPath)
				if err != nil {
					common.Message("Cannot load toml file")
					return err
				}
				// Delete extraneous keys and apply defaults
				err = transformCliTomlToD(tr)
				if err != nil {
					common.Message("Cannot transform toml file", fileWithPath)
					return err
				}

				err = SaveToml(tr, outFileWithPath)
				if err != nil {
					common.Message("Cannot save toml file")
					return err
				}
			}
		} else {
			err := common.CopyFileContents(fileWithPath, outFileWithPath)
			if err != nil {
				common.Message("Cannot copy file from " + fileWithPath + " to " + outFileWithPath)
				return err
			}
		}
	}

	//	toml := filepath.Join(networkConfigDir, common.MonetcliTomlName+".toml")

	switch publishTarget {
	case "simple":
		// Publish the target

	default:
		return errors.New("unknown publish target")
	}

	return nil
}
