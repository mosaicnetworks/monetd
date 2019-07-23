//Package files provides standard file functions
package files

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/poa/common"
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

//CreateMonetConfigFolders creates the standard directory layout for a configuration folder
func CreateMonetConfigFolders(configDir string) error {

	return CreateDirsIfNotExists([]string{
		configDir,
		filepath.Join(configDir, common.BabbleDir),
		filepath.Join(configDir, common.EthDir),
		filepath.Join(configDir, common.KeyStoreDir),
		filepath.Join(configDir, common.EthDir, common.POADir),
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
