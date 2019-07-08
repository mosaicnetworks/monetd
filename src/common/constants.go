package common

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

//Constants used across the MonetCLI command packages

var (
	//MonetdTomlDir contains the directory name for the monetd config
	//It varies by OS
	MonetdTomlDir string
	//MonetcliTomlDir contains the directory name for the monetd config
	//It varies by OS
	MonetcliTomlDir string
	//EvmlcTomlDir contains the directory name for the monetd config
	//It varies by OS
	EvmlcTomlDir string
)

//Constants for files names and default values used in building configurations
const (
	MonetdTomlDirDot    = ".monet"
	MonetcliTomlDirDot  = ".monetcli"
	MonetdTomlDirCaps   = "MONET"
	MonetcliTomlDirCaps = "MONETCLI"

	BabbleDir = "babble"
	EthDir    = "eth"
	PwdFile   = "pwd.txt"

	EvmlcTomlDirCaps = "EVMLC"
	EvmlcTomlDirDot  = ".evmlc"

	TomlSuffix = ".toml"

	MonetdTomlName   = "monetd"
	MonetcliTomlName = "network"
	EvmlcTomlName    = "config"

	PeersJSON        = "peers.json"
	PeersGenesisJSON = "peers.genesis.json"
	GenesisJSON      = "genesis.json"

	PeersJSONTarget        = "babble/peers.json"
	PeersGenesisJSONTarget = "babble/peers.genesis.json"
	GenesisJSONTarget      = "eth/genesis.json"

	//TODO Change this value when monetd repo is public
	DefaultSolidityContract = "https://raw.githubusercontent.com/mosaicnetworks/evm-lite/poa/e2e/smart-contracts/monet.sol"
	//	DefaultSolidityContract = "https://raw.githubusercontent.com/mosaicnetworks/monetd/master/smart-contract/genesis.sol"
	TemplateContract       = "template.sol"
	GenesisContract        = "contract0.sol"
	GenesisABI             = "contract0.abi"
	DefaultAccountBalance  = "1234000000000000000000"
	DefaultContractAddress = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
)

func init() {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		MonetdTomlDir = MonetdTomlDirCaps
		MonetcliTomlDir = MonetcliTomlDirCaps
		EvmlcTomlDir = EvmlcTomlDirCaps
		return
	}
	MonetdTomlDir = MonetdTomlDirDot
	MonetcliTomlDir = MonetcliTomlDirDot
	EvmlcTomlDir = EvmlcTomlDirDot

}

//KeyValue is a struct for holding toml keys, default values and whether the default overrides the monetcli supplied value
type KeyValue struct {
	Key      string
	Value    interface{}
	Override bool
	Prompt   string
	Answers  []string
	DataType string
}

//GetMonetDefaultConfigKeys defines the default config values
func GetMonetDefaultConfigKeys(monetConfigDir string) []KeyValue {

	return []KeyValue{
		{Key: "datadir", Value: monetConfigDir, Override: true},
		{Key: "log", Value: "debug", Override: false, Prompt: "Logging Level", Answers: []string{"debug", "info", "warn", "error", "fatal", "panic"}},
		{Key: "eth.datadir", Value: filepath.Join(monetConfigDir, EthDir), Override: true},
		{Key: "eth.genesis", Value: filepath.Join(monetConfigDir, EthDir, GenesisJSON), Override: true},
		{Key: "eth.keystore", Value: filepath.Join(monetConfigDir, EthDir, "keystore"), Override: true},
		{Key: "eth.pwd", Value: filepath.Join(monetConfigDir, EthDir, PwdFile), Override: true},
		{Key: "eth.db", Value: filepath.Join(monetConfigDir, EthDir, "chaindata"), Override: true},
		{Key: "eth.listen", Value: ":8080", Override: false},
		{Key: "eth.cache", Value: "128", Override: false},

		{Key: "babble.datadir", Value: filepath.Join(monetConfigDir, BabbleDir), Override: true},
		{Key: "babble.listen", Value: ":1337", Override: false},
		{Key: "babble.service-listen", Value: ":8000", Override: false},
		{Key: "babble.heartbeat", Value: time.Duration(500 * time.Millisecond).String(), Override: false},
		{Key: "babble.timeout", Value: time.Duration(1000 * time.Millisecond).String(), Override: false},
		{Key: "babble.cache-size", Value: "50000", Override: false},
		{Key: "babble.sync-limit", Value: "1000", Override: false},

		{Key: "babble.fast-sync", Value: false, Override: false, Answers: []string{"false", "true"}, DataType: "bool"},
		{Key: "babble.max-pool", Value: "2", Override: false},
		{Key: "babble.store", Value: true, Override: true},
		{Key: "babble.bootstrap", Value: true, Override: false},
	}
}

//Colour constants used in the functions in src/common/logs.go
const (
	ColourInfo    = color.FgGreen
	ColourWarning = color.FgHiMagenta
	ColourError   = color.FgHiRed
	ColourPrompt  = color.FgHiYellow
	ColourOther   = color.FgYellow
	ColourOutput  = color.FgHiCyan
	ColourDebug   = color.FgCyan
)
