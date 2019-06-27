package network

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Wizard to set up a Monet Network",
	Long:  `Wizard to set up a Monet Network`,
	RunE:  runWizardCmd,
}

func runWizardCmd(cmd *cobra.Command, args []string) error {
	return runWizard()

}

func runWizard() error {

	home, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	userConfigDir := home

	// Infinite loop that we need to explicitly break out of
configloop:
	for {
		userConfigDir = requestFile("MonetCLI Configuration Directory Location: ", userConfigDir)

		DirExists := common.CheckIfExists(userConfigDir)

		if !DirExists {

			common.MessageWithType(common.MsgInformation, "Folder "+userConfigDir+" does not exist.")

			confirm := RequestSelect("Please select an option: ", []string{common.WizardExitWithoutSavingChanges, common.WizardChangeConfigDir, common.WizardTextCreateNewConfiguration}, common.WizardTextCreateNewConfiguration)

			switch confirm {
			case common.WizardExit:
				break configloop
			case common.WizardChangeConfigDir:
				continue configloop
			}

			err := common.CreateNewConfig(userConfigDir)
			if err != nil {
				common.Message("Error creating new config: ", err)
				return err
			}

		} else {
			common.MessageWithType(common.MsgInformation, "Folder "+userConfigDir+" already exists.")
			confirmSelection := RequestSelect("Please select an option: ",
				[]string{
					common.WizardExitWithoutSavingChanges,
					common.WizardChangeConfigDir,
					common.WizardEditExistingConfiguration,
					common.WizardRenameCurrentDirectoryandCreateNewConfiguration},
				common.WizardRenameCurrentDirectoryandCreateNewConfiguration)

			switch confirmSelection {
			case common.WizardExitWithoutSavingChanges:
				break configloop
			case common.WizardChangeConfigDir:
				continue configloop
			case common.WizardRenameCurrentDirectoryandCreateNewConfiguration:

				// Safe Rename
				err := common.SafeRenameDir(userConfigDir)
				if err != nil {
					common.Message("Error renaming old config: ", err)
					return err
				}

				err = common.CreateNewConfig(userConfigDir)
				if err != nil {
					common.Message("Error creating new config: ", err)
					return err
				}

			}

		}

		// If we reach this point, we have a configuration, either new or existing.

		common.Message("Editing: ", userConfigDir)

		shouldBreak, err := editWizard(userConfigDir)
		if err != nil {
			return err
		}
		if shouldBreak {
			break configloop
		}

	}
	return nil
}

func editWizard(configDir string) (bool, error) {

	configFile := filepath.Join(configDir, common.MonetcliTomlName+".toml")

editloop:
	for {

		common.MessageWithType(common.MsgInformation, "Edit menu for   "+configDir+" ")
		confirmSelection := RequestSelect("Please select an option: ",
			[]string{ //TODO remove these comments
				common.WizardAddKeyPair,         // Skeleton
				common.WizardCheckConfiguration, // Skeleton
				common.WizardCompile,            // Skeleton
				common.WizardGenerate,           // Skeleton
				common.WizardParams,             // Complete
				common.WizardShow,               // Complete
				common.WizardExit,               // Complete
			},
			common.WizardShow)

		switch confirmSelection {
		case common.WizardExit:
			return true, nil
			break editloop
		case common.WizardShow:
			_ = common.ShowConfigFile(configFile)
		case common.WizardParams:
			common.SetParamsWithParams(configDir)
		case common.WizardAddKeyPair:
			err := addPeerWizard(configDir)
			if err != nil {
				return false, err
			}
		case common.WizardCheckConfiguration:
			err := checkConfigurationWizard(configDir)
			if err != nil {
				return false, err
			}
		case common.WizardCompile:
			err := compileWizard(configDir)
			if err != nil {
				return false, err
			}

			err = monetDConfigWizard() // Move onto monetd config
			if err != nil {
				return false, err
			}

			break editloop

		case common.WizardGenerate:
			err := generateKeyPairWizard(configDir)
			if err != nil {
				return false, err
			}

		}

		_ = requestFile("Press Enter to Continue", "")

	}

	return false, nil
}

func monetDConfigWizard() error {
	//TODO
	return nil
}

func checkConfigurationWizard(configDir string) error {
	return CheckConfigWithParams(configDir)
}

func generateKeyPairWizard(configDir string) error {
	var nodename, ip string
	var isValidator bool

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

nodeloop:
	for {
		nodename = requestString("Node Name: ", "")
		//TODO Check if the moniker is already in use

		safeLabel := common.GetNodeSafeLabel(nodename)
		for _, node := range currentNodes {
			if node == safeLabel {
				common.MessageWithType(common.MsgError, "That Moniker has already been used", nodename)
				continue nodeloop
			}

		}

		break nodeloop
	}

iploop:
	for {
		ip = requestString("Node Address: ", nodename+":1337")
		//TODO Check if the ip is already in use
		break iploop
	}

	isValidator = requestBool("Is validator in initial peer set: ", true)
	//TODO Check if the ip is already in use

	//TODO generate passwordfile as the last parameter of this call
	return GenerateKeyPair(configDir, nodename, ip, isValidator, "")
}

func compileWizard(configDir string) error {
	return CompileConfigWithParam(configDir)
}

func addPeerWizard(configDir string) error {

	//TODO

	//	func addValidatorParamaterised(moniker string, addr string, pubkey string, ip string, isValidator bool) error {
	return nil
}
