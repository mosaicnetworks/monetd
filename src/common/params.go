package common

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

//SetParamsWithParams allows the interactive editing of parameters in the cli config file
func SetParamsWithParams(configDir string) error {
	var result string
	// We assume the default here, but all the file path stuff is overwritten in the
	// config publish actions any way so there is no lasting impact
	assumedMonetDir, _ := DefaultHomeDir(MonetdTomlDir)
	keySet := GetMonetDefaultConfigKeys(assumedMonetDir)
	configFile := filepath.Join(configDir, MonetcliTomlName+TomlSuffix)
	tree, err := LoadToml(configFile)
	if err != nil {
		Message("Cannot load configuration toml file: ", configFile)
		return err
	}

paramloop:
	for {

		// Iterate through each parameter and change values accordingly
		for _, keys := range keySet {
			if keys.Override {
				continue
			}

			defaultAnswer, ok := keys.Value.(string)
			if !ok {
				boolAnswer, ok := keys.Value.(bool)
				if ok {
					defaultAnswer = strconv.FormatBool(boolAnswer)
				}
			}

			if tree.Has(keys.Key) {
				defaultAnswer, ok = tree.Get(keys.Key).(string)
				if !ok {
					boolAnswer, ok := tree.Get(keys.Key).(bool)
					if ok {
						defaultAnswer = strconv.FormatBool(boolAnswer)
					}
				}
			}
			validateFunc := func(input string) error {
				if strings.TrimSpace(input) == "" {
					return errors.New("Cannot be blank")
				}
				return nil
			}

			promptText := keys.Key
			if keys.Prompt != "" {
				promptText = keys.Prompt
			}

			if keys.Answers != nil {

				prompt := promptui.Select{
					Label: promptText + "  ",
					Items: keys.Answers,
					Size:  10,
				}

				scrollPos := -1

				for i, v := range keys.Answers {
					if v == defaultAnswer {
						scrollPos = i
						break
					}
				}

				if scrollPos < 0 {
					_, result, err = prompt.Run()
				} else {
					_, result, err = prompt.RunCursorAt(scrollPos, 0)
				}

				if err != nil {
					Message("Error entering parameters: ", err)
					//We don't return this error, and just not save the mangled result
				} else {
					if keys.DataType == "bool" {
						boolResult, err := strconv.ParseBool(result)
						if err == nil {
							tree.Set(keys.Key, boolResult)
						}

					} else {
						tree.Set(keys.Key, result)
					}

				}
			} else {
				prompt := promptui.Prompt{
					Label:    promptText + "  ",
					Validate: validateFunc,
					Default:  defaultAnswer,
					Pointer:  promptui.PipeCursor,
				}

				result, err := prompt.Run()
				if err != nil {
					Message("Error entering parameters: ", err)
					//NB, the error is not returned, the mangled data just is not saved.
				} else {
					tree.Set(keys.Key, result)
				}
			}
		}

		prompt := promptui.Select{
			Label: "Do you want to save your changes",
			Items: []string{WizardCancelChanges, WizardSaveChanges, WizardEditAgain},
		}

		_, result, err = prompt.Run()

		if err != nil {
			Message("Error in params confirmation: ", err)
			return err
		}

		// Infinite loop. Need to save or cancel to break out it
		switch result {
		case WizardSaveChanges:
			SaveToml(tree, configFile)
			break paramloop
		case WizardCancelChanges:
			break paramloop
		}

	}
	return nil
}
