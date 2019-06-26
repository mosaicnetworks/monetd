package common

import "fmt"

var (
	//VerboseLogging controls where Message produces any output
	VerboseLogging = true
)

//Message is a simple wrapper for stdout logging. Setting VerboseLayout to false disables its output
func Message(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = fmt.Println(a...)
		return n, err
	}

	return 0, nil
}
