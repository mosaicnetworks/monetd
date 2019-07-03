package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fatih/color"
)

//General file wrapper code

//DefaultHomeDir returns a default location for a configuration file.
func DefaultHomeDir(tomlDir string) (string, error) {
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

//CheckIsDir checks if file is a directory
func CheckIsDir(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

//CheckIfExists checks if a file / directory exists
func CheckIfExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

//WriteToFile writes a string variable to a file.
//It overwrites any pre-existing data silently.
func WriteToFile(filename string, data string) error {
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

//SafeRenameDir renames a folder to folder.~n~ where n is the lowest value
//where the folder does not already exist.
//n is capped at 100 - which would require the user to manually tidy the parent folder.
func SafeRenameDir(origDir string) error {
	const maxloops = 100

	for i := 1; i < 100; i++ {
		newDir := origDir + ".~" + strconv.Itoa(i) + "~"
		if CheckIfExists(newDir) {
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

//CopyFileContents writes the contents to src to a new file dst.
//This operation is silently destructive
func CopyFileContents(src, dst string) (err error) {

	Message("Copying from " + src + " to " + dst)
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

//ShowConfigFile echoes the given file to screen
func ShowConfigFile(filename string) error {

	color.Set(ColourOutput)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Print(string(b))

	color.Unset()
	return nil
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
