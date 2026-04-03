package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"abplatform/internal/kafka"
	"abplatform/internal/model"
)
type EventHandler struct {
	producer *kafka.Producer
}

func NewEventHandler(producer *kafka.Producer) *EventHandler {
	return &EventHandler{producer: producer}
}
func (h *EventHandler) TrackEvent(c *gin.Context) {
	var event model.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if event.ExperimentID == 0 || event.UserID == "" || event.EventName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing fields",
		})
		return
	}

	err := h.producer.SendEvent(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to send event",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "event sent",
	})
}