package network

type configurationRecord struct {
	//location of the network.toml file
	config     *configRecord    `mapstructure:"config"`
	validators *validatorRecord `mapstructure:"validators"`
	poa        *poaRecord       `mapstructure:"config"`
}

type configRecord struct {
	dataDir string `mapstructure:"datadir"`
}

type poaRecord struct {
	contractAddress string `mapstructure:"contractaddress"`
	contractName    string `mapstructure:"contractname"`
	compilerVersion string `mapstructure:"compilerversion"`
	byteCode        string `mapstructure:"bytecode"`
	abi             string `mapstructure:"abi"`
}

type validatorRecord struct {
	moniker            string `mapstructure:"monikers"`
	address            string `mapstructure:"addresses"`
	ip                 string `mapstructure:"ips"`
	isInitialValidator string `mapstructure:"isvalidator"`
}

var (
	config       configurationRecord
	configConfig configRecord
	poa          poaRecord
)

const (
	defaultContractAddress = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
	defaultContractName    = "genesis_array.sol"
)

func defaultConfig() {

	home, err := defaultHomeDir()
	if err == nil {
		networkViper.SetDefault("config.datadir", home)
	}
	networkViper.SetDefault("poa.contractaddress", defaultContractAddress)
	networkViper.SetDefault("poa.contractname", defaultContractName)
	networkViper.SetDefault("poa.compilerversion", "")
	networkViper.SetDefault("validators.monikers", "")
	networkViper.SetDefault("validators.addresses", "")
	networkViper.SetDefault("validators.ips", "")
	networkViper.SetDefault("validators.isvalidator", "")
}

func newConfigurationRecord() *configurationRecord {

	configConfig = configRecord{dataDir: ""}
	poa = poaRecord{contractAddress: defaultContractAddress, contractName: defaultContractName}

	config = configurationRecord{config: &configConfig, poa: &poa}

	return &config

}
