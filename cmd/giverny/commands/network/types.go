package network

// node is a struct for defining individual nodes in a network.
type node struct {
	Moniker   string
	NetAddr   string
	Validator bool
	Tokens    string
	PubKeyHex string
	Address   string
}
