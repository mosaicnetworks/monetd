package commands

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/poa/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mconfig "github.com/mosaicnetworks/monetd/cmd/monetd/config"
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
		//			NewClearCmd(),
		//			NewPullCmd(),
		newBuildCmd(),
	)

	//TODO remove - temporary debug out to preserve the import - we will need it shortly
	fmt.Print(mconfig.Config.DataDir)

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
	common.InfoMessage(mconfig.Config.DataDir)
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

func buildConfig(cmd *cobra.Command, args []string) error {
	return pconfig.BuildConfig(mconfig.Config.DataDir, nodeParam, addressParam, passwordFile)
}
