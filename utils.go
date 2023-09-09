package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

func BuildImage(srcPath string, imageTag string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	buildContext, buildContextErr := archive.TarWithOptions(srcPath, &archive.TarOptions{})
	if buildContextErr != nil {
		panic(buildContextErr)
	}

	imageBuildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{strings.ToLower(imageTag)},
	}

	buildResponse, buildErr := cli.ImageBuild(ctx, buildContext, imageBuildOptions)
	if buildErr != nil {
		panic(buildErr)
	}
	defer buildResponse.Body.Close()

	_, copyErr := io.Copy(os.Stdout, buildResponse.Body)
	if copyErr != nil {
		panic(copyErr)
	}
}

func RunImage(image string, count int, initialPort int) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// pulling image is disabled considering image will be built manually only
	for i := 0; i < count; i++ {
		hostPort, err := nat.NewPort("tcp", fmt.Sprint(initialPort))
		if err != nil {
			panic(err)
		}
		targetPort, err := nat.NewPort("tcp", fmt.Sprint(initialPort+i))
		if err != nil {
			panic(err)
		}
		containerConfig := container.Config{
			Image:        image,
			ExposedPorts: nat.PortSet{hostPort: struct{}{}},
		}

		hostConfig := container.HostConfig{
			PortBindings: nat.PortMap{
				hostPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: targetPort.Port(),
					},
				},
			},
		}

		resp, err := cli.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}
}
