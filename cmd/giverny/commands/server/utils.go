package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/src/configuration"
)

var pidFile string
var logOut string

// "/tmp/daemonize.pid"

func init() {

	configDir := filepath.Join(configuration.Global.DataDir, ServerDir)
	pidFile = filepath.Join(configDir, ServerPIDFile)
	logOut = filepath.Join(configDir, ServerLogFile)

	files.CreateDirsIfNotExists([]string{configDir})
}

func savePID(pid string) error {

	file, err := os.Create(pidFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		return err
	}

	defer file.Close()

	_, err = file.WriteString(pid)

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		return err
	}

	file.Sync() // flush to disk
	return nil
}
