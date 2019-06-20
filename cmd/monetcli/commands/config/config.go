package config

import (
	"github.com/spf13/cobra"
)

//NetworkCmd controls network configuration
var ConfigCmd = &cobra.Command{
	Use:              "config",
	Short:            "manage monetd configuration",
	TraverseChildren: true,
}

func init() {
	//Subcommands
	ConfigCmd.AddCommand(
		NewCheckCmd(),
	)

	//Commonly used command line flags
	//	NetworkCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	//	NetworkCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	//	viper.BindPFlags(NetworkCmd.Flags())
}

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check the configuration",
		Long: `
Check the configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: checkconfig,
	}
	return cmd
}

func checkconfig(cmd *cobra.Command, args []string) error {
	return nil
}
