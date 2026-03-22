package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
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
		Port:        "5000",
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

func (s *Service) StartAndStreamMicroservice(ctx context.Context, id int, containerId string) (io.ReadCloser, error) {

	//Start Container
	err := s.client.StartMicroservice(ctx, containerId)
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
	ms, err := s.repo.GetAllMicroservices()
	if ms == nil {
		ms = []Microservice{}
	}
	return ms, err
}

func (s *Service) StopMicroservice(ctx context.Context, id int) error {
	containerId, err := s.repo.GetMicroserviceContainerID(id)
	if err != nil {
		return err
	}

	err = s.client.StopMicroservice(ctx, containerId)
	if err != nil {
		return err
	}

	err = s.repo.UpdateMicroserviceStatus(id, StatusStopped)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateMicroserviceContainerID(id int) (string, error) {
	containerId, err := s.repo.GetMicroserviceContainerID(id)
	if err != nil {
		return "", err
	}
	return containerId, nil
}

func (s *Service) RemoveMicroservice(ctx context.Context, id int) error {
	ms, err := s.repo.GetMicroserviceByID(id)
	if err != nil {
		return err
	}

	err = s.client.RemoveMicroservice(ctx, ms.ContainerId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	// Remove source code locally
	internalDir := filepath.Join("microservices", ms.Name)
	err = os.RemoveAll(internalDir)
	if err != nil {
		return err
	}

	err = s.repo.DeleteMicroservice(id)
	if err != nil {
		return err
	}

	return nil
}
