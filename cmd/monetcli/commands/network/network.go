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
		Use:              "network",
		Short:            "manage monet network configuration",
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
Review the Peers list.`,
		Args: cobra.ExactArgs(0),
		RunE: reviewPeers,
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

//check add generate compile

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [moniker] [publickey] [ip] [isValidator]",
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

	cmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	viper.BindPFlags(cmd.Flags())

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
