package network

import (
	"errors"
	"os"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "generate new configuration",
		Long: `monetcli network new

Creates a new configuration.`,
		Args: cobra.ExactArgs(0),
		RunE: newConfig,
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "force the creation of a new config file")

	return cmd
}

// func newconfig(cmd *cobra.Command, args []string) error {
func newConfig(cmd *cobra.Command, args []string) error {

	//	fullConfigDir := filepath.Join(configDir, tomlDir)
	fullConfigDir := configDir
	message("Creating a new network configuration")

	if configDir == "" {
		return errors.New("No configuration directory specified. Aborting")
	}

	message("Requested directory", fullConfigDir)

	createDir := true
	// Makes sure that we either have an empty directory or we error out
	if common.CheckIfExists(fullConfigDir) {
		message("Directory Exists")
		isDir, err := common.CheckIsDir(fullConfigDir)
		if !isDir {
			return errors.New("requested directory is an extant file")
		}

		isEmpty, err := isEmptyDir(fullConfigDir)
		if err != nil {
			return err
		}
		if !isEmpty {
			message("Directory is not empty")
			if !force {
				message("No force option specified")
				if err == nil {
					return errors.New("configuration directory exists. \nUse the --force option to override")
				}
				return err
			}
			// Rename the existing file
			err := common.SafeRenameDir(fullConfigDir)
			if err != nil {
				return err

			}

		} else {
			message("Directory exists and is empty")
			createDir = false
		}
	} else {
		message("Directory does not exist")
	}

	if createDir {
		err := os.MkdirAll(fullConfigDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	common.CreateNewConfig(configDir)

	return nil
}
