package network

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/babble/src/babble"
	conf "github.com/mosaicnetworks/monetd/cmd/monetcli/commands/config"
	"github.com/mosaicnetworks/monetd/src/common"
	monet "github.com/mosaicnetworks/monetd/src/version"
	"github.com/spf13/cobra"
)

//WizardCmd is the wizard subcommand
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
	common.BannerTitle("wizard")
	home, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	userConfigDir := home

	// Infinite loop that we need to explicitly break out of
configloop:
	for {
		userConfigDir = common.RequestFile("MonetCLI Configuration Directory Location", userConfigDir)

		if !common.CheckIfExists(userConfigDir) {
			common.MessageWithType(common.MsgInformation, "Folder "+userConfigDir+" does not exist.")
			confirm := common.RequestSelect("Please select an option", []string{
				common.WizardExitWithoutSavingChanges,
				common.WizardChangeConfigDir,
				common.WizardTextCreateNewConfiguration,
			},
				common.WizardTextCreateNewConfiguration)

			switch confirm {
			case common.WizardExitWithoutSavingChanges:
				return nil
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
			confirmSelection := common.RequestSelect("Please select an option",
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

		common.MessageWithType(common.MsgDebug, "Editing monetcli config: ", userConfigDir)

		shouldBreak, err := editWizard(userConfigDir)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error in editing monetcli configuration: ", err)
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

	for {
		common.BannerTitle("edit")
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
				common.WizardVersion,
				common.WizardExit,
			},
			common.WizardShow)

		switch confirmSelection {
		case common.WizardExit:
			return true, nil
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
		case common.WizardVersion:
			common.BannerTitle("monetd")
			fmt.Print(monet.FullVersion())

		case common.WizardCompile:
			err := compileWizard(configDir)
			if err != nil {
				common.MessageWithType(common.MsgError, "Error Compiling: ", err)
				return false, err
			}

			monetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)
			err = monetDConfigWizard(configDir, monetConfigDir) // Move onto monetd config
			if err != nil {
				return false, err
			}
			common.ContinuePrompt()
			return true, nil

		case common.WizardGenerate:
			err := generateKeyPairWizard(configDir)
			if err != nil {
				return false, err
			}

		}

		common.ContinuePrompt()

	}

}

func monetDConfigWizard(networkConfigDir string, monetConfigDir string) error {

	// The delivery process putting files into the correct place is by and large complete for the
	// testnet command - parameterise it and move it into common

	common.MessageWithType(common.MsgInformation, "MonetD Config publish to ", monetConfigDir)

	// Here we have genesis.json, peers.json network.toml in the .monetcli directory ready to go.

	if common.CheckIfExists(monetConfigDir) {
		confirm := common.RequestSelect("Overwrite Monet Configuration: ", []string{"No", "Yes"}, "No")

		if confirm == "No" {
			common.MessageWithType(common.MsgWarning, "MonetD Config publishing cancelled")
			return errors.New("Cancelling publish")
		}

		conf.Force = true
	}

	err := conf.PublishAllConfigWithParams(networkConfigDir, monetConfigDir)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Error in monetDConfigWizard: ", err)
		return err
	}

	common.MessageWithType(common.MsgInformation, "All configurations exported.")

	for {

		selection := common.RequestSelect("Choose an option", []string{common.WizardShow, common.WizardExit}, common.WizardExit)
		if selection == common.WizardExit {
			return nil
		}

		_ = conf.ShowConfigParams(monetConfigDir)

		common.ContinuePrompt()

	}

	//	return nil
}

func checkConfigurationWizard(configDir string) error {
	return CheckConfigWithParams(configDir)
}

func generateKeyPairWizard(configDir string) error {
	var nodename, ip, safeLabel string
	var isValidator bool

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

nodeloop:
	for {
		nodename = common.RequestString("Node Name", "")

		safeLabel = common.GetNodeSafeLabel(nodename)
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
		ip = common.RequestString("Node Address", nodename+":1337")
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
		password = common.RequestPassword("Enter Keystore Passphrase", "")
		password2 := common.RequestPassword("Confirm Keystore Passphrase", "")

		if password == password2 {
			break passwordloop
		}
	}

	outDir := filepath.Join(configDir, safeLabel)
	passwordFile := filepath.Join(outDir, common.PwdFile)
	privkeyfile := filepath.Join(outDir, babble.DefaultKeyfile)

	if !common.CheckIfExists(outDir) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			common.Message("Error creating directory: ", outDir)
			return err
		}
	}

	err = common.WriteToFile(passwordFile, password)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error saving password: ", err)
		return err
	}

	return GenerateKeyPair(configDir, nodename, ip, isValidator, passwordFile, privkeyfile)
}

func compileWizard(configDir string) error {
	return CompileConfigWithParam(configDir)
}

func addPeerWizard(configDir string) error {

	var nodename, ip, pubkey, safeLabel string
	var isValidator = true

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

nodeloop:
	for {
		nodename = common.RequestString("Node Name", "")

		safeLabel = common.GetNodeSafeLabel(nodename)
		common.MessageWithType(common.MsgDebug, "SafeLabel: ", safeLabel)

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
		ip = common.RequestString("Node Address", nodename+":1337")
		//Enhancement -  Check if the ip is already in use
		break iploop
	}

	pubkey = common.RequestString("Pub Key", "")

	err = AddValidatorParamaterised(configDir, nodename, safeLabel, "", pubkey, ip, isValidator)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error adding a peer: ", err)
	}
	return nil
}
