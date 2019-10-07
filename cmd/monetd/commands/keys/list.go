package keys

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// newListCmd returns the command that lists keyfiles
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list keyfiles",
		Long:  "List keyfiles in <keystore>",
		RunE:  list,
	}

	return cmd
}

// list prints all the *.json files in [datadir]/keystore
func list(cmd *cobra.Command, args []string) error {

	files, err := ioutil.ReadDir(_keystore)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fmt.Println(
				strings.TrimSuffix(
					file.Name(),
					filepath.Ext(file.Name()),
				),
			)
		}
	}

	return nil
}
