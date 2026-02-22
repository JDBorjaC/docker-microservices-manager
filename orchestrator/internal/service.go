package internal

import "context"

type Service struct {
	client *DockerClient
}

func NewService(client *DockerClient) *Service {
	return &Service{client: client}
}

func (s *Service) PullImage(ctx context.Context, imageId string) error {
	return s.client.PullImage(ctx, imageId)
}
