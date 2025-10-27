package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "SUCCESS",
		Response: gin.H{
			"status":    "UP",
			"message":   "FlexFume Backend is running",
			"timestamp": time.Now().UnixMilli(),
		},
	})
}
