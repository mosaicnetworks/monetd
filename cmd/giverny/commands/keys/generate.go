package keys

import (
	"path/filepath"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/files"

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
		Short: "generate multiple keys",
		Long: `
The generate sub command is intended only for tests. It generates a number of
keys and writes them to <keystore>. The keyfiles are named <prefix><suffix> 
where prefix is set by --prefix (default "Account") and suffix is a number 
between --min-suffix and --max-suffix inclusive.
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

	for i := minSuffix; i <= maxSuffix; i++ {
		moniker := prefix + strconv.Itoa(i)
		common.DebugMessage("Account ", moniker)

		// check if the keyfile already exists
		keyfile := filepath.Join(_keystore, moniker+".json")
		if files.CheckIfExists(keyfile) {
			common.ErrorMessage(moniker + " already exists. skipping ")
			continue
		}

		// write the password file
		pwdfile := filepath.Join(_keystore, moniker+".txt")
		if err := files.WriteToFilePrivate(pwdfile, defaultPassword); err != nil {
			common.ErrorMessage(err)
			continue
		}

		// write the keyfile
		_, err := crypto.NewKeyfile(_keystore, moniker, pwdfile)
		if err != nil {
			common.ErrorMessage(err)
		}
	}

	return nil
}
