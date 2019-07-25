package server

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the giverny server",
		Long: `
Start the giverny server.
`,
		Args: cobra.ExactArgs(0),
		RunE: startServer,
	}

	return cmd
}

func startServer(cmd *cobra.Command, args []string) error {

	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("Already running or " + pidFile + " file exist.")
		os.Exit(1)
	}

	daemonize("/dev/null", logOut, logErr)

	return nil
}
