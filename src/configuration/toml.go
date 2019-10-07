package configuration

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
)

const configTOML = `
verbose = "{{.Verbose}}"
api-listen = "{{.APIAddr}}"

[babble]
  listen = "{{.Babble.BindAddr}}"
  heartbeat = "{{.Babble.Heartbeat}}"
  timeout = "{{.Babble.TCPTimeout}}"
  cache-size = {{.Babble.CacheSize}}
  sync-limit = {{.Babble.SyncLimit}}
  max-pool = {{.Babble.MaxPool}}
  bootstrap = {{.Babble.Bootstrap}}
  moniker = "{{.Babble.Moniker}}"

[eth]
  cache = {{.Eth.Cache}}
  min-gas-price = {{.Eth.MinGasPrice}}
`

// GlobalTOML converts the global Config object into a TOML string
func GlobalTOML() (string, error) {
	configTmpl, err := template.New("monetd.toml").Parse(configTOML)
	if err != nil {
		return "", fmt.Errorf("Error parsing monetd.toml template: %v", err)
	}

	var buf bytes.Buffer
	err = configTmpl.Execute(&buf, Global)
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

	if err := ioutil.WriteFile(tomlPath, []byte(tomlString), 0644); err != nil {
		return fmt.Errorf("Failed to write %s: %v", tomlPath, err)
	}

	return nil
}
