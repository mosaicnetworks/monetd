package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the giverny server",
		Long: `
Start the giverny server.
`,
		Args: cobra.ArbitraryArgs,
		RunE: startServer,
	}

	return cmd
}

func startServer(cmd *cobra.Command, args []string) error {

	// this is a hidden parameter. Whilst not ideal it is a least bad pragmatic
	// choice. The regular start command performs checks for the pid file. It then
	// launches a background process of "giverny server start main". This actually
	// starts the server. We did not want a method of bypassing the already running
	// checks in the documentation.
	if len(args) > 0 && args[0] == "main" {
		return maincmd()
	}

	return start()
}

func start() error {

	// check if daemon already running.
	if _, err := os.Stat(pidFile); err == nil {
		return errors.New("already running or " + pidFile + " exists")
	}

	cmd := exec.Command(os.Args[0], "server", "start", "main")
	cmd.Start()
	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
	savePID(strconv.Itoa(cmd.Process.Pid))

	return nil
}

func maincmd() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	f, err := os.OpenFile(logOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	go func() {
		signalType := <-ch
		signal.Stop(ch)
		log.Println("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		log.Println("Received signal type : ", signalType)

		// remove PID file
		os.Remove(pidFile)

		os.Exit(0)

	}()

	servermain()

	return nil
}
