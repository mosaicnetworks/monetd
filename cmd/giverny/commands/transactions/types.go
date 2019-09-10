package transactions

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
	Amount int64
}

type fulltransaction struct {
	Node     string
	NodeName string
	From     string
	FromName string
	To       string
	Amount   int64
}

type delta struct {
	Moniker      string
	Address      string
	TransCredit  int64
	TransDebit   int64
	TransNet     int64
	FaucetCredit int64
	TotalNet     int64
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
