package docker

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ShowLogs shows a docker container logs
func ShowLogs(client *client.Client, containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reader, err := client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: false})
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		return err
	}

	reader.Close()
	return nil
}
