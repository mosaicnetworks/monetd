package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

//Constants for the labels used in Wizard select options.
//Constants are much safer than string literals for parsing the results of select commands.
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
	WizardView                                            = "View"
	WizardEdit                                            = "Edit"
	WizardDelete                                          = "Delete"
	WizardSaveChanges                                     = "Save Changes"
	WizardCancelChanges                                   = "Cancel Changes"
	WizardEditAgain                                       = "Edit Again"
	WizardPeers                                           = "Peers"
)

//RequestFile prompts the user for a file name.
func RequestFile(promptText string, defaultValue string) string {

	prompt := promptui.Prompt{
		Label:    promptText + "  ",
		Validate: validateFile,
		Default:  defaultValue,
		Pointer:  promptui.PipeCursor,
	}

	result, _ := prompt.Run()

	return result
}

//RequestString prompts the user for a string
func RequestString(promptText string, defaultValue string) string {

	prompt := promptui.Prompt{
		Label:    promptText + "  ",
		Validate: validateString,
		Default:  defaultValue,
		Pointer:  promptui.PipeCursor,
	}

	result, _ := prompt.Run()

	return result
}

//RequestPassword prompts the user for input using a mask to hide their answer
func RequestPassword(promptText string, defaultValue string) string {

	prompt := promptui.Prompt{
		Label:    promptText + "  ",
		Validate: validateString,
		Default:  defaultValue,
		Mask:     '#',
		Pointer:  promptui.PipeCursor,
	}

	result, _ := prompt.Run()

	return result
}

//RequestBool prompts the user for a boolean answer and parses the string to
//return a bool type
func RequestBool(promptText string, defaultValue bool) bool {

	defaultStr := strconv.FormatBool(defaultValue)
	strVal := RequestSelect(promptText, []string{"true", "false"}, defaultStr)

	rtnVal, _ := strconv.ParseBool(strVal)

	return rtnVal
}

func validateFile(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("Cannot be blank")
	}
	return nil
}

func validateString(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("Cannot be blank")
	}
	return nil
}

//RequestSelect is a wrapper to a promptui selecting one from an option
func RequestSelect(promptText string, answerSet []string, defaultValue string) string {
	var result string
	var err error

	prompt := promptui.Select{
		Label: promptText + "  ",
		Items: answerSet,
		Size:  10,
	}

	scrollPos := -1
	if defaultValue != "" {
		for i, v := range answerSet {
			if v == defaultValue {
				scrollPos = i
				break
			}
		}
	}

	if scrollPos < 0 {
		_, result, err = prompt.Run()
	} else {
		_, result, err = prompt.RunCursorAt(scrollPos, 0)
	}

	if err != nil {
		Message("A Select Error: ", err)
	}
	return result
}
