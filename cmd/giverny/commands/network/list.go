package network

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/spf13/cobra"
)

// newListCmd returns the command that lists keyfiles
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list configured networks",
		RunE:  list,
	}

	return cmd
}

// list prints the names of all the folders in [datadir]/networks
func list(cmd *cobra.Command, args []string) error {

	networksDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir)

	files, err := ioutil.ReadDir(networksDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}
