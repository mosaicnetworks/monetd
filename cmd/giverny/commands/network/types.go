package network

// node is a struct for defining individual nodes in a network.
type node struct {
	Moniker   string `toml:"moniker"`
	NetAddr   string `toml:"netaddr"`
	Validator bool   `toml:"validator"`
	Tokens    string `toml:"tokens"`
	PubKeyHex string `toml:"pubkey"`
	Address   string `toml:"address"`
}

//Config defines the nodes.toml file structure
type Config struct {
	Network networkConfig `toml:"network"`
	Nodes   []node        `toml:"nodes,squash"`
}

type networkConfig struct {
	Name string `toml:"name"`
}
