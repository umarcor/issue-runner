package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func docker(img, dir string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}

	adir, err := os.Getwd()
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      img,
			Cmd:        []string{"./run"},
			WorkingDir: "/src",
			Tty:        true,
		},
		&container.HostConfig{
			Binds: []string{path.Join(adir, dir) + ":/src"},
		},
		nil,
		"",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	exitCode := 0
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		log.Printf("status.StatusCode: %#+v\n", status.StatusCode)
		exitCode = int(status.StatusCode)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("container exit %d", exitCode)
	}
	return nil
}
