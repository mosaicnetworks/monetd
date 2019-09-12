package transactions

import "math/big"

// CLI Parameters
var networkName = ""
var ips = ""
var faucet = "Faucet"
var totalTransactions = 20
var surplusCredit = 1000000

type ipmapping map[string]string

type account struct {
	Moniker   string
	Tokens    string
	PubKeyHex string
	Address   string
}

type transaction struct {
	From   int
	To     int
	Amount *big.Int
}

type fulltransaction struct {
	Node     string
	NodeName string
	From     string
	FromName string
	To       string
	Amount   *big.Int
}

type delta struct {
	Moniker      string
	Address      string
	TransCredit  *big.Int
	TransDebit   *big.Int
	TransNet     *big.Int
	FaucetCredit *big.Int
	TotalNet     *big.Int
}

type node struct {
	Moniker string
	NetAddr string
}

type nodeTransactions struct {
	Moniker      string
	Address      string
	Transactions []fulltransaction
}
