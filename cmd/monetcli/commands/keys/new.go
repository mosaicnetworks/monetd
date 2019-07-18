package keys

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//AddNewFlags adds flags to the New command
func AddNewFlags(cmd *cobra.Command) {

	defaultMonetcliConfigDir, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	cmd.PersistentFlags().StringVarP(&monetCliConfigDir, "monetcli-config-dir", "m", defaultMonetcliConfigDir, "the directory containing monet nodes configurations")
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

	accountsDir := filepath.Join(monetCliConfigDir, common.MonetAccountsSubFolder)

	// Create Monet Config if missing
	if !common.CheckIfExists(monetCliConfigDir) {
		err := os.MkdirAll(monetCliConfigDir, os.ModePerm)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error creating directory: ", monetCliConfigDir)
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
	nodeFile := filepath.Join(nodeDir, common.NodeFile)

	key, err := GenerateKeyPair(keyFile, passwordFile)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error creating key pair: ", keyFile)
		return err
	}

	common.MessageWithType(common.MsgDebug, "Generated Address      : ", key.Address.Hex())
	common.MessageWithType(common.MsgDebug, "Generated PubKey       : ", hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)))
	common.MessageWithType(common.MsgDebug, "Generated ID           : ", key.Id)

	//TODO Remove this line
	common.MessageWithType(common.MsgDebug, "Generated Private Key  : ", hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))

	// write peers config
	tree, err := toml.Load("")
	tree.SetPath([]string{"node", "moniker"}, monikerParam)
	tree.SetPath([]string{"node", "label"}, safeLabel)
	tree.SetPath([]string{"node", "address"}, key.Address.Hex())
	tree.SetPath([]string{"node", "pubkey"}, hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)))

	//TODO Remove this line
	tree.SetPath([]string{"node", "privatekey"}, hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))

	common.SaveToml(tree, nodeFile)

	return nil
}
