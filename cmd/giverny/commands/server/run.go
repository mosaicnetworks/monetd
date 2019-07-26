package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inBackground = false

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run the giverny server",
		Long: `
Run the giverny server.
`,
		Args: cobra.ArbitraryArgs,
		RunE: runServer,
	}

	addRunFlags(cmd)
	return cmd
}

func addRunFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&inBackground, "background", "b", inBackground, "flag to launch server in the background")
	viper.BindPFlags(cmd.Flags())
}

func runServer(cmd *cobra.Command, args []string) error {
	return maincmd()
}

func maincmd() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	if inBackground {
		f, err := os.OpenFile(logOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
