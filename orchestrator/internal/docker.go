package internal

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerClient struct {
	client *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerClient{client: apiClient}, nil
}

func (d *DockerClient) Close() error {
	return d.client.Close()
}

// Pulls a Docker image from a registry. If the image already exists locally, skips the pull.
func (d *DockerClient) PullImage(ctx context.Context, imageId string) error {

	//B: check if the image exists in docker host before pull
	_, err := d.client.ImageInspect(ctx, imageId)
	if err == nil {
		return nil
	}

	reader, err := d.client.ImagePull(ctx, imageId, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) CreateMicroservice(ctx context.Context, dir string, ms Microservice) (*client.ContainerCreateResult, error) {

	//Use absolute path to ensure that the Bind Volume references
	//the right host OS directory (i.e. in case the app runs inside of a container)
	absPath := filepath.Join(os.Getenv("HOST_SOURCE_PATH"), dir)

	resp, err := d.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: ms.Image,
		Name:  ms.Name,
		HostConfig: &container.HostConfig{
			Binds: []string{
				absPath + ":/app",
			},
		},
		Config: &container.Config{
			Tty: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil

}

func (d *DockerClient) StartMicroservice(ctx context.Context, id string) error {
	_, err := d.client.ContainerStart(ctx, id, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) LogMicroservice(ctx context.Context, id string) (*client.ContainerLogsResult, error) {
	out, err := d.client.ContainerLogs(ctx, id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return nil, err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return &out, nil
}

func (d *DockerClient) StopMicroservice(ctx context.Context, id string) error {
	_, err := d.client.ContainerStop(ctx, id, client.ContainerStopOptions{})
	if err != nil {
		return err
	}
	return nil
}
