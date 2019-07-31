package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

//GetNetworks lists networks
func GetNetworks(cli *client.Client, output bool) (map[string]string, error) {

	rtn := make(map[string]string)

	ctx := context.Background()
	arrRes, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return rtn, err
	}

	for _, net := range arrRes {
		if output {
			fmt.Printf("%s   %s  %s\n", net.Name, net.ID, net.Driver)
		}
		rtn[net.Name] = net.ID
	}

	return rtn, nil
}

// CreateNetwork creates a docker network
func CreateNetwork(cli *client.Client, networkName, subnet, iprange, gateway string) (string, error) {
	ctx := context.Background()

	fmt.Println("Creating options")
	opts := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		IPAM: &network.IPAM{
			//			Driver: "bridge",
			Config: []network.IPAMConfig{
				network.IPAMConfig{
					Subnet:  subnet,
					IPRange: iprange,
					Gateway: gateway,
				},
			},
		},
	}

	fmt.Println("Creating network")
	netresp, err := cli.NetworkCreate(ctx, networkName, opts)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ID: " + netresp.ID)
	if netresp.Warning != "" {
		fmt.Println("Warning: " + netresp.Warning)
	}

	return netresp.ID, err
}

// RemoveNetwork removes a network
func RemoveNetwork(cli *client.Client, networkID string) error {
	ctx := context.Background()
	return cli.NetworkRemove(ctx, networkID)
}
