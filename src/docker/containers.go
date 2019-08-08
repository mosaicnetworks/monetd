package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"

	"github.com/docker/docker/client"
)

//CreateContainerFromImage creates a container, returning its ID
func CreateContainerFromImage(cli *client.Client, imageName string, isImageRemote bool,
	nodeName string, cmd strslice.StrSlice, start bool) (string, error) {

	ctx := context.Background()

	if isImageRemote {
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		io.Copy(os.Stdout, out)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Cmd:   cmd,
		Image: imageName,
	}, nil, nil, nodeName)
	if err != nil {
		return "", err
	}

	if start {
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return "", err
		}
	}
	//	fmt.Println(resp.ID)
	return resp.ID, nil
}

//StartContainer starts a container previously created by CreateContainerFromImage
func StartContainer(cli *client.Client, containerID string) error {

	ctx := context.Background()
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

//GetContainers lists containers.
func GetContainers(cli *client.Client, output bool) (map[string]string, error) {

	rtn := make(map[string]string)

	ctx := context.Background()
	arrRes, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		fmt.Println(err.Error())
		return rtn, err
	}

	for _, net := range arrRes {
		if output {
			fmt.Printf("%s   %s  %s\n", net.Names[0], net.ID, net.Status)
		}
		if len(net.Names) > 0 {
			rtn[strings.TrimLeft(net.Names[0], "/")] = net.ID
		}
	}

	return rtn, nil
}

// StopContainer stops a container
func StopContainer(cli *client.Client, containerID string) error {
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		return err
	}
	return nil
}

// RemoveContainer removes a container
func RemoveContainer(cli *client.Client, containerID string, force, removelinks, removevolumes bool) error {
	ctx := context.Background()
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force:         force,
		RemoveLinks:   removelinks,
		RemoveVolumes: removevolumes,
	}); err != nil {
		return err
	}
	return nil
}

// ConnectContainerToNetwork connects a created container to an extant network
func ConnectContainerToNetwork(cli *client.Client, networkID string, containerID string, ip string) error {

	ctx := context.Background()

	return cli.NetworkConnect(ctx, networkID, containerID,
		&network.EndpointSettings{
			IPAMConfig: &network.EndpointIPAMConfig{IPv4Address: ip},
			IPAddress:  ip,
		})
}

