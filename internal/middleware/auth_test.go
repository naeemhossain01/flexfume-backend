package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := auth.NewJWTManager("test-secret-key", 24*time.Hour)

	t.Run("Success - Valid token", func(t *testing.T) {
		token, err := jwtManager.GenerateToken("test-user-123", "+1234567890", "USER")
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		AuthMiddleware(jwtManager)(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, "test-user-123", c.GetString("user_id"))
		assert.Equal(t, "+1234567890", c.GetString("phone_number"))
		assert.Equal(t, "USER", c.GetString("user_role"))
	})

	t.Run("Error - Missing Authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		AuthMiddleware(jwtManager)(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Invalid token format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat")

		AuthMiddleware(jwtManager)(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		AuthMiddleware(jwtManager)(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Expired token", func(t *testing.T) {
		expiredManager := auth.NewJWTManager("test-secret-key", -1*time.Hour)
		token, err := expiredManager.GenerateToken("test-user-123", "+1234567890", "USER")
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		AuthMiddleware(jwtManager)(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
