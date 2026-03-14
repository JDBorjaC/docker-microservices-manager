package internal

import (
	"context"
	"os"
	"path/filepath"
)

type Service struct {
	client *DockerClient
	repo   *Repository
}

func NewService(client *DockerClient, repo *Repository) *Service {
	return &Service{client: client, repo: repo}
}
func (s *Service) PullImage(ctx context.Context, imageId string) error {
	return s.client.PullImage(ctx, imageId)
}
func (s *Service) CreateMicroservice(ctx context.Context, req CreateMicroserviceRequest) error {
	//TODO: prevent folder overwriting for an existing name
	dir := filepath.Join("microservices", req.Name)
	os.MkdirAll(dir, 0755)
	os.WriteFile(
		filepath.Join(dir, "main.py"),
		[]byte(req.Code),
		0644,
	)
	s.client.CreateMicroservice(ctx, dir, Microservice{
		Name:        req.Name,
		Image:       req.Image,
		Description: req.Description,
	})
	//TODO: add microservice to repository
	return nil
}
