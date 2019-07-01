package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/src/consensus/babble"
	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/spf13/cobra"
)

/*******************************************************************************
RunCmd
*******************************************************************************/

//NewRunCmd returns the command that starts the daemon
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a MONET node",
		Long: `Run a MONET node.
	
	  Start a daemon which acts as a full node on a MONET network. All data and 
	  configuration are stored under a directory [datadir] controlled by the 
	  --datadir flag ($HOME/.monet by default on UNIX systems). 
	  
	  [datadir] must contain a set of files defining the network that this node is 
	  attempting to join or create. Please refer to monetcli for a tool to manage 
	  this configuration. 
	  
	  Further options pertaining to the operation of the node are read from the 
	  [datadir]/monetd.toml file, or overwritten by the following flags.`,

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			logger.WithField("Base", fmt.Sprintf("%+v", config.BaseConfig)).Debug("Config Base")
			logger.WithField("Babble", fmt.Sprintf("%+v", config.Babble)).Debug("Config Babble")
			logger.WithField("Eth", fmt.Sprintf("%+v", config.Eth)).Debug("Config Eth")

			return nil
		},

		RunE: runBabble,
	}

	bindFlags(cmd)

	return cmd
}

func bindFlags(cmd *cobra.Command) {
	// Babble config
	cmd.Flags().String("babble.listen", config.Babble.BindAddr, "IP:PORT of Babble node")
	cmd.Flags().String("babble.service-listen", config.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Duration("babble.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("babble.timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache-size", config.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync-limit", config.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("babble.max-pool", config.Babble.MaxPool, "Max number of pool connections")

	// Eth config
	cmd.Flags().String("eth.listen", config.Eth.EthAPIAddr, "IP:PORT of Monet HTTP API service")
	cmd.Flags().Int("eth.cache", config.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")
}

/*******************************************************************************
READ CONFIG AND RUN
*******************************************************************************/

// Run the EVM-Lite / Babble engine
func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(config.Babble, logger)
	engine, err := engine.NewEngine(*config, babble, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
