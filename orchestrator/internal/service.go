package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

type Service struct {
	client *DockerClient
	repo   *Repository

	statusListeners   map[chan Microservice]struct{}
	statusListenersMu sync.RWMutex
}

func NewService(client *DockerClient, repo *Repository) *Service {
	return &Service{
		client:          client,
		repo:            repo,
		statusListeners: make(map[chan Microservice]struct{}),
	}
}

type langConfig struct {
	BaseImage    string
	Dependencies string
	SourceFile   string
	DestPath     string
	InternalPort string
	SyntaxCheck  string
	RunCmd       string
}

var supportedLanguages = map[string]langConfig{
	LangFlask: {
		BaseImage:    "python:3.9-slim",
		Dependencies: "RUN pip install flask",
		SourceFile:   "app.py",
		DestPath:     "/app/app.py",
		InternalPort: "5000",
		SyntaxCheck:  "python -m py_compile /app/app.py",
		RunCmd:       `CMD ["python", "app.py"]`,
	},
	LangExpress: {
		BaseImage:    "node:22-alpine",
		Dependencies: "RUN npm init -y && npm install express",
		SourceFile:   "app.js",
		DestPath:     "/app/app.js",
		InternalPort: "3000",
		SyntaxCheck:  "node -c /app/app.js",
		RunCmd:       `CMD ["node", "app.js"]`,
	},
	LangGin: {
		BaseImage:    "golang:1.25.6",
		Dependencies: "RUN go mod init runner && go get github.com/gin-gonic/gin",
		SourceFile:   "main.go",
		DestPath:     "/app/main.go",
		InternalPort: "8080",
		SyntaxCheck:  "go build -o /dev/null /app/main.go",
		RunCmd:       `CMD ["go", "run", "main.go"]`,
	},
	LangCargo: {
		BaseImage:    "rust:1.85-slim",
		Dependencies: "RUN cargo init --name runner . && echo 'actix-web = \"4\"' >> Cargo.toml && echo 'serde = { version = \"1\", features = [\"derive\"] }' >> Cargo.toml && echo 'serde_json = \"1\"' >> Cargo.toml && cargo build --release && rm src/main.rs",
		SourceFile:   "main.rs",
		DestPath:     "/app/src/main.rs",
		InternalPort: "8080",
		SyntaxCheck:  "cargo build --release",
		RunCmd:       `CMD ["cargo", "run", "--release"]`,
	},
}

func (s *Service) CreateMicroservice(ctx context.Context, req CreateMicroserviceRequest) (*Microservice, error) {

	// Capture duplicate name
	existingMs, err := s.repo.GetMicroserviceBy("name", req.Name)
	if err != nil {
		return nil, err // Unexpected DB error
	}
	if existingMs != nil {
		return nil, fmt.Errorf("microservice with name '%s' already exists", req.Name)
	}

	// Resolve language config
	lang, exists := supportedLanguages[req.Language]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedLang, req.Language)
	}

	// Write source code + Dockerfile locally
	buildDir := filepath.Join("microservices", req.Name)
	os.MkdirAll(buildDir, 0755)
	defer os.RemoveAll(buildDir)

	os.WriteFile(
		filepath.Join(buildDir, lang.SourceFile),
		[]byte(req.Code),
		0644,
	)

	dockerfile := fmt.Sprintf("FROM %s\nWORKDIR /app\nRUN (useradd -m runner || adduser -D runner) || true\n%s\nCOPY %s %s\nRUN chown -R runner /app\nUSER runner\nRUN %s\n%s\n",
		lang.BaseImage, lang.Dependencies, lang.SourceFile, lang.DestPath, lang.SyntaxCheck, lang.RunCmd)
	os.WriteFile(
		filepath.Join(buildDir, "Dockerfile"),
		[]byte(dockerfile),
		0644,
	)

	// Build the custom image
	customImageName := "msm-" + req.Name + ":latest"
	if err := s.client.BuildImage(ctx, buildDir, customImageName); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildFailed, err)
	}

	ms := &Microservice{
		Name:        req.Name,
		Description: req.Description,
		Language:    req.Language,
		Image:       customImageName,
		Status:      ContainerCreated,
		Code:        req.Code,
		CreatedAt:   time.Now(),
	}

	// Create Container from the custom-built image
	result, err := s.client.CreateMicroservice(ctx, *ms, lang.InternalPort)
	if err != nil {
		// ROLLBACK: remove the built image
		s.client.RemoveImage(context.Background(), customImageName)
		return nil, fmt.Errorf("%w: %v", ErrContainerFailed, err)
	}

	// Save Container ID and Insert into DB
	ms.ContainerId = result.ID

	if err := s.repo.InsertMicroservice(ms); err != nil {
		// ROLLBACK (Docker container + image)
		s.client.RemoveMicroservice(context.Background(), ms.ContainerId)
		s.client.RemoveImage(context.Background(), customImageName)
		return nil, err
	}

	return ms, nil
}

