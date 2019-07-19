package keys

import (
	"encoding/json"
	"fmt"
)

// mustPrintJSON prints the JSON encoding of the given object and
// exits the program with an error message when the marshaling fails.
func mustPrintJSON(jsonObject interface{}) error {
	str, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON object: %v", err)
	}
	fmt.Println(string(str))
	return nil
}
