package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateMicroservice godoc
// @Summary Create and start a microservice container
// @Description Writes user code to disk and starts a runner container with a bind mount. Returns the assigned host port.
// @Tags microservices
// @Accept json
// @Produce json
// @Param request body CreateMicroserviceRequest true "Microservice definition"
// @Success 201 {object} Microservice
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices [post]
func (h *Handler) CreateMicroservice(c *gin.Context) {
	var req CreateMicroserviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ms, err := h.service.CreateMicroservice(c.Request.Context(), req)
	if err != nil {
		if fmt.Sprintf("microservice with name '%s' already exists", req.Name) == err.Error() || errors.Is(err, ErrDuplicate) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrUnsupportedLang) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrBuildFailed) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "image build failed: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ms)
}

// GetMicroservices godoc
// @Summary Get all microservices
// @Description Retrieves all microservices stored in the orchestrator database
// @Tags microservices
// @Produce json
// @Success 200 {array} Microservice
// @Failure 500 {object} map[string]string
// @Router /microservices [get]
func (h *Handler) GetMicroservices(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")

	microservices, err := h.service.GetAllMicroservices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, microservices)
}

// GetMicroserviceByID godoc
// @Summary Get a microservice by ID
// @Description Retrieves a microservice by its internal ID
// @Tags microservices
// @Produce json
// @Param id path int true "Microservice Internal ID"
// @Success 200 {object} Microservice
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/{id} [get]
func (h *Handler) GetMicroserviceByID(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	ms, err := h.service.GetMicroserviceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ms == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "microservice not found"})
		return
	}
	c.JSON(http.StatusOK, ms)
}

// GetMicroserviceLogs godoc
// @Summary Get static logs for a microservice
// @Description Fetches the recent log output of a microservice container.
// @Tags microservices
// @Produce text/plain
// @Param id path int true "Microservice Internal ID"
// @Success 200 {string} string "Container Logs"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/logs/{id} [get]
func (h *Handler) GetMicroserviceLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	containerId, err := h.service.ValidateMicroserviceContainerID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logs, err := h.service.GetMicroserviceLogs(c.Request.Context(), containerId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Container not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching logs: " + err.Error()})
		return
	}

	c.String(http.StatusOK, logs)
}

// StopMicroservice godoc
// @Summary Stop a microservice container
// @Description Stops a running microservice container and updates its status in the db.
// @Tags microservices
// @Produce json
// @Param id path int true "Microservice Internal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/stop/{id} [patch]
func (h *Handler) StopMicroservice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	err = h.service.StopMicroservice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stop signal sent"})
}

// StartMicroservice godoc
// @Summary Start a microservice container
// @Description Starts a stopped microservice container and updates its status in the db.
// @Tags microservices
// @Produce json
// @Param id path int true "Microservice Internal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/start/{id} [patch]
func (h *Handler) StartMicroservice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	err = h.service.StartMicroservice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "start signal sent"})
}

// StreamStatusUpdates godoc
// @Summary Stream all microservice status updates via SSE
// @Description Real-time stream of status changes for all managed microservices.
// @Tags microservices
// @Produce text/event-stream
// @Success 200 {string} string "SSE Stream of status updates"
// @Router /microservices/status/events [get]
func (h *Handler) StreamStatusUpdates(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported"})
		return
	}

	subChan := h.service.SubscribeStatus()
	defer h.service.UnsubscribeStatus(subChan)

	// Context for client disconnection
	clientCtx := c.Request.Context()

	// Notify connection established
	fmt.Fprintf(c.Writer, "event: info\ndata: Connected to status stream\n\n")
	flusher.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case ms, ok := <-subChan:
			if !ok {
				return
			}
			data, _ := json.Marshal(ms)
			fmt.Fprintf(c.Writer, "event: status_update\ndata: %s\n\n", string(data))
			flusher.Flush()
		case <-ticker.C:
			// Enviar keep-alive comment
			fmt.Fprintf(c.Writer, ": keep-alive\n\n")
			flusher.Flush()
		case <-clientCtx.Done():
			fmt.Printf("Status stream for a client disconnected: %v\n", clientCtx.Err())
			return
		}
	}
}

// UpdateMicroservice godoc
// @Summary Update a microservice's code
// @Description Rebuilds the microservice with new source code. Stops and removes the old container/image, builds a new image and creates a fresh container.
// @Tags microservices
// @Accept json
// @Produce json
// @Param id path int true "Microservice Internal ID"
// @Param request body UpdateMicroserviceRequest true "New code and optional description"
// @Success 200 {object} Microservice
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/{id} [put]
func (h *Handler) UpdateMicroservice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	var req UpdateMicroserviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ms, err := h.service.UpdateMicroservice(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "microservice not found"})
			return
		}
		if errors.Is(err, ErrUnsupportedLang) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrBuildFailed) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "image build failed: " + err.Error()})
			return
		}
		if errors.Is(err, ErrContainerFailed) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "container creation failed: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ms)
}

// DeleteMicroservice godoc
// @Summary Delete a microservice container and its data
// @Description Stops (if running), removes the container, deletes source code and db records.
// @Tags microservices
// @Produce json
// @Param id path int true "Microservice Internal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices/{id} [delete]
func (h *Handler) DeleteMicroservice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	err = h.service.RemoveMicroservice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "microservice deleted"})
}

// DeleteAllMicroservices godoc
// @Summary Delete all microservices
// @Description Stops and removes all microservice containers, custom images, and database entries.
// @Tags microservices
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices [delete]
func (h *Handler) DeleteAllMicroservices(c *gin.Context) {
	err := h.service.DeleteAllMicroservices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "some microservices failed to delete: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all microservices deleted successfully"})
}
