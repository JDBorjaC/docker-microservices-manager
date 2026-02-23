package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type PullImageRequest struct {
	ImageId string `json:"imageId"`
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
