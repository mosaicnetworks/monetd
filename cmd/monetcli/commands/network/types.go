package network

import (
	"github.com/mosaicnetworks/monetd/src/common"
)

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
	compilerVersion string `mapstructure:"compilerversion"`
	byteCode        string `mapstructure:"bytecode"`
	abi             string `mapstructure:"abi"`
}

type validatorRecord struct {
	moniker     string `mapstructure:"moniker"`
	label       string `mapstructure:"label"`
	address     string `mapstructure:"address"`
	pubkey      string `mapstructure:"pubkey"`
	ip          string `mapstructure:"ip"`
	isValidator bool   `mapstructure:"validator"`
}

type validatorRecordList map[string]*validatorRecord

var (
	config       configurationRecord
	configConfig configRecord
	poa          poaRecord
)

/* func defaultConfig() {

	home, err := common.DefaultHomeDir(common.MonetcliTomlDir)
	if err == nil {
		networkViper.SetDefault("config.datadir", home)
	}
	networkViper.SetDefault("poa.contractaddress", common.DefaultContractAddress)
	networkViper.SetDefault("poa.compilerversion", "")
	networkViper.SetDefault("validators.monikers", "")
	networkViper.SetDefault("validators.addresses", "")
	networkViper.SetDefault("validators.pubkeys", "")
	networkViper.SetDefault("validators.ips", "")
	networkViper.SetDefault("validators.isvalidator", "")
} */

func newConfigurationRecord() *configurationRecord {

	configConfig = configRecord{dataDir: ""}
	poa = poaRecord{contractAddress: common.DefaultContractAddress}

	config = configurationRecord{config: &configConfig, poa: &poa}

	return &config

}
