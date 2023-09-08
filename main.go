package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	buildContext, buildContextErr := archive.TarWithOptions("./worker", &archive.TarOptions{})
    if buildContextErr != nil {
        panic(buildContextErr)
    }

	buildResponse, buildErr := cli.ImageBuild(
        ctx,
        buildContext,
        types.ImageBuildOptions{
            Dockerfile: "Dockerfile",
            Tags:       []string{"worker-image"},
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
