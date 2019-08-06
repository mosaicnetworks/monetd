package keys

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newPasswordFile string

// newUpdateCmd returns the command that changes the passphrase of a keyfile
func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [moniker]",
		Short: "change the passphrase on a keyfile",
		Args:  cobra.ExactArgs(1),
		RunE:  update,
	}

	addUpdateFlags(cmd)

	return cmd
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&newPasswordFile, "new-passfile", "", "the file containing the new passphrase")
	viper.BindPFlags(cmd.Flags())
}

func update(cmd *cobra.Command, args []string) error {
	moniker := args[0]

	err := crypto.UpdateKeysMoniker(configuration.Global.DataDir, moniker, PasswordFile, newPasswordFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
