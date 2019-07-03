package common

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/mbndr/figlet4go"
)

var (
	//VerboseLogging controls where Message produces any output
	VerboseLogging = true
	//HideBanners Suppresses the figlet banner
	HideBanners = false
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

//Message is a simple wrapper for stdout logging. Setting VerboseLayout to false disables its output
func Message(a ...interface{}) (n int, err error) {
	if VerboseLogging {
		n, err = MessageWithType(MsgDebug, a...)
		return n, err
	}

	return 0, nil
}

//MessageWithType is a central point for cli logging messages
//It colour codes the output, suppressing Debug messages if --verbose option not enabled
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

//Banner outputs a title.
//It can be suppressed with a -q flag
func Banner(msg string) {
	if HideBanners {
		fmt.Println(msg)
		fmt.Println("")
		return
	}
	options := figlet4go.NewRenderOptions()
	options.FontName = "larry3d"
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorRed,
		figlet4go.ColorBlue,
		figlet4go.ColorGreen,
	}

	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.RenderOpts(msg, options)
	fmt.Print(renderStr)
}

//BannerTitle wraps ClearScreen, Banner and BlankLine to show a title on a
//cleared page with a blankline underneath
func BannerTitle(msg string) {
	ClearScreen()
	Banner(msg)
	BlankLine()
}

//BlankLine is a Wrapper function for a Blank Line on the terminal
func BlankLine() {
	fmt.Println("")
}
