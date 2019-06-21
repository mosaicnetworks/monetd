package network

import (
	"github.com/spf13/cobra"
)

//NetworkCmd controls network configuration
var NetworkCmd = &cobra.Command{
	Use:              "network",
	Short:            "manage monet network configuration",
	TraverseChildren: true,
}

func init() {
	//Subcommands
	NetworkCmd.AddCommand(
		NewNewCmd(),
		NewCheckCmd(),
		NewAddCmd(),
		NewGenerateCmd(),
		NewCompileCmd(),
	)

	//Commonly used command line flags
	//	NetworkCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	//	NetworkCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	//	viper.BindPFlags(NetworkCmd.Flags())
}

//check add generate compile

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check configuration",
		Long: `
Check configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: checkconfig,
	}
	return cmd
}

func checkconfig(cmd *cobra.Command, args []string) error {
	return nil
}

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add key pair",
		Long: `
Add a key pair to the configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: addkeypair,
	}
	return cmd
}

func addkeypair(cmd *cobra.Command, args []string) error {
	return nil
}

func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate and add key pair",
		Long: `
Generate and add a key pair to the configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: generatekeypair,
	}
	return cmd
}

func generatekeypair(cmd *cobra.Command, args []string) error {
	return nil
}

func NewCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "compile configuration",
		Long: `
compile network configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: compileconfig,
	}
	return cmd
}

func compileconfig(cmd *cobra.Command, args []string) error {
	return nil
}
