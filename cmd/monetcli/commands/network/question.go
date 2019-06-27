package network

import (
	"errors"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/monetd/src/common"
)

func requestFile(promptText string, defaultValue string) string {

	prompt := promptui.Prompt{
		Label:    promptText + "  ",
		Validate: validateFile,
		Default:  defaultValue,
	}

	result, _ := prompt.Run()

	return result
}

func requestString(promptText string, defaultValue string) string {

	prompt := promptui.Prompt{
		Label:    promptText + "  ",
		Validate: validateString,
		Default:  defaultValue,
	}

	result, _ := prompt.Run()

	return result
}

func requestBool(promptText string, defaultValue bool) bool {

	defaultStr := strconv.FormatBool(defaultValue)
	strVal := RequestSelect(promptText, []string{"true", "false"}, defaultStr)

	rtnVal, _ := strconv.ParseBool(strVal)

	return rtnVal
}

func validateFile(input string) error {
	//TODO actually implement

	return nil
}

func validateString(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("Cannot be blank")
	}
	return nil
}

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
		common.Message("A Select Error: ", err)
	}
	return result
}
