//Package files provides standard file functions
package files

import (
	"io"
	"os"
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
