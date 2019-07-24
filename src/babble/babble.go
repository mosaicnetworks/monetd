package babble

import (
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// InmemBabble implementes EVM-Lite's Consensus interface.
// It uses an inmemory Babble node.
type InmemBabble struct {
	config     *babble.BabbleConfig
	babble     *babble.Babble
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Logger
}

// NewInmemBabble instantiates a new InmemBabble consensus system
func NewInmemBabble(config *babble.BabbleConfig, logger *logrus.Logger) *InmemBabble {
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
	ib.logger.Debug("INIT")

	ib.ethState = state
	ib.ethService = service

	ib.config.Proxy = NewInmemProxy(state, service, service.GetSubmitCh(), ib.logger)

	babble := babble.NewBabble(ib.config)

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
