package main

import (
	"context"
	"fmt"
	"github.com/tesla59/goloadbalancer/loadbalancer"
	"io"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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
		hostPort, err := nat.NewPort(DefaultProtocol, fmt.Sprint(WorkerPort))
		if err != nil {
			panic(err)
		}
		targetPort, err := nat.NewPort(DefaultProtocol, fmt.Sprint(initialPort+i))
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
			Mounts: []mount.Mount{
				{
					Type:     mount.TypeBind,
					Source:   path.Join(PWD, DatabaseFileName),
					Target:   DatabaseTargetPath,
					ReadOnly: false,
				},
				{
					Type:     mount.TypeBind,
					Source:   path.Join(PWD, DefaultConfigFileName),
					Target:   ConfigTargetPath,
					ReadOnly: true,
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
		ContainerIDs = append(ContainerIDs, resp.ID)
	}
}

func StopContainers(IDs []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for i := range IDs {
		if err := cli.ContainerStop(ctx, IDs[i], container.StopOptions{}); err != nil {
			panic(err)
		}
	}
}

func RemoveContainers(IDs []string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for i := range IDs {
		if err := cli.ContainerRemove(ctx, IDs[i], types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
	}
}

func GetServerPool(n int) (Pool []string) {
	precedingURL := "http://localhost:"
	for i := 0; i < n; i++ {
		Pool = append(Pool, precedingURL+fmt.Sprint(config.WorkerPort+i))
	}
	return
}

func RunLoadBalancers(count int, initialPort int) {
	loadbalancer.InitLoadBalancer(GetServerPool(config.Worker))
	for i := 0; i < count; i++ {
		go loadbalancer.NewLoadBalancer(":" + fmt.Sprint(initialPort+i))
	}
}
