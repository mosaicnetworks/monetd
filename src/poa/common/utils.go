package common

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

// MustPrintJSON prints the JSON encoding of the given object and
// exits the program with an error message when the marshaling fails.
func MustPrintJSON(jsonObject interface{}) error {
	str, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON object: %v", err)
	}
	fmt.Println(string(str))
	return nil
}

//GetNodeSafeLabel converts a free format string into a node label friendly format
//Anything other than an alphanumeric is converted to _
func GetNodeSafeLabel(moniker string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(moniker, "_")
}
