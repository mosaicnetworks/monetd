package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ListImages lists docker images. It is a wrapper to GetImages
func ListImages(cli *client.Client) error {
	_, err := GetImages(cli, true)
	return err
}

// GetImages retrieves a list of tags from the local docker repository
func GetImages(cli *client.Client, showOutput bool) (map[string]string, error) {

	rtn := make(map[string]string)
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return rtn, err
	}

	for _, image := range images {
		if showOutput {
			fmt.Println("Image   : " + image.ID + " ")
			fmt.Println("          " + strings.Join(image.RepoTags, "\n          ") + "\n")
		}
		for _, tag := range image.RepoTags {
			rtn[tag] = image.ID
		}

	}

	return rtn, nil
}
