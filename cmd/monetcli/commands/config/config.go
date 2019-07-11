package config

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//ConfigCmd implements the config CLI subcommand
	ConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "manage monetd configuration",
		Long: `monetcli config
		
The config subcommands manage the monet configuration file, as used by 
the monetd server process. `,
		TraverseChildren: true,
	}

	nodeName         string
	monetConfigDir   string
	networkConfigDir string

	//Force is a flag to allow overwriting of config files without warning. Can be
	//set programmatically or with --force flag
	Force bool
)

func init() {
	//Subcommands
	ConfigCmd.AddCommand(
		NewCheckCmd(),
		NewPublishCmd(),
		NewLocationCmd(),
		NewShowCmd(),
	)

	defaultConfigDir, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)

	ConfigCmd.PersistentFlags().StringVarP(&monetConfigDir, "monet-config-dir", "m", defaultMonetConfigDir, "the directory containing monet nodes configurations")
	ConfigCmd.PersistentFlags().StringVarP(&networkConfigDir, "config-dir", "c", defaultConfigDir, "the directory containing monet nodes configurations")
	viper.BindPFlags(ConfigCmd.Flags())
}

//NewLocationCmd shows the config file path
func NewLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `monetcli config location
Shows the location of the configuration files for the monetd server.`,
		Args: cobra.ArbitraryArgs,
		RunE: locationConfig,
	}
	return cmd
}

//NewShowCmd echoes the config file to screen
func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show the configuration files",
		Long: `monetcli config location
Shows the monetd configuration files for the monetd server.`,
		Args: cobra.ArbitraryArgs,
		RunE: showConfig,
	}
	return cmd
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

	cmd.PersistentFlags().StringVarP(&nodeName, "node-name", "n", "", "the node to publish")
	cmd.Flags().BoolVarP(&Force, "force", "f", false, "force the creation of a new config file")
	viper.BindPFlags(cmd.Flags())

	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {
	common.MessageWithType(common.MsgInformation, "The Monet Configuration files are located at:")
	common.MessageWithType(common.MsgInformation, monetConfigDir)
	return nil
}

func showConfig(cmd *cobra.Command, args []string) error {
	ShowConfigParams(monetConfigDir)
	return nil
}

//ShowConfigParams outputs the monetd configuration file for monetConfigDir
func ShowConfigParams(monetConfigDir string) error {
	filename := filepath.Join(monetConfigDir, common.MonetdTomlName+common.TomlSuffix)
	common.MessageWithType(common.MsgInformation, "Displaying file: ", filename)
	return common.ShowConfigFile(filename)
}
