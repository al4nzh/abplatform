package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"abplatform/internal/service"
)

type ExperimentHandler struct {
	service *service.ExperimentService
}

func NewExperimentHandler(service *service.ExperimentService) *ExperimentHandler {
	return &ExperimentHandler{service: service}
}

type CreateExperimentRequest struct {
	Name string `json:"name"`
}

func (h *ExperimentHandler) CreateExperiment(c *gin.Context) {
	var req CreateExperimentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name is required",
		})
		return
	}

	exp, err := h.service.CreateExperiment(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create experiment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, exp)
}