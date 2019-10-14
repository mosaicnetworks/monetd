package config

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
)

const configTOML = `
# Set to true for extended logging
verbose = "{{.Verbose}}"
# The IP:PORT of the HTTP API service (defaults to :8080)
api-listen = "{{.APIAddr}}"

[babble]
# Advertise IP:PORT of Babble node. This is shared in the peers.json
# file and should not be an internal IP
  listen = "{{.Babble.BindAddr}}"
# The heartbeat timer, the time in milliseconds between gossips  
  heartbeat = "{{.Babble.Heartbeat}}"
# TCP timeout  
  timeout = "{{.Babble.TCPTimeout}}"
# Number of items in the LRU cache  
  cache-size = {{.Babble.CacheSize}}
# Max number of events per sync  
  sync-limit = {{.Babble.SyncLimit}}
# Max number of pool connections  
  max-pool = {{.Babble.MaxPool}}
# Bootstrap Babble from database  
  bootstrap = {{.Babble.Bootstrap}}
# Moniker for this node  
  moniker = "{{.Babble.Moniker}}"

[eth]
# megabytes of memory allocated to internal caching 
#    (min 16MB / database forced) (default 128)
  cache = {{.Eth.Cache}}
# minimum gasprice of transactions submitted through this node (eg 1T) (default "0")  
  min-gas-price = {{.Eth.MinGasPrice}}
`

// GlobalTOML converts the global Config object into a TOML string
func GlobalTOML() (string, error) {
	configTmpl, err := template.New("monetd.toml").Parse(configTOML)
	if err != nil {
		return "", fmt.Errorf("Error parsing monetd.toml template: %v", err)
	}

	var buf bytes.Buffer
	err = configTmpl.Execute(&buf, configuration.Global)
	if err != nil {
		return "", fmt.Errorf("Error executing monetd.toml template: %v", err)
	}

	return buf.String(), nil
}

// DumpGlobalTOML takes the global Config object, encodes it into a TOML string,
// and writes it to a file.
func DumpGlobalTOML(configDir, fileName string) error {
	tomlString, err := GlobalTOML()
	if err != nil {
		return err
	}

	tomlPath := filepath.Join(configDir, fileName)

	if err := files.ProcessFileOptions(tomlPath, files.BackupExisting); err != nil {
		return err
	}

	if err := ioutil.WriteFile(tomlPath, []byte(tomlString), 0644); err != nil {
		return fmt.Errorf("Failed to write %s: %v", tomlPath, err)
	}

	configuration.ShowIPWarnings()
	return nil
}
