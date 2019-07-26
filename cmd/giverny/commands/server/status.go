package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status of the giverny server",
		Long: `
Report status of the giverny server.
`,
		Args: cobra.ExactArgs(0),
		RunE: statusServer,
	}

	return cmd
}

func statusServer(cmd *cobra.Command, args []string) error {

	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("PID file exists\nServer should be running")
	} else {
		fmt.Println("PID file does not exist\nServer should not be running")
	}

	url := "http://localhost:8088/ispublished"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Server is not responding")
		return nil
	}
	defer resp.Body.Close()

	fmt.Println("Server is responding")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	fmt.Printf("Is Published Status is %s\n", body)

	return nil
}
