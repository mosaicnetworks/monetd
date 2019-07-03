package config

import (
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NetworkCmd controls network configuration
var (
	//ConfigCmd implements the config CLI subcommand
	ConfigCmd = &cobra.Command{
		Use:              "config",
		Short:            "manage monetd configuration",
		TraverseChildren: true,
	}

	publishTarget    string
	monetConfigDir   string
	networkConfigDir string
	Force            bool
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
	viper.BindPFlags(ConfigCmd.Flags())
}

//NewCheckCmd defines the CLI command config check
func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check the configuration",
		Long: `
Check the configuration of the monetd server to ensure that it is consistent.`,
		Args: cobra.ArbitraryArgs,
		RunE: checkConfig,
	}
	return cmd
}

//NewPublishCmd implements the "config publish" CLI subcommand
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
	cmd.Flags().BoolVarP(&Force, "force", "f", false, "force the creation of a new config file")
	viper.BindPFlags(cmd.Flags())

	return cmd
}
