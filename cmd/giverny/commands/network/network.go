package network

import (
	"path/filepath"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	givernyNetworksDir  = "networks"
	givernyKeystoreDir  = "keystore"
	givernyDockerDir    = "docker"
	givernyTmpDir       = ".tmp"
	defaultTokens       = "1234567890000000000000"
	networkTomlFileName = "network.toml"
)

var (
	numberOfNodes = -1
	networkName   = "network0"
)

//NetworkCmd is the CLI subcommand
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Advanced Network Configuration",
	Long: `Network
	
Advanced Network Config Manager. `,

	TraverseChildren: true,
}

func init() {

	//Subcommands
	NetworkCmd.AddCommand(
		newAddCmd(),
		newBuildCmd(),
		newDumpCmd(),
		newExportCmd(),
		newImportCmd(),
		newNewCmd(),
		newPushCmd(),
		newStartCmd(),
		newStatusCmd(),
		newStopCmd(),
		newLocationCmd(),
		newListCmd(),
	)

	viper.BindPFlags(NetworkCmd.Flags())

	// make sure the giverny config folders exist.
	createGivernyRootNetworkFolders()

}

func createGivernyRootNetworkFolders() error {

	files.CreateDirsIfNotExists([]string{
		configuration.GivernyConfigDir,
		filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir),
	})

	return nil
}
