package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {
	//Need both...this might make things interesting
	imageName := "alpine"
	imageCanconicalPath := "docker.io/library/" + imageName

	//This should work in theory for all of them....auto detecting which to run might be interesting...
	//or we cheat and run them all until one works :)
	packageManCommand := "apk info -vv | sort"

	//container name probably needs to be randomly generated and include ahab
	containerName := "ahab-testing"

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, imageCanconicalPath, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"sh", "-c", packageManCommand},
		Tty:   false,
	},
	nil,
	nil,
		containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	status, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		panic(err)
	}
}