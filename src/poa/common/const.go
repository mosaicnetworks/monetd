package common

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Directory Constants
const (
	BabbleDir = "babble"
	EthDir    = "eth"
)

// Monetd Configuration Directory
const (
	monetdTomlDirDot  = ".monet"
	monetdTomlDirCaps = "MONET"
)

//DefaultMonetConfigDir is a wrapper for DefaultConfigDir, but returns the
//location of the monetd configuration file.
func DefaultMonetConfigDir() (string, error) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return DefaultConfigDir(monetdTomlDirCaps)
	}
	return DefaultConfigDir(monetdTomlDirDot)
}

//DefaultConfigDir returns a the full path for the default location for a configuration file.
func DefaultConfigDir(configDir string) (string, error) {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", configDir), nil
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", configDir), nil
		} else {
			return filepath.Join(home, configDir), nil
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return "", errors.New("network: cannot determine a sensible default")
}

/* Helper Functions */
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
