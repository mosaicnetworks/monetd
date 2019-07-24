package keys

import (
	"errors"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/mosaicnetworks/monetd/src/configuration"
	monetcrypto "github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//privateKeyfile is the default name of the file containing a imported private key

	privateKeyfile string
)

//addImportFlags adds flags to the Import command
func addImportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&privateKeyfile, "privatekey", "", "file containing a raw private key to encrypt")

	viper.BindPFlags(cmd.Flags())
}

//newImportCmd returns the command that creates a Ethereum keyfile
func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [moniker] [keyfile]",
		Short: "import a private key to import a new keyfile",
		Long: `
Import a new keyfile.

If you want to encrypt an existing private key, it can be specified by setting
--privatekey with the location of the file containing the private key.
`,
		Args: cobra.ExactArgs(1),
		RunE: importKey,
	}

	addImportFlags(cmd)

	return cmd
}

func importKey(cmd *cobra.Command, args []string) error {
	// Check if keyfile path given and make sure it doesn't already exist.

	moniker := args[0]
	common.InfoMessage("Priv ", privateKeyfile)
	if privateKeyfile == "" {
		return errors.New("you have not specified a file to import using the --privatekey parameter")
	}

	_, err := monetcrypto.NewKeyPairFull(configuration.Configuration.DataDir, moniker, passwordFile, privateKeyfile, outputJSON)

	return err
}
