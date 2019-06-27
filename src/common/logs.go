package common

import "fmt"

var (
	//VerboseLogging controls where Message produces any output
	VerboseLogging = true
)

const (
	MsgInformation = 0
	MsgWarning     = 1
	MsgError       = 2
	MsgPrompt      = 3
	MsgDebug       = 4
)

//Message is a simple wrapper for stdout logging. Setting VerboseLayout to false disables its output
func Message(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = fmt.Println(a...)
		return n, err
	}

	return 0, nil
}

//TODO Change the prefix to be colour codes
func MessageWithType(msgType int, a ...interface{}) (n int, err error) {
	var prefix = ""

	switch msgType {
	case MsgInformation:
		prefix = "Info: "
	case MsgWarning:
		prefix = "Warn: "
	case MsgError:
		prefix = "Error: "
	case MsgPrompt:
		prefix = ""
	case MsgDebug:
		prefix = "Debug: "
	}

	n, err = fmt.Println(append([]interface{}{prefix}, a...))
	return n, err
}
