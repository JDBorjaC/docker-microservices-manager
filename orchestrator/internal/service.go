package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
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
func (s *Service) CreateMicroservice(ctx context.Context, req CreateMicroserviceRequest) (*Microservice, error) {

	// Capture duplicate name
	existingMs, err := s.repo.GetMicroserviceByName(req.Name)
	if err != nil {
		return nil, err // Unexpected DB error
	}
	if existingMs != nil {
		return nil, fmt.Errorf("microservice with name '%s' already exists", req.Name)
	}

	// Map language to docker image
	imageMap := map[string]string{
		LangFlask:   "msm-runner-flask",
		LangExpress: "msm-runner-express",
		LangGin:     "msm-runner-gin",
		LangCargo:   "msm-runner-cargo",
	}

	imageName, exists := imageMap[req.Language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", req.Language)
	}

	// Write source code locally
	internalDir := filepath.Join("microservices", req.Name)
	os.MkdirAll(internalDir, 0755)
	os.WriteFile(
		filepath.Join(internalDir, "app.py"),
		[]byte(req.Code),
		0644,
	)

	ms := &Microservice{
		Name:        req.Name,
		Description: req.Description,
		Image:       imageName,
		Status:      StatusCreated,
		CreatedAt:   time.Now(),
	}

	// Create Container
	result, err := s.client.CreateMicroservice(ctx, internalDir, *ms)
	if err != nil {
		return nil, err
	}

	// Save Container ID and Insert into DB
	ms.ContainerId = result.ID

	if err := s.repo.InsertMicroservice(ms); err != nil {
		// ROLLBACK (Docker)
		s.client.RemoveMicroservice(context.Background(), ms.ContainerId)
		return nil, err
	}

	return ms, nil
}

func (s *Service) StartAndStreamMicroservice(ctx context.Context, id int) (io.ReadCloser, error) {

	//Get Container ID from DB
	containerId, err := s.repo.GetMicroserviceContainerID(id)
	log.Printf("Container ID: %s", containerId)
	if err != nil {
		return nil, err
	}

	//Start Container
	err = s.client.StartMicroservice(ctx, containerId)
	if err != nil {
		s.repo.UpdateMicroserviceStatus(id, StatusFailed)
		return nil, err
	}

	//Update Status to Running
	s.repo.UpdateMicroserviceStatus(id, StatusRunning)

	//Stream Logs
	return s.client.LogMicroservice(ctx, containerId, true)
}

func (s *Service) GetAllMicroservices() ([]Microservice, error) {
	return s.repo.GetAllMicroservices()
}
