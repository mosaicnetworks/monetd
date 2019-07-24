package commands

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"
	"github.com/mosaicnetworks/monetd/src/poa/files"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mconfig "github.com/mosaicnetworks/monetd/src/configuration"
	pconfig "github.com/mosaicnetworks/monetd/src/poa/config"
)

var (
	//ConfigCmd implements the config CLI subcommand
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "manage monetd configuration",
		Long: `monetd config
		
The config subcommands manage the monet configuration file, as used by 
the monetd server process. `,
		TraverseChildren: true,
	}
)

func init() {
	//Subcommands
	configCmd.AddCommand(
		//			NewCheckCmd(),
		//			NewPublishCmd(),
		newLocationCmd(),
		//			NewShowCmd(),
		newClearCmd(),
		newPullCmd(),
		newBuildCmd(),
	)

	//TODO remove - temporary debug out to preserve the import - we will need it shortly
	fmt.Print(mconfig.Configuration.DataDir)

	// datadir is now the config for everything...

	//	viper.BindPFlags(ConfigCmd.Flags())
}

// newLocationCmd shows the config file path
func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `monetd config location
Shows the location of the configuration files for the monetd server.`,
		Args: cobra.ArbitraryArgs,
		RunE: locationConfig,
	}
	return cmd
}

func locationConfig(cmd *cobra.Command, args []string) error {
	common.InfoMessage("The Monet Configuration files are located at:")
	common.InfoMessage(mconfig.Configuration.DataDir)
	return nil
}

//NewBuildCmd echoes the config file to screen
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build the configuration files",
		Long: `monetd config build
Builds the monetd configuration files for the monetd server.`,
		Args: cobra.ArbitraryArgs,
		RunE: buildConfig,
	}

	cmd.PersistentFlags().StringVarP(&nodeParam, "node", "n", "", "the directory name containing monet nodes configurations")
	cmd.PersistentFlags().StringVarP(&addressParam, "address", "a", "", " ip address/host name of this node")
	cmd.PersistentFlags().StringVarP(&passwordFile, "passfile", "p", "", "the file that contains the passphrase for the keyfile")
	//	KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")

	viper.BindPFlags(cmd.Flags())

	//--node node0  --address 192.168.1.4 --peers node1,node2,node3 --peer-address host1,host2,host3

	return cmd
}

//newClearCmd shows the config file path
func newClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "clears any pre-existing configuration files",
		Long: `monetd config clear
Cleans up configuration files for the monetd server by renaming any pre-existing configuration.

Clearly this will disable any pre-existing configuration.`,
		Args: cobra.ArbitraryArgs,
		RunE: clearConfig,
	}
	return cmd
}

func newPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "pull the configuration files from a node",
		Long: `monetd config pull
Pulls the monetd configuration files from an existing peer.`,
		Args: cobra.ArbitraryArgs,
		RunE: pullConfig,
	}

	cmd.PersistentFlags().StringVarP(&nodeParam, "node", "n", "", "the directory name containing monet nodes configurations")
	cmd.PersistentFlags().StringVarP(&addressParam, "address", "a", "", " ip address/host name of this node")
	cmd.PersistentFlags().StringVarP(&passwordFile, "passfile", "p", "", "the file that contains the passphrase for the keyfile")
	cmd.PersistentFlags().StringVar(&existingPeer, "peer", "", "the address of an existing peer")

	//	KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func buildConfig(cmd *cobra.Command, args []string) error {
	return pconfig.BuildConfig(mconfig.Configuration.DataDir, nodeParam, addressParam, passwordFile)
}

func clearConfig(cmd *cobra.Command, args []string) error {

	if files.CheckIfExists(mconfig.Configuration.DataDir) {
		files.SafeRenameDir(mconfig.Configuration.DataDir)
	}

	//	ShowConfigParams(monetConfigDir)
	return nil
}

func pullConfig(cmd *cobra.Command, args []string) error {
	return nil
}
