package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"

	v "github.com/spf13/viper"
)

func checkInDocker() error {
	v.Set("indocker", false)
	if runtime.GOOS != "windows" {
		cmd := exec.Command("cat", "/proc/self/cgroup")
		o, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if strings.Contains(string(o), "docker") {
			fmt.Println("it seems you are running issue-runner inside a Docker container")
			if v.GetBool("no-docker") {
				fmt.Println("but execution of sibling containers is disabled through '--no-docker'")
			} else {
				sock := "/var/run/docker.sock"
				_, err := os.Stat(sock)
				if os.IsNotExist(err) {
					return fmt.Errorf("'%s' does not exist", sock)
				} else if err != nil {
					return err
				}

				fmt.Println("to run issues in sibling containers, ensure that 'issues:/volume' is bind")
				vol := "/volume"
				_, err = os.Stat(vol)
				if os.IsNotExist(err) {
					return fmt.Errorf("'%s' does not exist", vol)
				} else if err != nil {
					return err
				}

				v.Set("indocker", true)
			}
		}
	}
	return nil
}

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
