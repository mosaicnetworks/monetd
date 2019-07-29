package config

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"

	"github.com/mosaicnetworks/monetd/src/configuration"
)

//ImportZip imports a monetd config zip file (src) and writes to the appropriate
//subfolder in dest
func ImportZip(src string, dest string) error {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		var fpath string

		switch f.Name {
		case configuration.GenesisJSON:
			fpath = filepath.Join(dest, configuration.EthDir, configuration.GenesisJSON)
		case configuration.MonetTomlFile:
			fpath = filepath.Join(dest, configuration.MonetTomlFile)
		case configuration.PeersJSON:
			fpath = filepath.Join(dest, configuration.BabbleDir, configuration.PeersJSON)
		default:
			fpath = filepath.Join(dest, configuration.KeyStoreDir, f.Name)
		}

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		common.DebugMessage("Writing file " + fpath)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
