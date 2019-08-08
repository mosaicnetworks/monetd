package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

//CopyToContainer copies a directory / file to a container
func CopyToContainer(cli *client.Client, containerID, localSrcPath, containerDestPath string) error {
	ctx := context.Background()

	archive, err := newTarArchiveFromPath(localSrcPath)
	if err != nil {
		return err
	}

	err = cli.CopyToContainer(ctx, containerID, containerDestPath, archive, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}

//CopyFromContainer copies a file / directory from a container
func CopyFromContainer(cli *client.Client, containerID, containerSrcPath, localDestPath string) error {

	ctx := context.Background()

	content, stat, err := cli.CopyFromContainer(ctx, containerID, containerSrcPath)
	if err != nil {
		return err
	}
	defer content.Close()

	srcInfo := archive.CopyInfo{
		Path:       containerSrcPath,
		Exists:     true,
		IsDir:      stat.Mode.IsDir(),
		RebaseName: "",
	}

	preArchive := content
	return archive.CopyTo(preArchive, srcInfo, localDestPath)
}

func newTarArchiveFromPath(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})

	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}
