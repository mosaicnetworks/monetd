// Package common provides constants and utility functions that are shared
// across the poa packages
package common

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

var (
	//VerboseLogging is a globals that when set to false suppresses
	//Debug Message output.
	//Should be set by the --verbose option in monetcli.
	VerboseLogging = false
)

//Colour definitions as used in the Message functions in this unit.
const (
	ColourInfo    = color.FgGreen
	ColourWarning = color.FgHiMagenta
	ColourError   = color.FgHiRed
	ColourPrompt  = color.FgHiYellow
	ColourOther   = color.FgYellow
	ColourOutput  = color.FgHiCyan
	ColourDebug   = color.FgCyan
)

//Log level constants
const (
	MsgInformation = 0
	MsgWarning     = 1
	MsgError       = 2
	MsgPrompt      = 3
	MsgDebug       = 4
	MsgOther       = 5
)

//InfoMessage is a simple wrapper for stdout logging of Information Messages
func InfoMessage(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = MessageWithType(MsgInformation, a...)
		return n, err
	}
	return 0, nil
}

//ErrorMessage is a simple wrapper for stdout logging for Error Messages.
func ErrorMessage(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = MessageWithType(MsgError, a...)
		return n, err
	}
	return 0, nil
}

//DebugMessage is a simple wrapper for stdout logging. Setting VerboseLayout to false disables its output
func DebugMessage(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = MessageWithType(MsgDebug, a...)
		return n, err
	}
	return 0, nil
}

//MessageWithType is a central point for cli logging messages
//It colour codes the output, suppressing Debug messages if VerboseLogging is false
func MessageWithType(msgType int, a ...interface{}) (n int, err error) {

	color.Set(ColourOther)

	var prefix = ""

	switch msgType {
	case MsgInformation:
		color.Set(ColourInfo)
		//		prefix = "Info: "
	case MsgWarning:
		color.Set(ColourWarning)
		//		prefix = "Warn: "
	case MsgError:
		color.Set(ColourError)
		//		prefix = "Error: "
	case MsgPrompt:
		color.Set(ColourPrompt)
		//		prefix = ""
	case MsgDebug:
		if !VerboseLogging {
			return 0, nil
		}
		color.Set(ColourDebug)
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

//ClearScreen clears the CLI screen. Implementation is OS-specific
func ClearScreen() {
	// Attempt to clear cli screen.
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		return
	}

	print("\033[H\033[2J")
}