func (s *Service) GetMicroserviceLogs(ctx context.Context, containerId string) (string, error) {
	stream, err := s.client.LogMicroservice(ctx, containerId, false)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	bytesLog, err := io.ReadAll(stream)
	if err != nil {
		return "", fmt.Errorf("could not read log stream: %w", err)
	}

	return string(bytesLog), nil
}

func (s *Service) StartMicroservice(ctx context.Context, id int) error {
	ms, err := s.repo.GetMicroserviceBy("id", id)
	if err != nil {
		return err
	}
	if ms == nil {
		return ErrNotFound
	}

	return s.client.StartMicroservice(ctx, ms.ContainerId)
}

func (s *Service) GetAllMicroservices() ([]Microservice, error) {
	ms, err := s.repo.GetAllMicroservices()
	if ms == nil {
		ms = []Microservice{}
	}
	return ms, err
}

func (s *Service) GetMicroserviceByID(ctx context.Context, id int) (*Microservice, error) {
	return s.repo.GetMicroserviceBy("id", id)
}

func (s *Service) StopMicroservice(ctx context.Context, id int) error {
	ms, err := s.repo.GetMicroserviceBy("id", id)
	if err != nil {
		return err
	}
	if ms == nil {
		return ErrNotFound
	}

	// Marcar como exited previniendo el check de OnContainerDie
	s.repo.UpdateMicroservice(id, map[string]any{"status": ContainerExited})

	err = s.client.StopMicroservice(ctx, ms.ContainerId)
	if err != nil {
		s.repo.UpdateMicroservice(id, map[string]any{"status": ms.Status})
		return err
	}

	ms.Status = ContainerExited
	s.broadcastStatus(*ms)

	return nil
}

func (s *Service) ValidateMicroserviceContainerID(id int) (string, error) {
	ms, err := s.repo.GetMicroserviceBy("id", id)
	if err != nil {
		return "", err
	}
	if ms == nil {
		return "", ErrNotFound
	}
	return ms.ContainerId, nil
}

func (s *Service) RemoveMicroservice(ctx context.Context, id int) error {
	ms, err := s.repo.GetMicroserviceBy("id", id)
	if err != nil {
		return err
	}
	if ms == nil {
		return ErrNotFound
	}

	err = s.client.RemoveMicroservice(ctx, ms.ContainerId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	// Remove the custom-built Docker image
	err = s.client.RemoveImage(ctx, ms.Image)
	if err != nil {
		fmt.Printf("Warning: could not remove image %s: %v\n", ms.Image, err)
	}

	err = s.repo.DeleteMicroservice(id)
	if err != nil {
		return err
	}

	ms.Status = ContainerRemoved
	s.broadcastStatus(*ms)

	return nil
}

func (s *Service) DeleteAllMicroservices(ctx context.Context) error {
	microservices, err := s.GetAllMicroservices()
	if err != nil {
		return err
	}

	var errs []string
	for _, ms := range microservices {
		if err := s.RemoveMicroservice(ctx, ms.Id); err != nil {
			errs = append(errs, fmt.Sprintf("failed to delete MS %d: %v", ms.Id, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during deletion: %s", strings.Join(errs, ", "))
	}
	return nil
}

func (s *Service) UpdateMicroservice(ctx context.Context, id int, req UpdateMicroserviceRequest) (*Microservice, error) {
	// Fetch existing microservice
	ms, err := s.repo.GetMicroserviceBy("id", id)
	if err != nil {
		return nil, err
	}
	if ms == nil {
		return nil, ErrNotFound
	}

	// Update description if provided
	if req.Description != nil {
		ms.Description = *req.Description
		if err := s.repo.UpdateMicroservice(id, map[string]any{"description": ms.Description}); err != nil {
			return nil, err
		}
	}

	// If no code is updated, updated
	if req.Code == nil {
		return ms, nil
	}

	// Resolve final Language to use
	lang, exists := supportedLanguages[ms.Language]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedLang, ms.Language)
	}

	// We must have code to rebuild.
	codeToInject := *req.Code

	// Write new files & Build BEFORE destroying old container
	buildDir := filepath.Join("microservices", ms.Name)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, err
	}
	defer os.RemoveAll(buildDir)

	if err := os.WriteFile(filepath.Join(buildDir, lang.SourceFile), []byte(codeToInject), 0644); err != nil {
		return nil, err
	}

	dockerfile := fmt.Sprintf("FROM %s\nWORKDIR /app\nRUN (useradd -m runner || adduser -D runner) || true\n%s\nCOPY %s %s\nRUN chown -R runner /app\nUSER runner\nRUN %s\n%s\n",
		lang.BaseImage, lang.Dependencies, lang.SourceFile, lang.DestPath, lang.SyntaxCheck, lang.RunCmd)
	if err := os.WriteFile(filepath.Join(buildDir, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
		return nil, err
	}

	// Build new image (Docker will strip the tag from the old image if successful)
	customImageName := "msm-" + ms.Name + ":latest"
	if err := s.client.BuildImage(ctx, buildDir, customImageName); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildFailed, err)
	}

	// Set status to updating to decouple lifecycle events
	s.repo.UpdateMicroservice(id, map[string]any{"status": ContainerUpdating})
	ms.Status = ContainerUpdating

	// Now it's safe to destroy old container
	s.client.StopMicroservice(ctx, ms.ContainerId)
	s.client.RemoveMicroservice(ctx, ms.ContainerId)

	// Create new container
	ms.Image = customImageName
	ms.Status = ContainerCreated
	result, err := s.client.CreateMicroservice(ctx, *ms, lang.InternalPort)
	if err != nil {
		s.client.RemoveImage(context.Background(), customImageName)
		return nil, fmt.Errorf("%w: %v", ErrContainerFailed, err)
	}

	// Update repository
	ms.ContainerId = result.ID
	ms.Image = customImageName
	ms.Status = ContainerCreated
	ms.Code = codeToInject

	updates := map[string]any{
		"container_id": ms.ContainerId,
		"image":        ms.Image,
		"status":       ContainerCreated,
		"code":         codeToInject,
	}
	if err := s.repo.UpdateMicroservice(id, updates); err != nil {
		return nil, err
	}

	s.broadcastStatus(*ms)

	return ms, nil
}

