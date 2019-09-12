package transactions

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/cmd/giverny/commands/network"

	"github.com/BurntSushi/toml"
	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate transactions",
		Long: `
Generate transactions.
`,
		Args: cobra.ArbitraryArgs,
		RunE: generateTransactions,
	}

	addGenerateFlags(cmd)
	return cmd
}

func addGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&networkName, "network", "n", "", "network name")
	cmd.Flags().StringVar(&ips, "ips", "", "ips.dat file path")
	cmd.Flags().StringVar(&faucet, "faucet", faucet, "faucet account moniker")
	cmd.Flags().IntVar(&totalTransactions, "count", totalTransactions, "number of tranactions to generate")
	cmd.Flags().IntVar(&surplusCredit, "surplus", surplusCredit, "additional credit to allocate each account from the faucet above the bare minimum")
	viper.BindPFlags(cmd.Flags())
}

func generateTransactions(cmd *cobra.Command, args []string) error {

	// N.B. this function used int64 to generate transaction amounts.
	// It is good enough - there is sufficent range to be a realistic test
	// We could use bit.Int - but it appears unnecessary.

	var ipmap ipmapping

	var surplusCreditBig = new(big.Int).SetInt64(int64(surplusCredit))

	if !common.CheckMoniker(networkName) {
		return errors.New("network name must only contains characters in the range 0-9 or A-Z or a-z")
	}

	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	if !files.CheckIfExists(networkDir) {
		return errors.New("network does not exist: " + networkDir)
	}

	keystore := filepath.Join(networkDir, givernyKeystoreDir)
	if !files.CheckIfExists(networkDir) {
		return errors.New("keystore does not exist: " + keystore)
	}

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
}
