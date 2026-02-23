package internal

import (
	"context"
	"io"
	"os"

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
