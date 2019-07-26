package server

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the giverny server in the background",
		Long: `
Start the giverny server as a background process.
`,
		Args: cobra.ArbitraryArgs,
		RunE: startServer,
	}

	return cmd
}

func startServer(cmd *cobra.Command, args []string) error {

	return start()

}

func start() error {

	// check if daemon already running.
	if _, err := os.Stat(pidFile); err == nil {
		return errors.New("already running or " + pidFile + " exists")
	}

	cmd := exec.Command(os.Args[0], "server", "run", "--background")
	cmd.Start()
	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
	savePID(strconv.Itoa(cmd.Process.Pid))

	return nil
}
