package configuration

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// GivernyConfigDir is the absolute path of the giverny configuration directory
var GivernyConfigDir = defaultGivernyDir()

const (
	// GivernyNetworkDir is the networks subfolder of the Giverny config folder
	GivernyNetworkDir = "networks"
)

// defaultGivernyDir returns the full path for Giverny's data directory.
func defaultGivernyDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Giverny")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Giverny")
		} else {
			return filepath.Join(home, ".giverny")
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
