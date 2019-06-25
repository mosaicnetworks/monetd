package commands

import "github.com/spf13/cobra"

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive configuration wizard",
	RunE:  initialise,
}

func initialise(cmd *cobra.Command, args []string) error {
	// get datadir

	// generate a new key?

	// add to peers.json?

	// add to genesis.json?

	// is validator?

	return nil
}
