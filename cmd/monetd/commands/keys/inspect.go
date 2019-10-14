package keys

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _showPrivate bool

// newInspectCmd returns the command that inspects an Ethereum keyfile
func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [moniker]",
		Short: "inspect a keyfile",
		Long: `
Display the contents of a keyfile.

The output contains the corresponding address and public key. If --private is
specified, the keyfile will be decrypted with the passphrase and the raw private
key will also be returned. If --passfile is not specified, the user will be
prompted to enter the passphrase manually.
		`,
		Args: cobra.ExactArgs(1),
		RunE: inspect,
	}

	addInspectFlags(cmd)

	return cmd
}

// addInspectFlags adds flags to the Inspect command
func addInspectFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&_showPrivate, "private", false, "include the private key in the output")
	viper.BindPFlags(cmd.Flags())
}

func inspect(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	err := crypto.InspectKeyByMoniker(_keystore, moniker, _passwordFile, _showPrivate, _outputJSON)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
