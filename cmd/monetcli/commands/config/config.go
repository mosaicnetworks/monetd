package config

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NetworkCmd controls network configuration
var (
	ConfigCmd = &cobra.Command{
		Use:              "config",
		Short:            "manage monetd configuration",
		TraverseChildren: true,
	}

	publishTarget    string
	monetConfigDir   string
	networkConfigDir string
	force            bool
)

func init() {
	//Subcommands
	ConfigCmd.AddCommand(
		NewCheckCmd(),
		NewPublishCmd(),
	)

	defaultConfigDir, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)

	ConfigCmd.PersistentFlags().StringVarP(&monetConfigDir, "monet-config-dir", "m", defaultMonetConfigDir, "the directory containing monet nodes configurations")
	ConfigCmd.PersistentFlags().StringVarP(&networkConfigDir, "config-dir", "c", defaultConfigDir, "the directory containing monet nodes configurations")
	//Commonly used command line flags
	//	NetworkCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	//	NetworkCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(ConfigCmd.Flags())
}

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check the configuration",
		Long: `
Check the configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: checkConfig,
	}
	return cmd
}

func NewPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish a monet node configuration from the monetcli configuration",
		Long: `
Publish a Monet Node configuration`,
		Args: cobra.ArbitraryArgs,
		RunE: publishConfig,
	}

	cmd.PersistentFlags().StringVarP(&publishTarget, "publish-target", "t", "simple", "the publish target. One of simple, ...")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force the creation of a new config file")
	viper.BindPFlags(cmd.Flags())

	return cmd
}
