package configuration

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Directory Constants
const (
	// config
	ConfigDir = "monetd-config"
	BabbleDir = "babble"
	EthDir    = "eth"
	POADir    = "poa"

	// data
	DatabaseDir = "monetd-data"

	// keystore
	KeyStoreDir = "keystore"
)

// Monetd Configuration Directory
const (
	MonetdTomlDirDot  = ".monet"
	MonetdTomlDirCaps = "Monet"
)

// Filename constants
const (
	PeersJSON        = "peers.json"
	PeersGenesisJSON = "peers.genesis.json"
	GenesisJSON      = "genesis.json"
	MonetTomlFile    = "monetd.toml"
	EthDB            = "eth-db"
	BabbleDB         = "babble-db"
	WalletTomlFile   = "wallet.toml"
	ServerPIDFile    = "server.pid"
	BabblePrivKey    = "priv_key"
)

// Network Constants
const (
	DefaultGossipPort = "1337"
	DefaultAPIAddr    = ":8080"
)

//Keys constants
const (
	DefaultKeyfile        = "keyfile.json"
	DefaultPrivateKeyFile = "priv_key"
)

// Genesis Constants
const (
	DefaultAccountBalance  = "1234567890000000000000"
	DefaultContractAddress = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
	GenesisContract        = "contract0.sol"
	GenesisABI             = "contract0.abi"
	CompileResultFile      = "compile.toml"
)

// DefaultConfigDir returns the full path of the config directory where static
// configuration files are stored.
func DefaultConfigDir() string {
	return filepath.Join(DefaultMonetDir(), ConfigDir)
}

// DefaultDataDir returns the full path of the data directory where databases
// are stored.
func DefaultDataDir() string {
	return filepath.Join(DefaultMonetDir(), DatabaseDir)
}

// DefaultKeystoreDir returns the full path of the keystore where encrypted
// keyfiles are stored.
func DefaultKeystoreDir() string {
	return filepath.Join(DefaultMonetDir(), "keystore")
}

// DefaultMonetDir returns a the full path for the default location Monet
// configuration files based on the underlying OS.
func DefaultMonetDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Monet")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Monet")
		} else {
			return filepath.Join(home, ".monet")
		}
	}
	return ""
}

// Guess a sensible default location from OS and environment variables.
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
