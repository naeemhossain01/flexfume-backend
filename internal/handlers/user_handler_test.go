package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupUserHandlerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownUserHandlerTestDB() {
	database.DB = nil
}

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupUserHandlerTestDB(t)
	defer teardownUserHandlerTestDB()

	userService := services.NewUserService()
	handler := NewUserHandler(userService)

	t.Run("Success - Get user by ID", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: user.ID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/user/"+user.ID, nil)

		handler.GetUserByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/user/non-existent-id", nil)

		handler.GetUserByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserHandler_GetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupUserHandlerTestDB(t)
	defer teardownUserHandlerTestDB()

	userService := services.NewUserService()
	handler := NewUserHandler(userService)

	t.Run("Success - Get current user", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", user.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/user/me", nil)

		handler.GetCurrentUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - Missing user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/user/me", nil)

		handler.GetCurrentUser(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
