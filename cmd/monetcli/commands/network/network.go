package network

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultSolidityContract = "https://raw.githubusercontent.com/mosaicnetworks/evm-lite/poa/e2e/smart-contracts/genesis_array.sol"
	templateContract        = "template.sol"
	genesisContract         = "contract0.sol"
	genesisABI              = "contract0.abi"
)

//NetworkCmd controls network configuration
var (
	NetworkCmd = &cobra.Command{
		Use:              "network",
		Short:            "manage monet network configuration",
		TraverseChildren: true,
	}

	configDir      string
	force          bool
	verboseLogging = true
)

func init() {
	//Subcommands
	NetworkCmd.AddCommand(
		newNewCmd(),   // Barebones implemented.
		newCheckCmd(), // Framework implemented. Only checks contract address.
		newAddCmd(),
		newShowCmd(), // Complete.
		newGenerateCmd(),
		newCompileCmd(),
	)

	var defaultConfigDir, err = defaultHomeDir()

	if err != nil {
		fmt.Println(err.Error())
	}

	NetworkCmd.PersistentFlags().StringVar(&configDir, "config-dir", defaultConfigDir, "the directory containing the network.toml file")
	NetworkCmd.PersistentFlags().BoolVar(&verboseLogging, "verbose", false, "verbose messages")

	//	NetworkCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
	viper.BindPFlags(NetworkCmd.Flags())
}

//check add generate compile

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [moniker] [address] [ip] [isValidator]",
		Short: "add key pair",
		Long: `
Add a key pair to the configuration.`,
		Args: cobra.ExactArgs(4),
		RunE: addValidator,
	}
	return cmd
}

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [moniker]",
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

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "compile configuration",
		Long: `
compile network configuration.`,
		Args: cobra.ExactArgs(0),
		RunE: compileConfig,
	}
	return cmd
}
