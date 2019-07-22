package keys

import (
	"github.com/mosaicnetworks/monetd/cmd/monetd/config"
	"github.com/mosaicnetworks/monetd/src/poa/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Global CLI parameters
var (
	PasswordFile      string
	OutputJSON        bool
	monikerParam      string
	monetCliConfigDir string
	newPasswordFile   string
	showPrivate       = false
)

//KeysCmd is an Ethereum key manager
var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "monet key manager",
	Long: `Keys
	
Monet Key Manager. `,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	KeysCmd.AddCommand(
		//		NewGenerateCmd(),
		newInspectCmd(),
		newUpdateCmd(),
		newNewCmd(),
	)

	//Commonly used command line flags
	KeysCmd.PersistentFlags().StringVar(&PasswordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
	KeysCmd.PersistentFlags().BoolVar(&OutputJSON, "json", false, "output JSON instead of human-readable format")
	//	KeysCmd.PersistentFlags().StringVar(&monikerParam, "moniker", "", "specify moniker for this key")

	viper.BindPFlags(KeysCmd.Flags())
}

//addInspectFlags adds flags to the Inspect command
func addInspectFlags(cmd *cobra.Command) {

	cmd.Flags().BoolVar(&showPrivate, "private", false, "include the private key in the output")
	viper.BindPFlags(cmd.Flags())
}

//NewInspectCmd returns the command that inspects an Ethereum keyfile
func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [moniker]",
		Short: "inspect a keyfile",
		Long: `
Print various information about the keyfile.

Private key information can be printed by using the --private flag;
make sure to use this feature with great caution!

You must specify the moniker for this node with the --moniker parameter.`,
		Args: cobra.ExactArgs(1),
		RunE: inspect,
	}

	addInspectFlags(cmd)

	return cmd
}

func inspect(cmd *cobra.Command, args []string) error {
	monikerParam = args[0]
	return crypto.InspectKeyMoniker(config.Config.DataDir, monikerParam, PasswordFile, showPrivate, OutputJSON)
}

//NewNewCmd returns the command that creates a new keypair
func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [moniker]",
		Short: "create a new keypair",
		Long: `
Creates a new key pair and stores in under specified moniker`,
		Args: cobra.ExactArgs(1),
		RunE: newkeys,
	}

	return cmd
}

func newkeys(cmd *cobra.Command, args []string) error {
	monikerParam = args[0]

	// key is returned, but we don't want to do anything with it.
	_, err := crypto.NewKeyPair(config.Config.DataDir, monikerParam, PasswordFile)

	return err
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&newPasswordFile, "new-passfile", "", "the file containing the new passphrase for the keyfile")
	viper.BindPFlags(cmd.Flags())
}

//NewUpdateCmd returns the command that changes the passphrase of a keyfile
func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [moniker]",
		Short: "change the passphrase on a keyfile",
		Long: `keys update
Update the passphrase for the keyfile.

Private key information can be printed by using the --private flag;
make sure to use this feature with great caution!

You must specify the moniker for this node with the --moniker parameter.`,
		Args: cobra.ExactArgs(1),
		RunE: update,
	}

	addUpdateFlags(cmd)

	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	monikerParam = args[0]

	return crypto.UpdateKeysMoniker(config.Config.DataDir, monikerParam, PasswordFile, newPasswordFile)
}
