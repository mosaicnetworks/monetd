package server

import "github.com/spf13/cobra"

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status of the giverny server",
		Long: `
Report status of the giverny server.
`,
		Args: cobra.ExactArgs(0),
		RunE: statusServer,
	}

	return cmd
}

func statusServer(cmd *cobra.Command, args []string) error {
	return nil
}
