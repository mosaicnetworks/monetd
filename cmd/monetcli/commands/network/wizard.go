package network

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "wizard to set up a Monet Network",
	Long: `Wizard to set up a Monet Network
	
	This command provides a wizard interface to the 
	"monetcli network" commands. This provides a guided interface
	through the process of configuring a network.`,
	RunE: runWizardCmd,
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
		userConfigDir = common.RequestFile("MonetCLI Configuration Directory Location: ", userConfigDir)

		DirExists := common.CheckIfExists(userConfigDir)

		if !DirExists {

			common.MessageWithType(common.MsgInformation, "Folder "+userConfigDir+" does not exist.")

			confirm := common.RequestSelect("Please select an option: ", []string{common.WizardExitWithoutSavingChanges, common.WizardChangeConfigDir, common.WizardTextCreateNewConfiguration}, common.WizardTextCreateNewConfiguration)

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
			confirmSelection := common.RequestSelect("Please select an option: ",
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
			common.MessageWithType(common.MsgError, "Error in compiling: ", err)
			return err
		}
		if shouldBreak {
			break configloop
		}

	}
	return nil
}

func editWizard(configDir string) (bool, error) {

	configFile := filepath.Join(configDir, common.MonetcliTomlName+common.TomlSuffix)

editloop:
	for {

		common.MessageWithType(common.MsgInformation, "Edit menu for   "+configDir+" ")
		confirmSelection := common.RequestSelect("Please select an option: ",
			[]string{
				common.WizardAddKeyPair,
				common.WizardCheckConfiguration,
				common.WizardCompile,
				common.WizardGenerate,
				common.WizardParams,
				common.WizardPeers,
				common.WizardShow,
				common.WizardExit,
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
		case common.WizardPeers:
			err := PeersWizard(configDir)
			if err != nil {
				common.MessageWithType(common.MsgError, "Error in Peers wizard: ", common.MsgError)
				return false, err
			}
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
				common.MessageWithType(common.MsgError, "Error Compiling: ", err)
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

		_ = common.RequestFile("Press Enter to Continue", "")

	}

	return false, nil
}

func monetDConfigWizard() error {
	//TODO - New code

	// The delivery process putting files into the correct place is by and large complete for the
	// testnet command - parameterise it and move it into common

	common.MessageWithType(common.MsgInformation, "MonetD Config will be here.")

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
		nodename = common.RequestString("Node Name: ", "")

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
		ip = common.RequestString("Node Address: ", nodename+":1337")
		//Enhancement -  Check if the ip is already in use
		break iploop
	}

	// We decided that non-validators is a configuration too far.
	// We will leave the parameter in place, but set it silently
	// just in case we change our minds later.
	isValidator = true
	// isValidator = common.RequestBool("Is validator in initial peer set: ", true)

	var password string

passwordloop:
	for {
		password = common.RequestPassword("Enter Keystore Password: ", "")
		password2 := common.RequestPassword("Confirm Keystore Password: ", "")

		if password == password2 {
			break passwordloop
		}
	}

	passwordFile := filepath.Join(configDir, common.PwdFile)

	err = common.WriteToFile(passwordFile, password)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error saving password: ", err)
		return err
	}

	return GenerateKeyPair(configDir, nodename, ip, isValidator, passwordFile)
}

func compileWizard(configDir string) error {
	return CompileConfigWithParam(configDir)
}

func addPeerWizard(configDir string) error {

	//TODO - hook up wizard to add peer function req

	//	func AddValidatorParamateriseda(configDir, moniker string, addr string, pubkey string, ip string, isValidator bool) error {
	return nil
}
