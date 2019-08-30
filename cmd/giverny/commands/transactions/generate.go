package transactions

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

	var nodes []string
	var accounts []account
	var debits []int64
	var credits []int64
	var faucetAccount *account
	var trans []transaction
	var fulltrans []fulltransaction
	var faucettrans []fulltransaction
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

				credits = append(credits, 0)
				debits = append(debits, 0)

				common.DebugMessage("account ", moniker, balance)
			}
		} else {
			if ipmap != nil {
				netaddr = ipmap[netaddr]
			}
			if !strings.Contains(netaddr, ":") && len(netaddr) > 0 {
				netaddr += monetconfig.DefaultAPIAddr
			}
			nodes = append(nodes, netaddr)
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

	for i := 0; i < totalTransactions; i++ {
		var fromacct, toacct int

		for {
			fromacct = rand.Intn(accountcnt)
			toacct = rand.Intn(accountcnt)
			if fromacct != toacct {
				break
			}
		}
		amt := int64((rand.Intn(1000) * 1000000) + (rand.Intn(1000) * 1000) + rand.Intn(999) + 1)

		credits[toacct] += amt
		debits[fromacct] += amt

		trans = append(trans, transaction{
			From:   fromacct,
			To:     toacct,
			Amount: amt,
		})

		newtrans := fulltransaction{
			Node:   nodes[rand.Intn(nodecnt)],
			From:   accounts[fromacct].Address,
			To:     accounts[toacct].Address,
			Amount: amt,
		}

		fulltrans = append(fulltrans, newtrans)
		nodeTrans[fromacct].Transactions = append(nodeTrans[fromacct].Transactions, newtrans)

		//		common.DebugMessage(strconv.Itoa(i) + ": " + strconv.Itoa(fromacct) + ", " + strconv.Itoa(toacct) + ", " + strconv.FormatInt(amt, 10))
	}

	for i := 0; i < accountcnt; i++ {

		faucettrans = append(faucettrans, fulltransaction{
			From:   faucetAccount.Address,
			To:     accounts[i].Address,
			Amount: debits[i] + int64(surplusCredit),
			Node:   nodes[rand.Intn(nodecnt)],
		})

		deltas = append(deltas, delta{
			Address:      accounts[i].Address,
			TransCredit:  credits[i],
			TransDebit:   debits[i],
			TransNet:     credits[i] - debits[i],
			FaucetCredit: debits[i] + int64(surplusCredit),
			TotalNet:     credits[i] + int64(surplusCredit),
		})

		common.DebugMessage("Account " + accounts[i].Moniker + " +" + strconv.FormatInt(credits[i], 10) + " -" + strconv.FormatInt(debits[i], 10))

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
