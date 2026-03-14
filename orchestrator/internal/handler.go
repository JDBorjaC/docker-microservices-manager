package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// PullImage godoc
// @Summary Pull a Docker image
// @Description Pulls a Docker image from a registry. If the image already exists locally, skips the pull.
// @Tags images
// @Accept json
// @Produce json
// @Param request body PullImageRequest true "Image to pull"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /images/pull [post]
func (h *Handler) PullImage(c *gin.Context) {
	var req PullImageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PullImage(c.Request.Context(), req.ImageId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success!!!"})
}

// CreateMicroservice godoc
// @Summary Create and start a microservice container
// @Description Writes user code to disk and starts a runner container with a bind mount. Returns the assigned host port.
// @Tags microservices
// @Accept json
// @Produce json
// @Param request body CreateMicroserviceRequest true "Microservice definition"
// @Success 201 {object} MicroserviceResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /microservices [post]
func (h *Handler) CreateMicroservice(c *gin.Context) {
	var req CreateMicroserviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateMicroservice(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "yipeee")
}
