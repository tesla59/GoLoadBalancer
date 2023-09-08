package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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

	buildResponse, buildErr := cli.ImageBuild(
		ctx,
		buildContext,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{strings.ToLower(imageTag)},
		},
	)
	if buildErr != nil {
		panic(buildErr)
	}
	defer buildResponse.Body.Close()

	_, copyErr := io.Copy(os.Stdout, buildResponse.Body)
	if copyErr != nil {
		panic(copyErr)
	}
}
