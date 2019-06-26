package network

import (
	"path/filepath"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

const (
	SaveChanges   = "Save Changes"
	CancelChanges = "Cancel Changes"
	EditAgain     = "Edit Again"
)

func setParams(cmd *cobra.Command, args []string) error {
	// Call a wrapper function to ease calling from outside cobra
	return setParamsWithParams()
}

func setParamsWithParams() error {
	var result string
	// We assume the default here, but all the file path stuff is overwritten in the
	// config publish actions any way so there is no lasting impact
	assumedMonetDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)
	keySet := common.GetMonetDefaultConfigKeys(assumedMonetDir)

	configFile := filepath.Join(configDir, common.MonetcliTomlName+".toml")
	tree, err := common.LoadToml(configFile)
	if err != nil {
		common.Message("Cannot load configuration toml file: ", configFile)
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
				//TODO actually implement
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
					common.Message("Error entering parameters: ", err)
					//TODO - do we return this error, or just not save the mangled result as per now.
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
				}

				result, err := prompt.Run()
				if err != nil {
					common.Message("Error entering parameters: ", err)
					//TODO - do we return this error, or just not save the mangled result as per now.
				} else {
					tree.Set(keys.Key, result)
				}
			}
		}

		prompt := promptui.Select{
			Label: "Do you want to save your changes",
			Items: []string{CancelChanges, SaveChanges, EditAgain},
		}

		_, result, err = prompt.Run()

		if err != nil {
			common.Message("Error in params confirmation: ", err)
			return err
		}

		// Infinite loop. Need to save or cancel to break out it
		switch result {
		case SaveChanges:
			common.SaveToml(tree, configFile)
			break paramloop
		case CancelChanges:
			break paramloop
		}

	}
	return nil
}
