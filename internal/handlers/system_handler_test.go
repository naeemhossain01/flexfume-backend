package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestSystemHandler_GetSystemInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	systemService := services.NewSystemService()
	handler := NewSystemHandler(systemService)

	t.Run("Success - Get system info", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/system/info", nil)

		handler.GetSystemInfo(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
