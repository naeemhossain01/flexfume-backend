package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// SystemHandler handles system-related endpoints
type SystemHandler struct {
	systemService *services.SystemService
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(systemService *services.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// WakeUp performs comprehensive system health checks
func (h *SystemHandler) WakeUp(c *gin.Context) {
	response := h.systemService.PerformWakeUpChecks()
	c.JSON(http.StatusOK, response)
}

// Ping returns a simple pong response
func (h *SystemHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "SUCCESS",
		Response: gin.H{
			"status":    "pong",
			"timestamp": time.Now().UnixMilli(),
			"message":   "Backend is awake",
		},
	})
}
