package common

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

//Constants used across the MonetCLI command packages

var (
	MonetdTomlDir   string
	MonetcliTomlDir string
	EvmlcTomlDir    string
)

const (
	MonetdTomlDirDot    = ".monet"
	MonetcliTomlDirDot  = ".monetcli"
	MonetdTomlDirCaps   = "MONET"
	MonetcliTomlDirCaps = "MONETCLI"

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

	DefaultSolidityContract = "https://raw.githubusercontent.com/mosaicnetworks/evm-lite/poa/e2e/smart-contracts/genesis_array.sol"
	TemplateContract        = "template.sol"
	GenesisContract         = "contract0.sol"
	GenesisABI              = "contract0.abi"
	DefaultAccountBalance   = "1234000000000000000000"
	DefaultContractAddress  = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
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
		KeyValue{Key: "datadir", Value: monetConfigDir, Override: true},
		KeyValue{Key: "log", Value: "debug", Override: false, Prompt: "Logging Level", Answers: []string{"debug", "info", "warn", "error", "fatal", "panic"}},
		KeyValue{Key: "eth.datadir", Value: filepath.Join(monetConfigDir, "eth"), Override: true},
		KeyValue{Key: "eth.genesis", Value: filepath.Join(monetConfigDir, "eth", "genesis.json"), Override: true},
		KeyValue{Key: "eth.keystore", Value: filepath.Join(monetConfigDir, "eth", "keystore"), Override: true},
		KeyValue{Key: "eth.pwd", Value: filepath.Join(monetConfigDir, "eth", "pwd.txt"), Override: true},
		KeyValue{Key: "eth.db", Value: filepath.Join(monetConfigDir, "eth", "chaindata"), Override: true},
		KeyValue{Key: "eth.listen", Value: ":8080", Override: false},
		KeyValue{Key: "eth.cache", Value: "128", Override: false},

		KeyValue{Key: "babble.datadir", Value: filepath.Join(monetConfigDir, "babble"), Override: true},
		KeyValue{Key: "babble.listen", Value: ":1337", Override: false},
		KeyValue{Key: "babble.service-listen", Value: ":8000", Override: false},
		KeyValue{Key: "babble.heartbeat", Value: time.Duration(500 * time.Millisecond).String(), Override: false},
		KeyValue{Key: "babble.timeout", Value: time.Duration(1000 * time.Millisecond).String(), Override: false},
		KeyValue{Key: "babble.cache-size", Value: "50000", Override: false},
		KeyValue{Key: "babble.sync-limit", Value: "1000", Override: false},

		KeyValue{Key: "babble.fast-sync", Value: false, Override: false, Answers: []string{"false", "true"}, DataType: "bool"},
		KeyValue{Key: "babble.max-pool", Value: "2", Override: false},
		KeyValue{Key: "babble.store", Value: true, Override: true},
	}
}

const (
	ColourInfo    = color.FgGreen
	ColourWarning = color.FgHiMagenta
	ColourError   = color.FgHiRed
	ColourPrompt  = color.FgHiYellow
	ColourOther   = color.FgYellow
	ColourOutput  = color.FgHiCyan
)
