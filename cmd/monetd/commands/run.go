package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/mosaicnetworks/monetd/src/babble"
<<<<<<< HEAD
	"github.com/mosaicnetworks/monetd/src/configuration"
=======
>>>>>>> origin/develop
	"github.com/spf13/cobra"
)

/*******************************************************************************
RunCmd
*******************************************************************************/

//newRunCmd returns the command that starts the daemon
func newRunCmd() *cobra.Command {
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
<<<<<<< HEAD
			_logger.WithField("Base", fmt.Sprintf("%+v", configuration.Configuration.BaseConfig)).Debug("Config Base")
			_logger.WithField("Babble", fmt.Sprintf("%+v", configuration.Configuration.Babble)).Debug("Config Babble")
			_logger.WithField("Eth", fmt.Sprintf("%+v", configuration.Configuration.Eth)).Debug("Config Eth")
=======
			_logger.WithField("Base", fmt.Sprintf("%+v", _config.BaseConfig)).Debug("Config Base")
			_logger.WithField("Babble", fmt.Sprintf("%+v", _config.Babble)).Debug("Config Babble")
			_logger.WithField("Eth", fmt.Sprintf("%+v", _config.Eth)).Debug("Config Eth")
>>>>>>> origin/develop

			return nil
		},

		RunE: runBabble,
	}

	bindFlags(cmd)

	return cmd
}

func bindFlags(cmd *cobra.Command) {
	// Babble config
	cmd.Flags().String("babble.listen", configuration.Configuration.Babble.BindAddr, "IP:PORT of Babble node")
	cmd.Flags().String("babble.service-listen", configuration.Configuration.Babble.ServiceAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Duration("babble.heartbeat", configuration.Configuration.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("babble.timeout", configuration.Configuration.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache-size", configuration.Configuration.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync-limit", configuration.Configuration.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("babble.max-pool", configuration.Configuration.Babble.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("babble.bootstrap", configuration.Configuration.Babble.Bootstrap, "Bootstrap Babble from database")

	// Eth config
	cmd.Flags().String("eth.listen", configuration.Configuration.Eth.EthAPIAddr, "IP:PORT of Monet HTTP API service")
	cmd.Flags().Int("eth.cache", configuration.Configuration.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")
}

/*******************************************************************************
READ CONFIG AND RUN
*******************************************************************************/

// Run the EVM-Lite / Babble engine
func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(configuration.Configuration.ToBabbleConfig(), _logger)
	engine, err := engine.NewEngine(*configuration.Configuration.ToEVMLConfig(), babble, _logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
