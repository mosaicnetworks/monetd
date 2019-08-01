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
	Docker  dockerConfig  `toml:"docker"`
	Nodes   []node        `toml:"nodes,squash"`
}

type networkConfig struct {
	Name string `toml:"name"`
}

type dockerConfig struct {
	Name    string `toml:"name"`
	BaseIP  string `toml:"baseip"`
	Subnet  string `toml:"subnet"`
	IPRange string `toml:"iprange"`
	Gateway string `toml:"gateway"`
}

type dockerNodeConfig struct {
	Moniker string `toml:"moniker"`
	NetAddr string `toml:"netaddr"`
}
