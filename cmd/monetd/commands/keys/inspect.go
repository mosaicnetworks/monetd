package keys

import (
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showPrivate bool

// newInspectCmd returns the command that inspects an Ethereum keyfile
func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [moniker]",
		Short: "inspect a keyfile",
		Args:  cobra.ExactArgs(1),
		RunE:  inspect,
	}

	addInspectFlags(cmd)

	return cmd
}

// addInspectFlags adds flags to the Inspect command
func addInspectFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&showPrivate, "private", false, "include the private key in the output")
	viper.BindPFlags(cmd.Flags())
}

func inspect(cmd *cobra.Command, args []string) error {
	moniker := args[0]
	return crypto.InspectKeyMoniker(configuration.Configuration.DataDir, moniker, PasswordFile, showPrivate, OutputJSON)
}
