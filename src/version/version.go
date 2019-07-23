package version

import (
	"fmt"

	geth "github.com/ethereum/go-ethereum/params"
	_babble "github.com/mosaicnetworks/babble/src/version"
	evm "github.com/mosaicnetworks/evm-lite/src/version"
)

//Maj is Major Version Number
const Maj = "0"

//Min is Minor Version Number
const Min = "1"

//Fix is the Patch Version
const Fix = "0"

var (
	//Version is the full version string
	Version = "0.1.0"

	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	// GitBranch is set with --ldflags "-X main.gitBranch=$(git symbolic-ref --short HEAD)"
	GitBranch string
)

func init() {
	// branch is only of interest if it is not the master branch
	if GitBranch != "" && GitBranch != "master" {
		Version += "-" + GitBranch
	}

	if GitCommit != "" {
		Version += "-" + GitCommit[:8]
	}
}

//FullVersion outputs version information
func FullVersion() string {
	return fmt.Sprintln("Monet Version: "+Version) +
		fmt.Sprintln("     EVM-Lite Version: "+evm.Version) +
		fmt.Sprintln("     Babble Version: "+_babble.Version) +
		fmt.Sprintln("     Geth Version: "+geth.Version)
}
