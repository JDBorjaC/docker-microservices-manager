package internal

import (
	"context"
	"log"
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

	// Write source code locally
	internalDir := filepath.Join("microservices", req.Name)
	os.MkdirAll(internalDir, 0755)
	os.WriteFile(
		filepath.Join(internalDir, "app.py"),
		[]byte(req.Code),
		0644,
	)

	result, err := s.client.CreateMicroservice(ctx, internalDir, Microservice{
		Name:        req.Name,
		Image:       req.Image,
		Description: req.Description,
	})

	log.Println(result)

	if err != nil {
		return err
	}
	//TODO: add microservice to repository
	return nil
}
