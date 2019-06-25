package network

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// This section defines some constants used throughout this package.
const (
	defaultSolidityContract = "https://raw.githubusercontent.com/mosaicnetworks/evm-lite/poa/e2e/smart-contracts/genesis_array.sol"
	templateContract        = "template.sol"
	genesisContract         = "contract0.sol"
	genesisABI              = "contract0.abi"
	defaultAccountBalance   = "1234000000000000000000"
	genesisFileName         = "genesis.json"
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
		newNewCmd(),      // Barebones implemented. Need to add more parameters
		newCheckCmd(),    // Framework implemented. Only checks contract address.
		newAddCmd(),      // Add Peers, framework in place
		newShowCmd(),     // Complete.
		newGenerateCmd(), // Complete
		newContractCmd(), // Complete
		newCompileCmd(),  // Solidity Compile in place, need to finish peers amendments to solidity, generation of .monetd files
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

func newContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [contract]",
		Short: "set solidity contract",
		Long: `
Sets the solidity contract to use for poa.`,
		Args: cobra.ExactArgs(1),
		RunE: setContract,
	}
	return cmd
}

func setContract(cmd *cobra.Command, args []string) error {
	sol := args[0]

	if !checkIfExists(sol) {
		message("Cannot read solidity contract file: ", sol)
		return errors.New("cannot read contract file")
	}

	targetFile := filepath.Join(configDir, templateContract)

	message("Copying sol file: ", sol, targetFile)

	// Cut and paste copy files
	err := copyFileContents(sol, targetFile)

	return err
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
		Use:   "generate [moniker] [ip] [isValidator]",
		Short: "generate and add key pair",
		Long: `
Generate and add a key pair to the configuration.`,
		Args: cobra.ExactArgs(3),
		RunE: generatekeypair,
	}
	return cmd
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
