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
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
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

	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to solo")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	cmd.Flags().IntVar(&maxTransValue, "max-trans-value", maxTransValue, "maximum transaction value")
	//	cmd.Flags().StringVarP(&networkName, "network", "n", "", "network name")
	//	cmd.Flags().StringVar(&ips, "ips", "", "ips.dat file path")
	viper.BindPFlags(cmd.Flags())
}

func soloTransactions(cmd *cobra.Command, args []string) error {

	// Sanity Checking

	surplusCreditBig := new(big.Int).SetInt64(int64(surplusCredit))

	common.DebugMessage("Surplus set to: ", surplusCreditBig)

	configDir := monetconfig.Global.DataDir
	keystore := filepath.Join(configDir, monetconfig.KeyStoreDir)

	if files.CheckIfExists(outputfile) {
		return errors.New("output file already exists: " + outputfile)
	}

	if !files.CheckIfExists(keystore) {
		return errors.New("keystore does not exist: " + keystore)
	}

	faucetKey := filepath.Join(keystore, faucet+".json")
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
		accKey := filepath.Join(keystore, acc+".json")
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
			fromacct = rand.Intn(accountCount) + 1
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

	/*
			networkTomlFile := filepath.Join(networkDir, networkTomlFileName)
			if !files.CheckIfExists(networkTomlFile) {
				return errors.New("toml file does not exist: " + networkTomlFile)
			}

			transDir := filepath.Join(networkDir, givernyTransactionsDir)
			if err := files.SafeRenameDir(transDir); err != nil {
				return err
			}

			if err := files.CreateDirsIfNotExists([]string{transDir}); err != nil {
				return err
			}

			faucetTransFile := filepath.Join(transDir, "faucet.json")
			transFile := filepath.Join(transDir, "trans.json")
			deltaFile := filepath.Join(transDir, "delta.json")

			var conf = network.Config{}

			tomlbytes, err := ioutil.ReadFile(networkTomlFile)
			if err != nil {
				return fmt.Errorf("Failed to read the toml file at '%s': %v", networkTomlFile, err)
			}

			err = toml.Unmarshal(tomlbytes, &conf)
			if err != nil {
				return nil
			}

			if ips != "" {
				ipmap, err = loadIPS()
				if err != nil {
					return err
				}
			}

			var nodes []node
			var accounts []account
			var debits []*big.Int
			var credits []*big.Int
			var faucetAccount *account
			var trans []transaction
			var deltas []delta
			var nodeTrans []nodeTransactions

			common.DebugMessage("Parsing network.toml for node and accounts")
			for _, n := range conf.Nodes {

				netaddr := n.NetAddr
				moniker := n.Moniker
				balance := n.Tokens

				if netaddr == "" {
					if moniker == faucet {
						faucetAccount = &account{
							Address:   n.Address,
							Moniker:   n.Moniker,
							PubKeyHex: n.PubKeyHex,
							Tokens:    balance,
						}

						common.DebugMessage("faucet ", faucetAccount.Moniker, balance)
					} else {
						accounts = append(accounts, account{
							Address:   n.Address,
							Moniker:   n.Moniker,
							PubKeyHex: n.PubKeyHex,
							Tokens:    balance,
						})

						credits = append(credits, new(big.Int))
						debits = append(debits, new(big.Int))

						common.DebugMessage("account ", moniker, balance)
					}
				} else {
					if ipmap != nil {
						netaddr = ipmap[netaddr]
					}
					if !strings.Contains(netaddr, ":") && len(netaddr) > 0 {
						netaddr += monetconfig.DefaultAPIAddr
					}
					nodes = append(nodes, node{
						NetAddr: netaddr,
						Moniker: moniker})
					common.DebugMessage("node ", moniker)
				}

			}

			nodecnt := len(nodes)
			accountcnt := len(accounts)

			common.InfoMessage("Read " + strconv.Itoa(nodecnt) + " nodes.")
			common.InfoMessage("Read " + strconv.Itoa(accountcnt) + " accounts.")

			if faucetAccount == nil {
				return errors.New("faucet account not found: " + faucet)
			}

			if accountcnt < 2 {
				return errors.New("you must have at least 2 accounts")
			}

			common.InfoMessage("Faucet account, " + faucet + ", found.")

			for i := 0; i < accountcnt; i++ {
				nodeTrans = append(nodeTrans, nodeTransactions{
					Address: accounts[i].Address,
					Moniker: accounts[i].Moniker,
				})
			}

			fulltrans := nodeTransactions{Address: "", Moniker: ""}
			faucettrans := nodeTransactions{Address: faucetAccount.Address, Moniker: faucetAccount.Moniker}

			for i := 0; i < totalTransactions; i++ {
				var fromacct, toacct int

				for {
					fromacct = rand.Intn(accountcnt)
					toacct = rand.Intn(accountcnt)
					if fromacct != toacct {
						break
					}
				}
				//	amt := int64((rand.Intn(1000) * 1000000) + (rand.Intn(1000) * 1000) + rand.Intn(999) + 1)

				amt := new(big.Int).SetInt64(int64(rand.Intn(99990) + 9))
				amt.Mul(amt, new(big.Int).SetInt64(100000000))
				amt.Mul(amt, new(big.Int).SetInt64(100000000))

				// TODO add some lower order bits noise

				credits[toacct].Add(credits[toacct], amt)
				debits[fromacct].Add(debits[fromacct], amt)

				trans = append(trans, transaction{
					From:   fromacct,
					To:     toacct,
					Amount: amt,
				})

				nodeno := rand.Intn(nodecnt)

				newtrans := fulltransaction{
					Node:     nodes[nodeno].NetAddr,
					NodeName: nodes[nodeno].Moniker,
					From:     accounts[fromacct].Address,
					FromName: accounts[fromacct].Moniker,
					To:       accounts[toacct].Address,
					Amount:   amt,
				}

				fulltrans.Transactions = append(fulltrans.Transactions, newtrans)
				nodeTrans[fromacct].Transactions = append(nodeTrans[fromacct].Transactions, newtrans)

				//		common.DebugMessage(strconv.Itoa(i) + ": " + strconv.Itoa(fromacct) + ", " + strconv.Itoa(toacct) + ", " + strconv.FormatInt(amt, 10))
			}

			for i := 0; i < accountcnt; i++ {
				nodeno := rand.Intn(nodecnt)

				faucettrans.Transactions = append(faucettrans.Transactions, fulltransaction{
					From:     faucetAccount.Address,
					FromName: faucetAccount.Moniker,
					To:       accounts[i].Address,
					Amount:   new(big.Int).Add(debits[i], surplusCreditBig),
					Node:     nodes[nodeno].NetAddr,
					NodeName: nodes[nodeno].Moniker,
				})

				deltas = append(deltas, delta{
					Moniker:      accounts[i].Moniker,
					Address:      accounts[i].Address,
					TransCredit:  credits[i],
					TransDebit:   debits[i],
					TransNet:     new(big.Int).Sub(credits[i], debits[i]),
					FaucetCredit: new(big.Int).Add(debits[i], surplusCreditBig),
					TotalNet:     new(big.Int).Add(credits[i], surplusCreditBig),
				})

				common.DebugMessage("Account " + accounts[i].Moniker + " +" +
					credits[i].String() + " -" + debits[i].String())

				var jsonData []byte

				nodefile := filepath.Join(transDir, accounts[i].Moniker+".json")
				jsonData, err = json.Marshal(nodeTrans[i])
				if err != nil {
					return err
				}
				files.WriteToFile(nodefile, string(jsonData))
				common.DebugMessage("Node File written: ", nodefile)

			}

			var jsonData []byte

			jsonData, err = json.Marshal(faucettrans)
			if err != nil {
				return err
			}
			files.WriteToFile(faucetTransFile, string(jsonData))
			common.DebugMessage("Faucet File written: ", faucetTransFile)

			jsonData, err = json.Marshal(fulltrans)
			if err != nil {
				return err
			}
			files.WriteToFile(transFile, string(jsonData))
			common.DebugMessage("Trans File written: ", transFile)

			jsonData, err = json.Marshal(deltas)
			if err != nil {
				return err
			}
			files.WriteToFile(deltaFile, string(jsonData))
			common.DebugMessage("Delta File written: ", deltaFile)

			return nil
		}

		func loadIPS() (map[string]string, error) {
			rtn := make(ipmapping)

			if !files.CheckIfExists(ips) {
				return nil, errors.New("ip mapping does not exist: " + ips)
			}

			file, err := os.Open(ips)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {

				ippair := strings.Fields(scanner.Text())
				if len(ippair) != 2 {
					return nil, errors.New("malformed ip mapping " + strconv.Itoa(len(ippair)))
				}
				rtn[ippair[0]] = ippair[1]
			}

			if err := scanner.Err(); err != nil {
				return nil, err
			}

			return rtn, nil
	*/
}
