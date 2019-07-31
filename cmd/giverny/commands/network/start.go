package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
giverny network start

Starts a network. Does not start individual nodes
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}

	addStartFlags(cmd)

	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	network := args[0]

	return nil
}

func startDockerNetwork(networkName string) err {

	return nil
}

func startDocker() error {

	/*
	   #!/bin/bash

	   set -eux

	   N=${1:-4}
	   FASTSYNC=${2:-false}
	   MPWD=$(pwd)


	   docker network create \
	     --driver=bridge \
	     --subnet=172.77.0.0/16 \
	     --ip-range=172.77.0.0/16 \""
	     --gateway=172.77.5.254 \
	     babblenet

	   for i in $(seq 1 $N)
	   do
	       docker run -d --name=client$i --net=babblenet --ip=172.77.10.$i -it mosaicnetworks/dummy:0.5.0 \
	       --name="client $i" \
	       --client-listen="172.77.10.$i:1339" \
	       --proxy-connect="172.77.5.$i:1338" \
	       --discard \
	       --log="debug"
	   done

	   for i in $(seq 1 $N)
	   do
	       docker create --name=node$i --net=babblenet --ip=172.77.5.$i mosaicnetworks/babble:0.5.0 run \
	       --heartbeat=100ms \
	       --moniker="node$i" \
	       --cache-size=50000 \
	       --listen="172.77.5.$i:1337" \
	       --proxy-listen="172.77.5.$i:1338" \
	       --client-connect="172.77.10.$i:1339" \
	       --service-listen="172.77.5.$i:80" \
	       --sync-limit=500 \
	       --fast-sync=$FASTSYNC \
	       --store \
	       --log="debug"

	       docker cp $MPWD/conf/node$i node$i:/.babble
	       docker start node$i
	   done


	*/
	return nil
}
