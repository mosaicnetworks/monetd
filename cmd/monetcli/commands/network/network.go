package network

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NetworkCmd controls network configuration
var (
	NetworkCmd = &cobra.Command{
		Use:   "network",
		Short: "manage monet network configuration",
		Long: `Network
		
The network subcommand is used to configure a network of hubs within the monetcli configuration. The compile option builds the genesis file and pushes it to a monetd configuration. The commands available from the network command are sequenced in the wizard, testnet and testjoin commands.`,
		TraverseChildren: true,
	}

	configDir    string
	force        bool
	passwordFile string
)

func init() {
	//Subcommands
	NetworkCmd.AddCommand(
		newNewCmd(),
		newCheckCmd(),
		NewLocationCmd(),
		newAddCmd(),
		newShowCmd(),
		newGenerateCmd(),
		newContractCmd(),
		newParamsCmd(),
		newCompileCmd(),
		newPeersCmd(),
	)

	defaultConfigDir, err := common.DefaultHomeDir(common.MonetcliTomlDir)
	if err != nil {
		fmt.Println(err.Error())
	}

	NetworkCmd.PersistentFlags().StringVarP(&configDir, "config-dir", "c", defaultConfigDir, "the directory containing the network.toml file holding the monetcli configuration")
	viper.BindPFlags(NetworkCmd.Flags())
}

func newPeersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peers",
		Short: "review peers list",
		Long: `
Interactively review the Peers list with the ability to edit / delete entries.`,
		Args: cobra.ArbitraryArgs,
		RunE: reviewPeers,
	}

	return cmd
}

//NewLocationCmd defines the CLI command config check
func NewLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `monetcli network location
Outputs the location of the configuration files for the monetcli network.`,
		Args: cobra.ArbitraryArgs,
		RunE: locationConfig,
	}
	return cmd
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

func newParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "update parameters interactively",
		Long: `
Update Parameters interactively`,
		Args: cobra.ExactArgs(0),
		RunE: setParams,
	}

	return cmd
}

//check add generate compile

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [moniker] [publickey] [ip] [isValidator]",
		Short: "add key pair",
		Long: `Network Add

Add a peer to the monetcli configuration.`,
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

	cmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	viper.BindPFlags(cmd.Flags())

	return cmd
}

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "compile configuration",
		Long: `network compile

Compile monetcli network configuration into a monet hub configuration. This includes building the solidity smart contract for proof of authority with the initial peer list baked in, compiling in and placing it in the genesis block. Additionally peers files are written and the mandatory monet configurations are applied on top of the monetcli parameters.`,
		Args: cobra.ExactArgs(0),
		RunE: compileConfig,
	}
	return cmd
}

func reviewPeers(cmd *cobra.Command, args []string) error {

	return PeersWizard(configDir)
}

func setContract(cmd *cobra.Command, args []string) error {
	sol := args[0]

	if !common.CheckIfExists(sol) {
		message("Cannot read solidity contract file: ", sol)
		return errors.New("cannot read contract file")
	}

	targetFile := filepath.Join(configDir, common.TemplateContract)

	message("Copying sol file: ", sol, targetFile)

	// Cut and paste copy files
	err := common.CopyFileContents(sol, targetFile)

	return err
}

func locationConfig(cmd *cobra.Command, args []string) error {
	common.MessageWithType(common.MsgInformation, "The Monetcli Network Configuration files are located at:")
	common.MessageWithType(common.MsgInformation, configDir)
	return nil
}
