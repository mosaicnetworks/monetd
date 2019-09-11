package keys

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/crypto"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//CLI Params
var prefix = "Account"
var minSuffix = 1
var maxSuffix = 5

const defaultPassword = "test"

//newGenerateCmd returns the command that creates a Ethereum keyfile
func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "bulk generate key pairs",
		Long: `
The generate sub command is intended only for test nets. It generates a 
number of key pairs and places them in the current monet keystore. The 
accounts are names <prefix><suffix> where prefix is set by --prefix (default
"Account") and suffix is a number between --min-suffix and --max-suffix 
inclusive. The defaults are 1 and 5.
`,
		Args: cobra.ArbitraryArgs,
		RunE: generateKey,
	}

	cmd.Flags().StringVar(&prefix, "prefix", prefix, "prefix for account monikers")
	cmd.Flags().IntVar(&minSuffix, "min-suffix", minSuffix, "minimum suffix for account monikers")
	cmd.Flags().IntVar(&maxSuffix, "max-suffix", maxSuffix, "maximum suffix for account monikers")

	viper.BindPFlags(cmd.Flags())

	return cmd
}

func generateKey(cmd *cobra.Command, args []string) error {

	common.DebugMessage("Config dir: ", configuration.Global.DataDir)
	keydir := filepath.Join(configuration.Global.DataDir, configuration.KeyStoreDir)

	for i := minSuffix; i <= maxSuffix; i++ {
		moniker := prefix + strconv.Itoa(i)
		common.DebugMessage("Account ", moniker)
		pwdfile := filepath.Join(keydir, moniker+".txt")
		pairfile := filepath.Join(keydir, moniker+".json")
		if files.CheckIfExists(pairfile) {
			common.ErrorMessage(moniker + " already exists. skipping ")
			continue
		}

		files.WriteToFilePrivate(pwdfile, defaultPassword)
		_, err := crypto.NewKeyPair(configuration.Global.DataDir, moniker, pwdfile)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
