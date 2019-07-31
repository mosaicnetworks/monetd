package docker

import (
	"github.com/docker/docker/client"
)

// GetDockerClient returns a docker client
func GetDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
