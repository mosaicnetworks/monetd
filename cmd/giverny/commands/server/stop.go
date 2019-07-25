package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the giverny server",
		Long: `
Stop the giverny server.
`,
		Args: cobra.ExactArgs(0),
		RunE: startServer,
	}

	return cmd
}

func stopServer(cmd *cobra.Command, args []string) error {

	// upon receiving the stop command
	// read the Process ID stored in PIDfile
	// kill the process using the Process ID
	// and exit. If Process ID does not exist, prompt error and quit

	if _, err := os.Stat(pidFile); err == nil {
		data, err := ioutil.ReadFile(pidFile)
		if err != nil {
			fmt.Println("Not running")
			os.Exit(1)
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			fmt.Println("Unable to read and parse process id found in ", pidFile)
			os.Exit(1)
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			os.Exit(1)
		}
		// remove PID file
		os.Remove(pidFile)

		fmt.Printf("Killing process ID [%v] now.\n", ProcessID)
		// kill process and exist immediately
		err = process.Kill()

		if err != nil {
			fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
			os.Exit(1)
		} else {
			fmt.Printf("Killed process ID [%v]\n", ProcessID)
			os.Exit(0)
		}

	} else {

		fmt.Println("Not running.")
		os.Exit(1)
	}

	return nil
}
