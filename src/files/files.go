//Package files provides standard file functions
package files

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
)

// Bits is used to hold bitwise options
type Bits uint8

// File options
const (
	BackupExisting Bits = 1 << iota
	PromptIfExisting
	OverwriteSilently
)

// ProcessFileOptions decides what to do with an existing file before any
// modifications (create, copy, move, delete, update), based on the options
// provided.
// If the OverwriteSilently bit is set, this function does nothing, allowing the
// file to be silently modified later. Otherwise, if the PromptIfExisting bit is
// set, the user is interactively prompted for confirmation (yes/no). If the
// BackupExisting flag is set, the file is also backed-up.
func ProcessFileOptions(filename string, options Bits) error {
	if CheckIfExists(filename) {
		if options&OverwriteSilently == 0 {
			if options&PromptIfExisting != 0 {
				for {
					reader := bufio.NewReader(os.Stdin)
					common.PromptMessage(fmt.Sprintf("File %s already exists. Overwrite (yes/no)?: ", filename))
					name, _ := reader.ReadString('\n')
					fmt.Println("")
					tidiedName := strings.ToLower(strings.TrimSpace(name))
					if tidiedName == "no" || tidiedName == "n" {
						return errors.New("you declined to overwrite " + filename)
					}
					if tidiedName == "yes" {
						break
					}
					common.PromptMessage("You must type 'yes' or 'no'")
				}
			}
			if options&BackupExisting != 0 {
				SafeRename(filename)
			}
		}
	}
	return nil
}

// WriteToFile calls ProcessFileOptions before writing a string variable to a
// file.
func WriteToFile(filename string, data string, options Bits) error {
	if err := ProcessFileOptions(filename, options); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, []byte(data), 0644); err != nil {
		return fmt.Errorf("Failed to write %s: %v", filename, err)
	}

	return nil
}

// WriteToFilePrivate writes a string variable to a file with 0600 permissions.
// It creates all directories along the path if they don't exist.
func WriteToFilePrivate(filename string, data string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0744); err != nil {
		return err
	}
	return ioutil.WriteFile(
		filename,
		[]byte(data),
		0600)
}

// CreateDirsIfNotExists takes an array of strings containing filepaths and for
// any path that contains directories which do not exist, it creates them.
func CreateDirsIfNotExists(d []string) error {
	for _, dir := range d {
		if err := os.MkdirAll(dir, 0744); err != nil {
			return err
		}
	}
	return nil
}

//CheckIfExists checks if a file / directory exists
func CheckIfExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// CopyFileContents writes the contents from src to a new file dst. This
// operation is silently destructive
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

// SafeRename renames a folder or file to <name>.~n~ where n is the lowest value
// where the folder does not already exist. The value of n is capped at 100,
// which would require the user to manually tidy the parent folder.
func SafeRename(origDir string) error {

	// no renaming is necessary if the original file/folder doesnt exist
	if !CheckIfExists(origDir) {
		return nil
	}

	const maxloops = 100

	for i := 1; i < 100; i++ {
		newDir := origDir + ".~" + strconv.Itoa(i) + "~"
		if CheckIfExists(newDir) {
			continue
		}
		common.DebugMessage("Renaming " + origDir + " to " + newDir)
		err := os.Rename(origDir, newDir)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("you have reached the maximum number of automatic backups. Try removing the .monet.~n~ files")
}

//DownloadFile downs a file from a URL and writes it to disk
func DownloadFile(url string, writefile string, interactive bool) error {
	b, err := getRequest(url)
	if err != nil {
		common.ErrorMessage("Error getting "+url, err)
		return err
	}

	// Set options for WriteToFile base on 'interactive' flag
	var options Bits
	if interactive {
		options = PromptIfExisting
	}

	err = WriteToFile(writefile, string(b), options)
	if err != nil {
		common.ErrorMessage("Error writing "+writefile, err)
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
