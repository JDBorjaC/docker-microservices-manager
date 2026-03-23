package internal

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
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

// BuildImage creates a tar archive from the microservice directory (containing
// the user code + generated Dockerfile) and sends it to the Docker daemon to
// build a custom image tagged as "msm-<name>:latest".
func (d *DockerClient) BuildImage(ctx context.Context, dir string, imageName string) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Walk the directory and add every file to the tar archive
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: relPath,
			Size: int64(len(data)),
			Mode: 0644,
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		_, err = tw.Write(data)
		return err
	})
	if err != nil {
		return fmt.Errorf("error creating tar context: %w", err)
	}
	tw.Close()

	resp, err := d.client.ImageBuild(ctx, buf, client.ImageBuildOptions{
		Tags:        []string{imageName},
		Remove:      true,
		ForceRemove: true,
	})
	if err != nil {
		return fmt.Errorf("error building image: %w", err)
	}
	defer resp.Body.Close()

	// Drain the build output to ensure the build completes and capture build errors
	decoder := json.NewDecoder(resp.Body)
	for {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding build stream: %w", err)
		}
		
		// If Docker daemon reports an error during build (e.g. missing base image), it comes here
		if errorMsg, ok := msg["error"]; ok {
			return fmt.Errorf("docker build failed: %v", errorMsg)
		}
	}

	return nil
}

func (d *DockerClient) CreateMicroservice(ctx context.Context, ms Microservice, port string) (*client.ContainerCreateResult, error) {

	resp, err := d.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: ms.Image,
		Name:  ms.Name,
		HostConfig: &container.HostConfig{
			NetworkMode: "msm-network",
		},
		Config: &container.Config{
			Tty: true,
			Labels: map[string]string{
				"traefik.enable": "true",
				"traefik.http.routers." + ms.Name + ".rule":                           "PathPrefix(`/services/" + ms.Name + "`)",
				"traefik.http.routers." + ms.Name + ".priority":                       "10",
				"traefik.http.services." + ms.Name + ".loadbalancer.server.port":      port,
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

// RemoveImage deletes a Docker image by name/tag.
func (d *DockerClient) RemoveImage(ctx context.Context, imageName string) error {
	_, err := d.client.ImageRemove(ctx, imageName, client.ImageRemoveOptions{Force: true, PruneChildren: true})
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil // Already gone, no problem
		}
		return err
	}
	return nil
}

func (d *DockerClient) IsContainerRunning(ctx context.Context, id string) (bool, error) {
	result, err := d.client.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return false, err
	}

	for _, c := range result.Items {
		if c.ID == id || (len(id) >= 12 && len(c.ID) >= len(id) && c.ID[:len(id)] == id) {
			return c.State == "running", nil
		}
	}

	return false, nil
}

func (d *DockerClient) WatchContainerEvents(ctx context.Context, onContainerDie func(containerID string)) {
	res := d.client.Events(ctx, client.EventsListOptions{})

	go func() {
		for {
			select {
			case err := <-res.Err:
				fmt.Printf("Error listener docker eventos: %v\n", err)
				return
			case msg := <-res.Messages:
				if msg.Type == "container" && (msg.Action == "die" || msg.Action == "destroy") {
					onContainerDie(msg.Actor.ID)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
