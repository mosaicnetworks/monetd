package network

import "github.com/spf13/cobra"

func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "generate new configuration",
		Long: `
Create a new configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: newconfig,
	}
	return cmd
}

func newconfig(cmd *cobra.Command, args []string) error {
	return nil
}
