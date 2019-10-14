package keys

import (
	monetcrypto "github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _passwordFile string

//newImportCmd returns the command that creates a Ethereum keyfile
func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [moniker] [keyfile]",
		Short: "create an encrypted keyfile from an existing private key",
		Args:  cobra.ExactArgs(2),
		RunE:  importKey,
	}

	cmd.Flags().StringVar(&_passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func importKey(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	privateKeyfile := args[1]

	_, err := monetcrypto.NewKeyfileFull(_keystore, moniker, _passwordFile, privateKeyfile, false)

	return err
}
