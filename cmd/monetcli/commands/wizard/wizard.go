package wizard

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/network"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

const (
	WizardExit                                            = "Exit"
	WizardTextCreateNewConfiguration                      = "Create New Configuration"
	WizardExitWithoutSavingChanges                        = "Exit Without Saving Changes"
	WizardEditExistingConfiguration                       = "Edit Existing Configuration"
	WizardChangeConfigDir                                 = "Change Configuration Directory"
	WizardRenameCurrentDirectoryandCreateNewConfiguration = "Rename Current Directory and Create New Configuration"
	WizardAddKeyPair                                      = "Add Key Pair"
	WizardCheckConfiguration                              = "Check Configuration"
	WizardCompile                                         = "Compile POA Contract"
	WizardGenerate                                        = "Generate Key Pair"
	WizardParams                                          = "Edit Params"
	WizardShow                                            = "Show Configuration"
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

			confirm := requestSelect("Please select an option: ", []string{WizardExitWithoutSavingChanges, WizardChangeConfigDir, WizardTextCreateNewConfiguration}, WizardTextCreateNewConfiguration)

			switch confirm {
			case WizardExit:
				break configloop
			case WizardChangeConfigDir:
				continue configloop
			}

			err := network.CreateNewConfig(userConfigDir)
			if err != nil {
				common.Message("Error creating new config: ", err)
				return err
			}

		} else {
			common.MessageWithType(common.MsgInformation, "Folder "+userConfigDir+" already exists.")
			confirmSelection := requestSelect("Please select an option: ",
				[]string{
					WizardExitWithoutSavingChanges,
					WizardChangeConfigDir,
					WizardEditExistingConfiguration,
					WizardRenameCurrentDirectoryandCreateNewConfiguration},
				WizardRenameCurrentDirectoryandCreateNewConfiguration)

			switch confirmSelection {
			case WizardExitWithoutSavingChanges:
				break configloop
			case WizardChangeConfigDir:
				continue configloop
			case WizardRenameCurrentDirectoryandCreateNewConfiguration:

				// Safe Rename
				err := common.SafeRenameDir(userConfigDir)
				if err != nil {
					common.Message("Error renaming old config: ", err)
					return err
				}

				err = network.CreateNewConfig(userConfigDir)
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
		confirmSelection := requestSelect("Please select an option: ",
			[]string{ //TODO remove these comments
				WizardAddKeyPair,         // Skeleton
				WizardCheckConfiguration, // Skeleton
				WizardCompile,            // Skeleton
				WizardGenerate,           // Skeleton
				WizardParams,             // Complete
				WizardShow,               // Complete
				WizardExit,               // Complete
			},
			WizardShow)

		switch confirmSelection {
		case WizardExit:
			return true, nil
			break editloop
		case WizardShow:
			_ = common.ShowConfigFile(configFile)
		case WizardParams:
			network.SetParamsWithParams(configDir)
		case WizardAddKeyPair:
			err := addPeerWizard(configDir)
			if err != nil {
				return false, err
			}
		case WizardCheckConfiguration:
			err := checkConfigurationWizard(configDir)
			if err != nil {
				return false, err
			}
		case WizardCompile:
			err := compileWizard(configDir)
			if err != nil {
				return false, err
			}

			err = monetDConfigWizard() // Move onto monetd config
			if err != nil {
				return false, err
			}

			break editloop

		case WizardGenerate:
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
	return network.CheckConfigWithParams(configDir)
}

func generateKeyPairWizard(configDir string) error {
	var nodename, ip string
	var isValidator bool

nodeloop:
	for {
		nodename = requestString("Node Name: ", "")
		//TODO Check if the moniker is already in use

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

	return network.GenerateKeyPair(configDir, nodename, ip, isValidator)
}

func compileWizard(configDir string) error {
	return network.CompileConfigWithParam(configDir)
}

func addPeerWizard(configDir string) error {

	//TODO

	//	func addValidatorParamaterised(moniker string, addr string, pubkey string, ip string, isValidator bool) error {
	return nil
}
