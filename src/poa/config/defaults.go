package config

import (
	"path/filepath"
	"time"

	"github.com/mosaicnetworks/monetd/src/poa/network"

	"github.com/mosaicnetworks/monetd/src/poa/common"
)

//KeyValue is a struct for holding toml keys, default values and whether the default overrides the interim supplied value
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
		{Key: "eth.datadir", Value: filepath.Join(monetConfigDir, common.EthDir), Override: true},
		{Key: "eth.genesis", Value: filepath.Join(monetConfigDir, common.EthDir, common.GenesisJSON), Override: true},
		{Key: "eth.keystore", Value: filepath.Join(monetConfigDir, common.EthDir, "keystore"), Override: true},
		//		{Key: "eth.pwd", Value: filepath.Join(monetConfigDir, common.EthDir, PwdFile), Override: true},
		{Key: "eth.db", Value: filepath.Join(monetConfigDir, common.EthDir, "chaindata"), Override: true},
		{Key: "eth.listen", Value: ":" + common.DefaultEVMLitePort, Override: false},
		{Key: "eth.cache", Value: "128", Override: false},

		{Key: "babble.datadir", Value: filepath.Join(monetConfigDir, common.BabbleDir), Override: true},
		{Key: "babble.listen", Value: network.GetMyIP() + ":" + common.DefaultGossipPort, Override: false},
		{Key: "babble.service-listen", Value: ":" + common.DefaultBabblePort, Override: false},
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
