package internal

import (
	"context"
	"io"
	"os"

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
	resp, err := d.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: ms.Image,
		Name:  ms.Name,
		HostConfig: &container.HostConfig{
			Binds: []string{
				dir + ":/app",
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

	/*
		if _, err := d.client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		wait := d.client.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{})
		select {
		case err := <-wait.Error:
			if err != nil {
				panic(err)
			}
		case <-wait.Result:
		}

		out, err := d.client.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	*/
}
