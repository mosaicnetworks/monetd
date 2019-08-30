// Package babble wraps the babble and EVM-Lite components.

package babble

import (
	"github.com/mosaicnetworks/babble/src/babble"
	babble_config "github.com/mosaicnetworks/babble/src/config"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// InmemBabble implementes EVM-Lite's Consensus interface.
// It uses an inmemory Babble node.
type InmemBabble struct {
	config     *babble_config.Config
	babble     *babble.Babble
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Entry
}

// NewInmemBabble instantiates a new InmemBabble consensus system
func NewInmemBabble(config *babble_config.Config, logger *logrus.Entry) *InmemBabble {
	return &InmemBabble{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a Babble inmemory node
func (ib *InmemBabble) Init(state *state.State, service *service.Service) error {
	ib.ethState = state
	ib.ethService = service

	babble := babble.NewBabble(ib.config)

	inmemProxy := NewInmemProxy(state,
		service,
		babble,
		service.GetSubmitCh(),
		ib.logger)

	ib.config.Proxy = inmemProxy

	err := babble.Init()
	if err != nil {
		return err
	}

	ib.babble = babble

	return nil
}

// Run starts the Babble node
func (ib *InmemBabble) Run() error {
	ib.babble.Run()
	return nil
}

// Info returns Babble stats
func (ib *InmemBabble) Info() (map[string]string, error) {
	info := ib.babble.Node.GetStats()
	info["type"] = "babble"
	return info, nil
}
