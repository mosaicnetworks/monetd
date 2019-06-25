package version

import (
	"fmt"

	geth "github.com/ethereum/go-ethereum/params"
	_babble "github.com/mosaicnetworks/babble/src/version"
	evm "github.com/mosaicnetworks/evm-lite/src/version"
)

const Maj = "0"
const Min = "0"
const Fix = "1"

var (
	// The full version string
	Version = "0.0.1"

	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}
}

func FullVersion() string {
	return fmt.Sprintln("Monet Version: "+Version) +
		fmt.Sprintln("     EVM-Lite Version: "+evm.Version) +
		fmt.Sprintln("     Babble Version: "+_babble.Version) +
		fmt.Sprintln("     Geth Version: "+geth.Version)
}
