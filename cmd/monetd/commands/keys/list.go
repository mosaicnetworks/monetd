package keys

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

// newListCmd returns the command that lists keyfiles
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list keyfiles",
		Long: `
The list command supplies a list of moniker for the keys in the keystore 
subfolder of the configuration folder. 

The monikers are in safe format where any character not matching [0-9A-Za-z]
is converted to an underscore. `,
		RunE: list,
	}

	return cmd
}

// list prints all the *.json files in [datadir]/keystore
func list(cmd *cobra.Command, args []string) error {

	keystore := filepath.Join(configuration.Configuration.DataDir, "keystore")

	files, err := ioutil.ReadDir(keystore)
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
