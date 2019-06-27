package common

import (
	"fmt"

	"github.com/fatih/color"
)

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
		n, err = MessageWithType(MsgDebug, a...)
		return n, err
	}

	return 0, nil
}

//TODO Change the prefix to be colour codes
func MessageWithType(msgType int, a ...interface{}) (n int, err error) {

	color.Set(color.FgYellow)

	var prefix = ""

	switch msgType {
	case MsgInformation:
		color.Set(color.FgGreen)
		//		prefix = "Info: "
	case MsgWarning:
		color.Set(color.FgHiMagenta)
		//		prefix = "Warn: "
	case MsgError:
		color.Set(color.FgHiRed)
		//		prefix = "Error: "
	case MsgPrompt:
		color.Set(color.FgHiYellow)
		//		prefix = ""
	case MsgDebug:
		if !VerboseLogging {
			return 0, nil
		}
		color.Set(color.FgCyan)
		//		prefix = "Debug: "
	}

	if prefix == "" {
		n, err = fmt.Println(a...)

	} else {
		n, err = fmt.Println(append([]interface{}{prefix}, a...))
	}
	color.Unset()

	return n, err
}
