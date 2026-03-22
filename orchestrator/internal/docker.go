package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/errdefs"
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
			NetworkMode: "msmanager-network",
		},
		Config: &container.Config{
			Tty: true,
			Labels: map[string]string{
				"traefik.enable": "true",
				"traefik.http.routers." + ms.Name + ".rule":                           "PathPrefix(`/services/" + ms.Name + "`)",
				"traefik.http.routers." + ms.Name + ".priority":                       "10",
				"traefik.http.services." + ms.Name + ".loadbalancer.server.port":      ms.Port,
				"traefik.http.middlewares." + ms.Name + "-strip.stripprefix.prefixes": "/services/" + ms.Name,
				"traefik.http.routers." + ms.Name + ".middlewares":                    ms.Name + "-strip",
			},
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
		if errdefs.IsNotFound(err) {
			return fmt.Errorf("%w: container with id %s", ErrNotFound, id)
		}
		return err
	}
	return nil
}

func (d *DockerClient) LogMicroservice(ctx context.Context, id string, follow bool) (io.ReadCloser, error) {
	out, err := d.client.ContainerLogs(ctx, id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
	})
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, fmt.Errorf("%w: container with id %s", ErrNotFound, id)
		}
		return nil, err
	}

	return out, nil
}

func (d *DockerClient) StopMicroservice(ctx context.Context, id string) error {
	_, err := d.client.ContainerStop(ctx, id, client.ContainerStopOptions{})
	if err != nil {
		if errdefs.IsNotFound(err) {
			return fmt.Errorf("%w: container with id %s", ErrNotFound, id)
		}
		return err
	}
	return nil
}

func (d *DockerClient) RemoveMicroservice(ctx context.Context, id string) error {
	_, err := d.client.ContainerRemove(ctx, id, client.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		if errdefs.IsNotFound(err) {
			return fmt.Errorf("%w: container with id %s", ErrNotFound, id)
		}
		return err
	}
	return nil
}
