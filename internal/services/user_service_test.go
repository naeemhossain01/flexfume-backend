package services

import (
	"testing"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// setupUserTestDB creates an in-memory SQLite database for testing
func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the global DB for testing
	database.DB = db

	return db
}

// teardownUserTestDB closes the database connection
func teardownUserTestDB() {
	database.DB = nil
}

func TestUserService_CreateUser(t *testing.T) {
	db := setupUserTestDB(t)
	defer teardownUserTestDB()

	service := NewUserService()

	t.Run("Success - Create user", func(t *testing.T) {
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}

		result, err := service.CreateUser(user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("Error - Nil user", func(t *testing.T) {
		result, err := service.CreateUser(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	db := setupUserTestDB(t)
	defer teardownUserTestDB()

	service := NewUserService()

	t.Run("Success - Get user by ID", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		result, err := service.GetUserByID(user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		result, err := service.GetUserByID("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	db := setupUserTestDB(t)
	defer teardownUserTestDB()

	service := NewUserService()

	t.Run("Success - Get user by email", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		result, err := service.GetUserByEmail("test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		result, err := service.GetUserByEmail("nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
