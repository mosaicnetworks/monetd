//Package configuration contains global configuration options for Giverny
package configuration

import (
	"runtime"

	"github.com/mosaicnetworks/monetd/src/configuration"
)

//GivernyConfigDir is the root config directory for Giverny.
var GivernyConfigDir = defaultGivernyConfigDir()

//GivernyExportDir is the exports subfolder of the Giverny config folder
var GivernyExportDir = "exports"

const (
	givernyTomlDirCaps = "Giverny"
	givernyTomlDirDot  = ".giverny"
	// GivernyServerPort is the default port for Giverny Server
	GivernyServerPort = "8088"
)

func defaultGivernyConfigDir() string {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		rtn, _ := configuration.DefaultConfigDir(givernyTomlDirCaps)
		return rtn
	}
	rtn, _ := configuration.DefaultConfigDir(givernyTomlDirDot)
	return rtn
}
