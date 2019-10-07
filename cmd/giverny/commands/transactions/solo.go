package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type soloAccount struct {
	Moniker      string
	Address      string
	NextNonce    int
	Credits      *big.Int
	Debits       *big.Int
	Delta        *big.Int
	Transactions []soloTransaction
}

type soloTransaction struct {
	To     string
	ToName string
	Nonce  int
	Amount *big.Int
}

//CLI params
var accounts string
var outputfile = "trans.json"
var maxTransValue = 10
var roundRobin = false

func newSoloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solo",
		Short: "solo transactions",
		Long: `
Solo transactions generate a transaction set without needing access 
to the network toml file. You just need a well funded faucet account.
The additional accounts can be generated using giverny keys generate
`,
		Args: cobra.ArbitraryArgs,
		RunE: soloTransactions,
	}

	addSoloFlags(cmd)
	return cmd
}

func addSoloFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&faucet, "faucet", faucet, "faucet account moniker")
	cmd.Flags().StringVar(&accounts, "accounts", accounts, "comma separated account list")
	cmd.Flags().StringVar(&outputfile, "output", outputfile, "output file")

	cmd.Flags().BoolVar(&roundRobin, "round-robin", roundRobin, "set sender accounts round robin")

	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to solo")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	cmd.Flags().IntVar(&maxTransValue, "max-trans-value", maxTransValue, "maximum transaction value")

	viper.BindPFlags(cmd.Flags())
}

func soloTransactions(cmd *cobra.Command, args []string) error {

	// Sanity Checking

	surplusCreditBig := new(big.Int).SetInt64(int64(surplusCredit))

	common.DebugMessage("Surplus set to: ", surplusCreditBig)

	if files.CheckIfExists(outputfile) {
		return errors.New("output file already exists: " + outputfile)
	}

	if !files.CheckIfExists(_keystore) {
		return errors.New("keystore does not exist: " + _keystore)
	}

	faucetKey := filepath.Join(_keystore, faucet+".json")
	if !files.CheckIfExists(faucetKey) {
		return errors.New("faucet key does not exist: " + faucetKey)
	}

	if accounts == "" {
		return errors.New("no accounts specified")
	}

	// Build accounts, the faucet account is the first item in the list.

	arrAccounts := strings.Split(accounts, ",")
	accountCount := len(arrAccounts)

	if accountCount < 2 {
		return errors.New("you must have at least 2 accounts")
	}

	arrAccounts = append([]string{faucet}, arrAccounts...)

	var Accounts []soloAccount

	for _, acc := range arrAccounts {
		accKey := filepath.Join(_keystore, acc+".json")
		if !files.CheckIfExists(accKey) {
			return errors.New("cannot find key for account " + accKey)
		}

		keyjson, err := ioutil.ReadFile(accKey)
		if err != nil {
			return fmt.Errorf("Failed to read the keyfile at '%s': %v", accKey, err)
		}

		m := make(map[string]interface{})
		if err := json.Unmarshal(keyjson, &m); err != nil {
			return err
		}

		addr := m["address"].(string)

		common.DebugMessage(acc)
		Accounts = append(Accounts, soloAccount{
			Moniker:   acc,
			NextNonce: 1,
			Address:   addr,
			Credits:   new(big.Int),
			Debits:    new(big.Int),
			Delta:     new(big.Int),
		})
	}

	// Build Transactions

	for i := 0; i < totalTransactions; i++ {
		var fromacct, toacct int

		for {
			if roundRobin {
				fromacct = (i % accountCount) + 1
			} else {
				fromacct = rand.Intn(accountCount) + 1
			}
			toacct = rand.Intn(accountCount) + 1
			if fromacct != toacct {
				break
			}
		}

		amt := new(big.Int).SetInt64(int64(rand.Intn((maxTransValue*100)-10) + 9))
		amt.Mul(amt, new(big.Int).SetInt64(100000000))
		amt.Mul(amt, new(big.Int).SetInt64(100000000))

		// TODO add some lower order bits noise

		Accounts[fromacct].Transactions = append(Accounts[fromacct].Transactions,
			soloTransaction{
				To:     Accounts[toacct].Address,
				ToName: Accounts[toacct].Moniker,
				Nonce:  Accounts[fromacct].NextNonce,
				Amount: amt,
			})

		Accounts[fromacct].NextNonce++
		Accounts[fromacct].Debits.Add(Accounts[fromacct].Debits, amt)
		Accounts[toacct].Credits.Add(Accounts[toacct].Credits, amt)
		Accounts[fromacct].Delta.Sub(Accounts[fromacct].Delta, amt)
		Accounts[toacct].Delta.Add(Accounts[toacct].Delta, amt)

	}

	// Transactions are generated. Set the faucet transctions

	for i := 1; i < len(Accounts); i++ {
		Accounts[0].Transactions = append(Accounts[0].Transactions,
			soloTransaction{
				To:     Accounts[i].Address,
				ToName: Accounts[i].Moniker,
				Nonce:  Accounts[0].NextNonce,
				Amount: Accounts[i].Debits,
			})

		Accounts[0].NextNonce++
		Accounts[0].Debits.Add(Accounts[0].Debits, Accounts[i].Debits)
		Accounts[0].Delta.Sub(Accounts[0].Delta, Accounts[i].Debits)
	}

	// Write Output

	var jsonData []byte

	jsonData, err := json.Marshal(Accounts)
	if err != nil {
		return err
	}
	files.WriteToFile(outputfile, string(jsonData))
	common.DebugMessage("Node File written: ", outputfile)

	common.InfoMessage(string(jsonData))

	return nil
}
