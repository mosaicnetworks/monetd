//Package files provides standard file functions
package files

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
)

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

//WriteToFilePrivate writes a string variable to a file with 0600 permissions
func WriteToFilePrivate(filename string, data string) error {

	return ioutil.WriteFile(
		filename,
		[]byte(data), 0600)
}

//CreateDirsIfNotExists takes an array of strings contain filepaths and
//any that for not exist are created.
func CreateDirsIfNotExists(d []string) error {

	for _, dir := range d {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				common.ErrorMessage("Error creating directory: ", dir)
				return err
			}
			common.DebugMessage("Created Directory: ", dir)
		} else {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateMonetConfigFolders creates the standard directory layout for a
// monet configuration folder
func CreateMonetConfigFolders(configDir string) error {
	return CreateDirsIfNotExists([]string{
		configDir,
		filepath.Join(configDir, configuration.BabbleDir),
		filepath.Join(configDir, configuration.EthDir),
		filepath.Join(configDir, configuration.KeyStoreDir),
		filepath.Join(configDir, configuration.EthDir, configuration.POADir),
	})
}

//CheckIfExists checks if a file / directory exists
func CheckIfExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

//CopyFileContents writes the contents to src to a new file dst.
//This operation is silently destructive
func CopyFileContents(src, dst string) (err error) {

	common.DebugMessage("Copying from " + src + " to " + dst)
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

//SafeRenameDir renames a folder to folder.~n~ where n is the lowest value
//where the folder does not already exist.
//n is capped at 100 - which would require the user to manually tidy the parent folder.
func SafeRenameDir(origDir string) error {

	// XXX no renaming to do if the original file/folder doesnt exist
	if !CheckIfExists(origDir) {
		return nil
	}

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

//DownloadFile downs a file from a URL and writes it to disk
func DownloadFile(url string, writefile string) error {
	b, err := getRequest(url)
	if err != nil {
		fmt.Println("Error getting "+url, err)
		return err
	}

	err = WriteToFile(writefile, string(b))
	if err != nil {
		fmt.Println("Error writing "+writefile, err)
		return err
	}
	return nil
}

//GetRequest gets a request from a URL
func getRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytes, nil
}
