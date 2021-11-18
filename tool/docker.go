package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"

	v "github.com/spf13/viper"
)

func docker(img, dir string) error {
	if v.GetBool("no-docker") {
		return errDockerExecDisabled
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	fd, isTerminal := term.GetFdInfo(os.Stdout)
	err = jsonmessage.DisplayJSONMessagesStream(reader, os.Stdout, fd, isTerminal, nil)
	if err != nil {
		return err
	}

	adir, err := os.Getwd()
	if err != nil {
		return err
	}

	wrk := "/src"
	if v.GetBool("indocker") {
		wrk = dir
	}

	bind := path.Join(adir, dir) + ":/src"
	if v.GetBool("indocker") {
		bind = "issues:/volume"
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      img,
			Cmd:        []string{"./run"},
			WorkingDir: wrk,
			// TODO
			// it'd be interesting to set Tty: True;
			// unfortunately MultiWriter below fails because bytes.Buffer is not a TTY
			// however, os.Stdout and os.Stderr are TTYs
			Tty: true,
		},
		&container.HostConfig{
			Binds: []string{bind},
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
		fmt.Printf("status.StatusCode: %#+v\n", status.StatusCode)
		exitCode = int(status.StatusCode)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(io.MultiWriter(os.Stderr, buf), out)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	if exitCode != 0 {
		if strings.Contains(buf.String(), "exec user process caused \"exec format error\"") {
			return errExecFormat
		}
		return errExecFailure
	}
	return nil
}
