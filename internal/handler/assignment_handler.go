package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"abplatform/internal/service"
)

type AssignmentHandler struct {
	service *service.AssignmentService
}

func NewAssignmentHandler(service *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

func (h *AssignmentHandler) Assign(c *gin.Context) {
	experimentIDStr := c.Query("experiment_id")
	userID := c.Query("user_id")

	if experimentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "experiment_id is required",
		})
		return
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	experimentID, err := strconv.Atoi(experimentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "experiment_id must be an integer",
		})
		return
	}

	assignment, err := h.service.AssignUser(c.Request.Context(), experimentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to assign user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, assignment)
}