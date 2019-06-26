package network

import (
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

const tomlDir = ".monetcli"
const (
	tomlName = "network"
)

// We declare our own instance of viper to avoid any possibility of a clash
var (
	networkViper *viper.Viper
)

func setUpConfigFile() {
	networkViper = viper.New()
	networkViper.SetConfigName(tomlName) // name of config file (without extension)
	networkViper.SetConfigType("toml")
	defaultConfig()
}

// Write configure to file
func writeConfig() {

	message("Writing toml file")
	err := networkViper.WriteConfig()
	if err != nil {
		message("writeConfig error: ", err)
	}
}

func safeWriteConfig() {

	message("Writing toml file")
	err := networkViper.SafeWriteConfig()
	if err != nil {
		message("safeWriteConfig error: ", err)
		message(networkViper.AllSettings())
		networkViper.Debug()
	}
}

func createEmptyFile(f string) {
	emptyFile, err := os.Create(f)

	if err != nil {
		message("Create empty file: ", f, err)
		return
	}
	emptyFile.Close()
}

func loadConfig() error {

	message("Starting to load configuration")
	setUpConfigFile()
	networkViper.AddConfigPath(configDir)
	message("Added viper config path: ", configDir)
	err := networkViper.ReadInConfig() // Find and read the config file

	if err != nil {
		message("loadConfig: ", err)
		return err
	}
	message("Loaded Config")
	return nil
}

func isEmptyDir(dir string) (bool, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
