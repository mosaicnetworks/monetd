package keys

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/mosaicnetworks/monetd/src/configuration"
	monetcrypto "github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/spf13/cobra"
)

//newImportCmd returns the command that creates a Ethereum keyfile
func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [moniker] [keyfile]",
		Short: "import a private key to import a new keyfile",
		Long: `
Import keys to [moniker] from private key file [keyfile].
`,
		Args: cobra.ExactArgs(2),
		RunE: importKey,
	}

	return cmd
}

func importKey(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	privateKeyfile := args[1]

	common.DebugMessage(fmt.Sprintf("Importing to node %s from %s", moniker, privateKeyfile))

	_, err := monetcrypto.NewKeyPairFull(configuration.Configuration.DataDir, moniker, passwordFile, privateKeyfile, outputJSON)

	return err
}