/*
SYNCHRONIZATION OF THE ORCHESTRATOR WITH THE ACTUAL STATE OF THE CONTAINERS
*/

func (s *Service) SyncStateOnStartup(ctx context.Context) error {
	microservices, err := s.repo.GetAllMicroservices()
	if err != nil {
		return err
	}

	for _, ms := range microservices {
		if ms.ContainerId == "" {
			continue
		}
		state, err := s.client.GetContainerState(ctx, ms.ContainerId)
		if err != nil {
			fmt.Printf("Error verificando contenedor %s: %v\n", ms.Name, err)
			continue
		}
		internalState := s.mapDockerStateToInternal(state)
		if internalState != ms.Status {
			s.repo.UpdateMicroservice(ms.Id, map[string]any{"status": internalState})
		}
	}
	return nil
}

func (s *Service) StartReconciliationLoop(ctx context.Context) {
	eventHandler := ContainerEventHandler{
		OnContainerStart: func(containerID string) {
			ms, err := s.repo.GetMicroserviceBy("container_id", containerID)
			if err != nil || ms == nil {
				return
			}
			s.repo.UpdateMicroservice(ms.Id, map[string]any{
				"status": ContainerRunning,
			})
			ms.Status = ContainerRunning
			s.broadcastStatus(*ms)
		},
		OnContainerDie: func(containerID string) {
			ms, _ := s.repo.GetMicroserviceBy("container_id", containerID)
			if ms == nil {
				return
			}
			if ms.Status == ContainerUpdating || ms.Status == ContainerExited {
				return
			}

			state, err := s.client.GetContainerState(ctx, containerID)
			if err != nil {
				return
			}

			newStatus := s.mapDockerStateToInternal(state)

			s.repo.UpdateMicroservice(ms.Id, map[string]any{
				"status": newStatus,
			})
			ms.Status = newStatus
			s.broadcastStatus(*ms)
		},
		OnContainerDestroy: func(containerID string) {
			ms, err := s.repo.GetMicroserviceBy("container_id", containerID)
			if err != nil || ms == nil {
				return
			}
			if ms.Status == ContainerUpdating {
				return
			}
			s.RemoveMicroservice(context.Background(), ms.Id)
		},
	}

	s.client.WatchContainerEvents(ctx, eventHandler)
}

func (s *Service) mapDockerStateToInternal(state *types.ContainerState) string {
	if state.Status == "exited" {
		if state.ExitCode == 0 {
			return ContainerExited
		} else {
			return ContainerCrashed
		}
	}
	switch state.Status {
	case "created":
		return ContainerCreated
	case "running":
		return ContainerRunning
	case "paused":
		return ContainerPaused
	case "restarting":
		return ContainerRestarting
	case "removing":
		return ContainerRemoving
	case "dead":
		return ContainerDead
	default:
		return ContainerCrashed
	}
}

// Subscription logic for real-time status updates
func (s *Service) SubscribeStatus() chan Microservice {
	s.statusListenersMu.Lock()
	defer s.statusListenersMu.Unlock()

	ch := make(chan Microservice, 10)
	s.statusListeners[ch] = struct{}{}
	return ch
}

func (s *Service) UnsubscribeStatus(ch chan Microservice) {
	s.statusListenersMu.Lock()
	defer s.statusListenersMu.Unlock()

	delete(s.statusListeners, ch)
	close(ch)
}

func (s *Service) broadcastStatus(ms Microservice) {
	s.statusListenersMu.RLock()
	defer s.statusListenersMu.RUnlock()

	for ch := range s.statusListeners {
		select {
		case ch <- ms:
		default:
			// Buffer full, skip this update for this listener to avoid blocking
		}
	}
}
