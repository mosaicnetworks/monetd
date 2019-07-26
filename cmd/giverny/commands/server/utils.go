package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/poa/files"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/poa/common"
)

var pidFile string
var logOut string
var logErr string

// "/tmp/daemonize.pid"

func init() {

	configDir := filepath.Join(configuration.Global.DataDir, common.ServerDir)
	pidFile = filepath.Join(configDir, common.ServerPIDFile)
	logOut = filepath.Join(configDir, "log.out")
	logErr = filepath.Join(configDir, "error.out")

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
