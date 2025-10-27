package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler()

	t.Run("Success - Health check", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/health", nil)

		handler.HealthCheck(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
