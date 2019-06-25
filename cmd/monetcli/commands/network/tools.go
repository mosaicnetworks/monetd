package network

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

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

func defaultHomeDir() (string, error) {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", tomlDir), nil
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", tomlDir), nil
		} else {
			return filepath.Join(home, tomlDir), nil
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return "", errors.New("network: cannot determine a sensible default")
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func checkIsDir(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

func checkIfExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func isEmptyDir(dir string) (bool, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func safeRenameDir(origDir string) error {
	const maxloops = 100

	for i := 1; i < 100; i++ {
		newDir := origDir + ".~" + strconv.Itoa(i) + "~"
		if checkIfExists(newDir) {
			continue
		}
		fmt.Println("Renaming " + origDir + " to " + newDir)
		err := os.Rename(origDir, newDir)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("you have reached the maximum number of automatic backups. Try removing the .monet.~n~ files")
}

//Check if is a valid ethereum style address
func isValidAddress(v string) bool {

	re := regexp.MustCompile("^(0[xX]){0,1}[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
