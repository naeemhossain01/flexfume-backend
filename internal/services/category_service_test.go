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

// setupCategoryTestDB creates an in-memory SQLite database for testing
func setupCategoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Category{}, &models.Product{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the global DB for testing
	database.DB = db

	return db
}

// teardownCategoryTestDB closes the database connection
func teardownCategoryTestDB() {
	database.DB = nil
}

// createTestCategory creates a test category in the database
func createTestCategory(t *testing.T, db *gorm.DB, name, description string) *models.Category {
	category := &models.Category{
		Name:        name,
		Description: description,
	}

	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	return category
}

func TestCategoryService_CreateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	service := NewCategoryService()

	t.Run("Success - Create valid category", func(t *testing.T) {
		category := &models.Category{
			Name:        "Perfumes",
			Description: "Luxury perfumes collection",
		}

		result, err := service.CreateCategory(category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "Perfumes", result.Name)
		assert.Equal(t, "Luxury perfumes collection", result.Description)
	})

	t.Run("Error - Nil category", func(t *testing.T) {
		result, err := service.CreateCategory(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("Error - Empty name", func(t *testing.T) {
		category := &models.Category{
			Name:        "",
			Description: "Test description",
		}

		result, err := service.CreateCategory(category)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCategoryNameRequired, err)
	})

	t.Run("Error - Duplicate name", func(t *testing.T) {
		// Create first category
		createTestCategory(t, db, "Electronics", "Electronic items")

		// Try to create duplicate
		category := &models.Category{
			Name:        "Electronics",
			Description: "Another electronics category",
		}

		result, err := service.CreateCategory(category)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCategoryAlreadyExists, err)
	})

	t.Run("Success - Name with whitespace trimmed", func(t *testing.T) {
		category := &models.Category{
			Name:        "  Clothing  ",
			Description: "Fashion items",
		}

		result, err := service.CreateCategory(category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Clothing", result.Name)
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	service := NewCategoryService()

	t.Run("Success - Update category name", func(t *testing.T) {
		existing := createTestCategory(t, db, "Books", "Book collection")

		updated := &models.Category{
			Name:        "Books & Magazines",
			Description: "Updated description",
		}

		result, err := service.UpdateCategory(existing.ID, updated)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Books & Magazines", result.Name)
		assert.Equal(t, "Updated description", result.Description)
	})

	t.Run("Error - Empty category ID", func(t *testing.T) {
		updated := &models.Category{
			Name: "Test",
		}

		result, err := service.UpdateCategory("", updated)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Error - Nil category data", func(t *testing.T) {
		result, err := service.UpdateCategory("some-id", nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		updated := &models.Category{
			Name: "Test",
		}

		result, err := service.UpdateCategory("non-existent-id", updated)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCategoryNotFound, err)
	})

	t.Run("Error - Duplicate name on update", func(t *testing.T) {
		cat1 := createTestCategory(t, db, "Sports", "Sports items")
		createTestCategory(t, db, "Fitness", "Fitness equipment")

		updated := &models.Category{
			Name: "Fitness", // Try to use existing name
		}

		result, err := service.UpdateCategory(cat1.ID, updated)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCategoryAlreadyExists, err)
	})

	t.Run("Success - Update only description", func(t *testing.T) {
		existing := createTestCategory(t, db, "Toys", "Children toys")

		updated := &models.Category{
			Description: "Updated toys description",
		}

		result, err := service.UpdateCategory(existing.ID, updated)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Toys", result.Name) // Name unchanged
		assert.Equal(t, "Updated toys description", result.Description)
	})
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	service := NewCategoryService()

	t.Run("Success - Get existing category", func(t *testing.T) {
		created := createTestCategory(t, db, "Home & Garden", "Home improvement")

		result, err := service.GetCategoryByID(created.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, created.ID, result.ID)
		assert.Equal(t, "Home & Garden", result.Name)
		assert.Equal(t, "Home improvement", result.Description)
	})

	t.Run("Error - Empty ID", func(t *testing.T) {
		result, err := service.GetCategoryByID("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		result, err := service.GetCategoryByID("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrCategoryNotFound, err)
	})
}

func TestCategoryService_GetAllCategories(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	service := NewCategoryService()

	t.Run("Success - Get all categories", func(t *testing.T) {
		createTestCategory(t, db, "Category1", "Description 1")
		createTestCategory(t, db, "Category2", "Description 2")
		createTestCategory(t, db, "Category3", "Description 3")

		result, err := service.GetAllCategories()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
	})

	t.Run("Success - Empty list when no categories", func(t *testing.T) {
		// Clear database
		db.Exec("DELETE FROM categories")

		result, err := service.GetAllCategories()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	service := NewCategoryService()

	t.Run("Success - Delete existing category", func(t *testing.T) {
		created := createTestCategory(t, db, "ToDelete", "Will be deleted")

		err := service.DeleteCategory(created.ID)

		assert.NoError(t, err)

		// Verify it's deleted
		var count int64
		db.Model(&models.Category{}).Where("id = ?", created.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Error - Empty ID", func(t *testing.T) {
		err := service.DeleteCategory("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID is required")
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		err := service.DeleteCategory("non-existent-id")

		assert.Error(t, err)
		assert.Equal(t, ErrCategoryNotFound, err)
	})
}
