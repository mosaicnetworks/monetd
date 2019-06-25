package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show configuration",
		Long: `
Show configuration.`,
		Args: cobra.ExactArgs(0),
		RunE: showConfig,
	}

	return cmd
}

func showConfig(cmd *cobra.Command, args []string) error {

	filename := filepath.Join(configDir, tomlName+".toml")
	message("Displaying file: ", filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
