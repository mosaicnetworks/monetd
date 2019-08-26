package configuration

import (
	"bytes"
	"fmt"
	"html/template"
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

[eth]
  cache = {{.Eth.Cache}}
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
