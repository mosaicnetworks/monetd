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
	"github.com/spf13/viper"
)

var inBackground = false
var logPath string

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

	addRunFlags(cmd)
	return cmd
}

func addRunFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&inBackground, "background", "b", inBackground, "flag to launch server in the background")
	cmd.Flags().StringVarP(&logPath, "logpath", "l", "", "path to log files")

	viper.BindPFlags(cmd.Flags())
}

func startServer(cmd *cobra.Command, args []string) error {

	if inBackground {
		return start()
	}

	return maincmd()
}

func start() error {

	// check if daemon already running.
	if _, err := os.Stat(pidFile); err == nil {
		return errors.New("already running or " + pidFile + " exists")
	}

	fmt.Println("Logging to : ", logOut)
	cmd := exec.Command(os.Args[0], "server", "start", "--logpath", logOut)
	cmd.Start()
	fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
	savePID(strconv.Itoa(cmd.Process.Pid))

	return nil
}

func maincmd() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

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
