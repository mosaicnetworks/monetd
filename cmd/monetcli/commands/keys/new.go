package keys

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//AddNewFlags adds flags to the New command
func AddNewFlags(cmd *cobra.Command) {

	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)
	cmd.PersistentFlags().StringVarP(&monetConfigDir, "monet-config-dir", "m", defaultMonetConfigDir, "the directory containing monet nodes configurations")
	cmd.Flags().StringVar(&monikerParam, "moniker", "", "specify moniker for this key")
	viper.BindPFlags(cmd.Flags())
}

//NewNewCmd returns the command that creates a new keypair
func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a new keypair",
		Long: `
Creates a new key pair and stores in under specified moniker`,
		Args: cobra.ArbitraryArgs,
		RunE: newkeys,
	}

	AddNewFlags(cmd)

	return cmd
}

func newkeys(cmd *cobra.Command, args []string) error {

	if strings.TrimSpace(monikerParam) == "" {
		return errors.New("You need to specify a moniker with --moniker ")
	}

	accountsDir := filepath.Join(monetConfigDir, common.MonetAccountsSubFolder)

	// Create Monet Config if missing
	if !common.CheckIfExists(monetConfigDir) {
		err := os.MkdirAll(monetConfigDir, os.ModePerm)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error creating directory: ", monetConfigDir)
			return err
		}
	}

	// Create accountsDir if missing
	if !common.CheckIfExists(accountsDir) {
		err := os.MkdirAll(accountsDir, os.ModePerm)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error creating directory: ", accountsDir)
			return err
		}
	}

	// We use a "safe" version of the node name for our accounts folder
	safeLabel := common.GetNodeSafeLabel(monikerParam)

	nodeDir := filepath.Join(accountsDir, safeLabel)

	if common.CheckIfExists(nodeDir) {
		return errors.New("that moniker has already been used")
	}

	err := os.MkdirAll(nodeDir, os.ModePerm)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error creating directory: ", nodeDir)
		return err
	}

	// Generate Key
	keyFile := filepath.Join(nodeDir, common.DefaultKeyfile)

	key, err := GenerateKeyPair(keyFile, passwordFile)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error creating key pair: ", keyFile)
		return err
	}

	fmt.Print(key.Address)

	return nil
}
